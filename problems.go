package problems

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/goccha/http-constants/pkg/headers"
	"github.com/goccha/http-constants/pkg/mimetypes"
	"github.com/goccha/logging/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultType = "about:blank"
)

type Renderer interface {
	JSON(ctx context.Context, w http.ResponseWriter)
	XML(ctx context.Context, w http.ResponseWriter)
	Wrap() error
}

func setHeader(ctx context.Context, w http.ResponseWriter, status int, mimetype string) {
	w.Header().Set(headers.ContentType, mimetype)
	if status > 0 {
		w.WriteHeader(status)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func WriteJson(ctx context.Context, w http.ResponseWriter, status int, v interface{}) {
	setHeader(ctx, w, status, mimetypes.ProblemJson)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.EmbedObject(ctx, log.Warn(ctx).Err(err)).Send()
	}
}

func WriteXml(ctx context.Context, w http.ResponseWriter, status int, v interface{}) {
	setHeader(ctx, w, status, mimetypes.ProblemXml)
	if err := xml.NewEncoder(w).Encode(v); err != nil {
		log.EmbedObject(ctx, log.Warn(ctx).Err(err)).Send()
	}
}

type Problem interface {
	ProblemStatus() int
	Wrap() error
	String() string
	Renderer
}

type DefaultParams interface {
	SetParams(url, detail string)
	SetType(url string)
	SetTitle(title string)
	SetDetail(detail string)
	SetInstance(instance string)
	Problem
}

type DefaultProblem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func (p *DefaultProblem) SetParams(url, detail string) {
	if p.Type == DefaultType {
		p.Type = url
	}
	if detail != "" {
		p.Detail = detail
	}
}
func (p *DefaultProblem) SetType(url string) {
	p.Type = url
}
func (p *DefaultProblem) SetTitle(title string) {
	p.Title = title
}
func (p *DefaultProblem) SetDetail(detail string) {
	p.Detail = detail
}
func (p *DefaultProblem) SetInstance(instance string) {
	p.Instance = instance
}
func (p *DefaultProblem) ProblemStatus() int {
	return p.Status
}

func (p *DefaultProblem) JSON(ctx context.Context, w http.ResponseWriter) {
	WriteJson(ctx, w, p.ProblemStatus(), p)
}
func (p *DefaultProblem) XML(ctx context.Context, w http.ResponseWriter) {
	WriteXml(ctx, w, p.ProblemStatus(), p)
}
func (p *DefaultProblem) Wrap() error {
	return &ProblemError{problem: p}
}
func (p *DefaultProblem) String() string {
	bytes, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func NewProblem(status int) *DefaultProblem {
	p := &DefaultProblem{Type: DefaultType}
	p.Title = http.StatusText(status)
	p.Status = status
	return p
}

func Wrap(p Problem) error {
	return &ProblemError{problem: p}
}

func WrapError(err error) error {
	return &ProblemError{err: err}
}

type ProblemError struct {
	Path    string
	problem Problem
	err     error
}

func (err *ProblemError) Problem() Problem {
	if err.err != nil {
		return New(Instance(err.Path)).InternalServerError(err.err.Error())
	}
	if v, ok := err.problem.(DefaultParams); ok {
		if err.Path != "" {
			v.SetInstance(err.Path)
		}
	}
	return err.problem
}
func (err *ProblemError) Error() string {
	if err.err != nil {
		return err.err.Error()
	}
	return err.problem.String()
}
func (err *ProblemError) Unwrap() error {
	if err.err != nil {
		return err.err
	}
	return nil
}

type BadRequest struct {
	*DefaultProblem
	InvalidParams []InvalidParam    `json:"invalid-params,omitempty"`
	Errors        []ValidationError `json:"errors,omitempty"`
}

func (p *BadRequest) JSON(ctx context.Context, w http.ResponseWriter) {
	WriteJson(ctx, w, p.ProblemStatus(), p)
}
func (p *BadRequest) XML(ctx context.Context, w http.ResponseWriter) {
	WriteXml(ctx, w, p.ProblemStatus(), p)
}
func (p *BadRequest) Wrap() error {
	return &ProblemError{problem: p}
}

// InvalidParams Create RFC7807-style validation error messages
func InvalidParams(err error, params ...InvalidParam) Option {
	var fields []InvalidParam
	ve := &validator.ValidationErrors{}
	ne := &strconv.NumError{}
	ute := &json.UnmarshalTypeError{}
	if errors.As(err, ve) {
		fields = make([]InvalidParam, 0, len(*ve))
		for _, v := range *ve {
			p := InvalidParam{v.Field(), v.Tag()}
			fields = append(fields, p)
		}
	} else if errors.As(err, &ne) {
		fields = []InvalidParam{
			{ne.Func, ne.Num},
		}
	} else if errors.As(err, &ute) {
		fields = []InvalidParam{
			{ute.Field, "Illegal value type"},
		}
	}
	fields = append(fields, params...)
	return func(p DefaultParams) Problem {
		if err != nil {
			p.SetDetail(err.Error())
		}
		switch bp := p.(type) {
		case *BadRequest:
			bp.InvalidParams = append(bp.InvalidParams, fields...)
			return bp
		case *DefaultProblem:
			return &BadRequest{
				DefaultProblem: p.(*DefaultProblem),
				InvalidParams:  fields,
			}
		}
		return p
	}
}

func convertNamespaceToJsonPointer(namespace string) string {
	names := strings.Split(namespace, ".")
	buf := strings.Builder{}
	for i, n := range names {
		if i == 0 {
			buf.WriteRune('#')
		} else {
			buf.WriteString("/")
			if strings.HasSuffix(n, "]") {
				n = strings.ReplaceAll(n, "[", "/")
				buf.WriteString(strings.ReplaceAll(n, "]", ""))
			} else {
				buf.WriteString(n)
			}
		}
	}
	return buf.String()
}

// ValidationErrors Create RFC9457-style validation error messages.
func ValidationErrors(err error, validErrors ...ValidationError) Option {
	var fields []ValidationError
	ve := &validator.ValidationErrors{}
	ne := &strconv.NumError{}
	ute := &json.UnmarshalTypeError{}
	if errors.As(err, ve) {
		fields = make([]ValidationError, 0, len(*ve))
		for _, v := range *ve {
			p := ValidationError{v.Tag(), convertNamespaceToJsonPointer(v.Namespace())}
			fields = append(fields, p)
		}
	} else if errors.As(err, &ne) {
		fields = []ValidationError{
			{ne.Num, ne.Func},
		}
	} else if errors.As(err, &ute) {
		fields = []ValidationError{
			{"Illegal value type", ute.Field},
		}
	}
	fields = append(fields, validErrors...)
	return func(p DefaultParams) Problem {
		if err != nil {
			p.SetDetail(err.Error())
		}
		switch bp := p.(type) {
		case *BadRequest:
			bp.Errors = append(bp.Errors, fields...)
			return bp
		case *DefaultProblem:
			return &BadRequest{
				DefaultProblem: p.(*DefaultProblem),
				Errors:         fields,
			}
		}
		return p
	}
}

type CodeProblem struct {
	*DefaultProblem
	Code string `json:"code"`
}

func (p *CodeProblem) JSON(ctx context.Context, w http.ResponseWriter) {
	WriteJson(ctx, w, p.ProblemStatus(), p)
}
func (p *CodeProblem) XML(ctx context.Context, w http.ResponseWriter) {
	WriteXml(ctx, w, p.ProblemStatus(), p)
}
func (p *CodeProblem) Wrap() error {
	return &ProblemError{problem: p}
}

func Code(code string) Option {
	return func(p DefaultParams) Problem {
		switch dp := p.(type) {
		case *CodeProblem:
			dp.Code = code
			return dp
		case *DefaultProblem:
			return &CodeProblem{
				DefaultProblem: dp,
				Code:           code,
			}
		}
		return p
	}
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type ValidationError struct {
	Detail  string `json:"detail"`
	Pointer string `json:"pointer"`
}

type MsgFunc func() string

func selectMsg(err error, f ...MsgFunc) MsgFunc {
	if len(f) > 0 {
		return f[0]
	} else {
		return func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}
	}
}

func Of(ctx context.Context, path string, err error, f ...MsgFunc) Problem {
	pe := &ProblemError{}
	if errors.As(err, &pe) {
		if pe.Path == "" {
			pe.Path = path
		}
		return pe.Problem()
	}
	msg := selectMsg(err, f...)
	if st, ok := status.FromError(errors.Unwrap(err)); ok {
		switch st.Code() {
		case codes.Unavailable:
			log.EmbedObject(ctx, log.Warn(ctx, 1)).Stack().Msgf("%+v", err)
			return New(Instance(path)).Unavailable(msg())
		}
	}
	log.EmbedObject(ctx, log.Error(ctx, 1)).Stack().Err(err).Msgf("%+v", err)
	return New(Instance(path)).InternalServerError(msg())
}

func Bind(ctx context.Context, status int, body []byte, f ...func(status int) Problem) (problem Problem, err error) {
	problem = newProblem(status, f...)
	if len(body) <= 0 {
		return
	}
	if err = json.Unmarshal(body, problem); err != nil {
		log.Error(ctx).Msg(string(body))
		return problem, fmt.Errorf("%w", err)
	}
	return
}

func newProblem(status int, f ...func(status int) Problem) (problem Problem) {
	if len(f) > 0 {
		problem = f[0](status)
	}
	if problem == nil {
		switch status {
		case http.StatusBadRequest:
			problem = &BadRequest{}
		default:
			problem = &DefaultProblem{}
		}
	}
	return
}

func Decode(ctx context.Context, status int, body io.Reader, f ...func(status int) Problem) (problem Problem, err error) {
	problem = newProblem(status, f...)
	if body == nil {
		return
	}
	if err = json.NewDecoder(body).Decode(&problem); err != nil {
		_, _ = io.Copy(io.Discard, body)
		return problem, fmt.Errorf("%w", err)
	}
	return
}
