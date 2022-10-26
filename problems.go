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
	SetParams(url, instance, detail string)
	SetType(url string)
	SetTitle(title string)
	SetDetail(detail string)
	SetInstance(instance string)
}

type DefaultProblem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func (p *DefaultProblem) SetParams(url, instance, detail string) {
	if p.Type == DefaultType {
		p.Type = url
	}
	p.Instance = instance
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
	return &ProblemError{P: p}
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

type ProblemError struct {
	Path string
	P    Problem
	Err  error
}

func (err *ProblemError) Problem() Problem {
	if err.Err != nil {
		return New(err.Path, nil).InternalServerError(err.Err.Error())
	}
	if v, ok := err.P.(DefaultParams); ok {
		if err.Path != "" {
			v.SetInstance(err.Path)
		}
	}
	return err.P
}
func (err *ProblemError) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.P.String()
}

type BadRequest struct {
	*DefaultProblem
	InvalidParams []InvalidParam `json:"invalid-params,omitempty"`
}

func (p *BadRequest) JSON(ctx context.Context, w http.ResponseWriter) {
	WriteJson(ctx, w, p.ProblemStatus(), p)
}
func (p *BadRequest) XML(ctx context.Context, w http.ResponseWriter) {
	WriteXml(ctx, w, p.ProblemStatus(), p)
}
func (p *BadRequest) Wrap() error {
	return &ProblemError{P: p}
}

func NewBadRequest(err error, params ...InvalidParam) func(p *DefaultProblem) Problem {
	var fields []InvalidParam
	switch err := err.(type) {
	case validator.ValidationErrors:
		fields = make([]InvalidParam, 0, len(err))
		for _, v := range err {
			p := InvalidParam{v.Field(), v.Tag()}
			fields = append(fields, p)
		}
	case *strconv.NumError:
		fields = []InvalidParam{
			{err.Func, err.Num},
		}
	case *json.UnmarshalTypeError:
		fields = []InvalidParam{
			{err.Field, "Illegal value type"},
		}
	}
	fields = append(fields, params...)
	return func(p *DefaultProblem) Problem {
		if p.Detail == "" && err != nil {
			p.Detail = err.Error()
		}
		return &BadRequest{
			DefaultProblem: p,
			InvalidParams:  fields,
		}
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
	return &ProblemError{P: p}
}

func NewCodeProblem(code string, typ ...string) func(p *DefaultProblem) Problem {
	return func(p *DefaultProblem) Problem {
		if len(typ) > 0 {
			p.Type = typ[0]
		}
		return &CodeProblem{
			DefaultProblem: p,
			Code:           code,
		}
	}
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
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

func ServerProblemOf(ctx context.Context, path string, err error, f ...MsgFunc) Problem {
	switch err := err.(type) {
	case *ProblemError:
		if err.Path == "" {
			err.Path = path
		}
		return err.Problem()
	default:
		msg := selectMsg(err, f...)
		if st, ok := status.FromError(errors.Unwrap(err)); ok {
			switch st.Code() {
			case codes.Unavailable:
				log.EmbedObject(ctx, log.Warn(ctx, 1)).Msgf("%+v", err)
				return New(path).Unavailable(msg())
			}
		}
		log.EmbedObject(ctx, log.Error(ctx, 1)).Msgf("%+v", err)
		return New(path).InternalServerError(msg())
	}
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
