package problems

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/goccha/errors"
	"github.com/goccha/http-constants/pkg/headers"
	"github.com/goccha/http-constants/pkg/mimetypes"
	"github.com/goccha/stackdriver/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

const (
	DefaultType = "about:blank"
)

type Renderer interface {
	JSON(ctx context.Context, w http.ResponseWriter)
	XML(ctx context.Context, w http.ResponseWriter)
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
	p.Type = url
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
	w.Header().Set(headers.ContentType, mimetypes.ProblemJson)
	if p.ProblemStatus() > 0 {
		w.WriteHeader(p.ProblemStatus())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Warn(ctx).Err(err).Send()
	}
}
func (p *DefaultProblem) XML(ctx context.Context, w http.ResponseWriter) {
	w.Header().Set(headers.ContentType, mimetypes.ProblemXml)
	if p.ProblemStatus() > 0 {
		w.WriteHeader(p.ProblemStatus())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := xml.NewEncoder(w).Encode(p); err != nil {
		log.Warn(ctx).Err(err).Send()
	}
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
	p.Title = http.StatusText(int(status))
	p.Status = status
	return p
}

type ProblemError struct {
	Path    string
	problem Problem
	Err     error
}

func (err *ProblemError) Problem() Problem {
	if err.Err != nil {
		return New(err.Path, nil).InternalServerError(err.Err.Error())
	}
	if v, ok := err.problem.(DefaultParams); ok {
		if err.Path != "" {
			v.SetInstance(err.Path)
		}
	}
	return err.problem
}
func (err *ProblemError) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.problem.String()
}

func New(path string, f ...func(p *DefaultProblem) Problem) *Builder {
	b := &Builder{
		url:  DefaultType,
		path: path,
	}
	if f != nil && len(f) > 0 {
		b.f = f[0]
	}
	return b
}

type Builder struct {
	url  string
	path string
	f    func(p *DefaultProblem) Problem
}

func (b *Builder) Type(format string, args ...interface{}) *Builder {
	b.url = fmt.Sprintf(format, args...)
	return b
}
func (b *Builder) build(status int, detail string, f func(p *DefaultProblem) Problem) (sp Problem) {
	if f != nil {
		sp = b.f(NewProblem(status))
	} else {
		sp = NewProblem(status)
	}
	if dp, ok := sp.(DefaultParams); ok {
		dp.SetParams(b.url, b.path, detail)
	}
	return sp
}
func (b *Builder) BadRequest(format string, args ...interface{}) Problem {
	return b.build(http.StatusBadRequest, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Unauthorized(format string, args ...interface{}) Problem {
	return b.build(http.StatusUnauthorized, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) PaymentRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusPaymentRequired, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Forbidden(format string, args ...interface{}) Problem {
	return b.build(http.StatusForbidden, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) NotFound(format string, args ...interface{}) Problem {
	return b.build(http.StatusNotFound, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) MethodNotAllowed(format string, args ...interface{}) Problem {
	return b.build(http.StatusMethodNotAllowed, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) NotAcceptable(format string, args ...interface{}) Problem {
	return b.build(http.StatusNotAcceptable, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) ProxyAuthRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusProxyAuthRequired, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) RequestTimeout(format string, args ...interface{}) Problem {
	return b.build(http.StatusRequestTimeout, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Conflict(format string, args ...interface{}) Problem {
	return b.build(http.StatusConflict, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Gone(format string, args ...interface{}) Problem {
	return b.build(http.StatusGone, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) LengthRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusLengthRequired, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) PreconditionFailed(format string, args ...interface{}) Problem {
	return b.build(http.StatusPreconditionFailed, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) RequestEntityTooLarge(format string, args ...interface{}) Problem {
	return b.build(http.StatusRequestEntityTooLarge, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) RequestURITooLong(format string, args ...interface{}) Problem {
	return b.build(http.StatusRequestURITooLong, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) UnsupportedMediaType(format string, args ...interface{}) Problem {
	return b.build(http.StatusUnsupportedMediaType, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) RequestedRangeNotSatisfiable(format string, args ...interface{}) Problem {
	return b.build(http.StatusRequestedRangeNotSatisfiable, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) ExpectationFailed(format string, args ...interface{}) Problem {
	return b.build(http.StatusExpectationFailed, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Teapot(format string, args ...interface{}) Problem {
	return b.build(http.StatusTeapot, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) MisdirectedRequest(format string, args ...interface{}) Problem {
	return b.build(http.StatusMisdirectedRequest, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) UnprocessableEntity(format string, args ...interface{}) Problem {
	return b.build(http.StatusUnprocessableEntity, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Locked(format string, args ...interface{}) Problem {
	return b.build(http.StatusLocked, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) FailedDependency(format string, args ...interface{}) Problem {
	return b.build(http.StatusFailedDependency, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) TooEarly(format string, args ...interface{}) Problem {
	return b.build(http.StatusTooEarly, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) UpgradeRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusUpgradeRequired, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) PreconditionRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusPreconditionRequired, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) TooManyRequests(format string, args ...interface{}) Problem {
	return b.build(http.StatusTooManyRequests, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) RequestHeaderFieldsTooLarge(format string, args ...interface{}) Problem {
	return b.build(http.StatusRequestHeaderFieldsTooLarge, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) UnavailableForLegalReasons(format string, args ...interface{}) Problem {
	return b.build(http.StatusUnavailableForLegalReasons, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) InternalServerError(format string, args ...interface{}) Problem {
	return b.build(http.StatusInternalServerError, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) NotImplemented(format string, args ...interface{}) Problem {
	return b.build(http.StatusNotImplemented, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) BadGateway(format string, args ...interface{}) Problem {
	return b.build(http.StatusBadGateway, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Unavailable(format string, args ...interface{}) Problem {
	return b.build(http.StatusServiceUnavailable, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) GatewayTimeout(format string, args ...interface{}) Problem {
	return b.build(http.StatusGatewayTimeout, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) HTTPVersionNotSupported(format string, args ...interface{}) Problem {
	return b.build(http.StatusHTTPVersionNotSupported, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) VariantAlsoNegotiates(format string, args ...interface{}) Problem {
	return b.build(http.StatusVariantAlsoNegotiates, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) InsufficientStorage(format string, args ...interface{}) Problem {
	return b.build(http.StatusInsufficientStorage, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) LoopDetected(format string, args ...interface{}) Problem {
	return b.build(http.StatusLoopDetected, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) NotExtended(format string, args ...interface{}) Problem {
	return b.build(http.StatusNotExtended, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) NetworkAuthenticationRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusNetworkAuthenticationRequired, fmt.Sprintf(format, args...), b.f)
}

type BadRequest struct {
	*DefaultProblem
	InvalidParams []InvalidParam `json:"invalid-params,omitempty"`
}

func NewBadRequest(err error) func(p *DefaultProblem) Problem {
	var fields []InvalidParam
	switch err.(type) {
	case validator.ValidationErrors:
		ve := err.(validator.ValidationErrors)
		fields = make([]InvalidParam, 0, len(ve))
		for _, v := range ve {
			p := InvalidParam{v.Field(), v.Tag()}
			fields = append(fields, p)
		}
	case *strconv.NumError:
		ne := err.(*strconv.NumError)
		fields = []InvalidParam{
			{ne.Func, ne.Num},
		}
	}
	return func(p *DefaultProblem) Problem {
		if err != nil {
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

func NewCodeProblem(code string) func(p *DefaultProblem) Problem {
	return func(p *DefaultProblem) Problem {
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

func ClientProblemOf(ctx context.Context, path string, err error) Problem {
	if err == nil {
		return New(path, NewBadRequest(err)).BadRequest("")
	}
	return New(path, NewBadRequest(err)).BadRequest(err.Error())
}

func ServerProblemOf(ctx context.Context, path string, err error) Problem {
	switch err.(type) {
	case *ProblemError:
		p := err.(*ProblemError)
		if p.Path == "" {
			p.Path = path
		}
		return p.Problem()
	default:
		st, ok := status.FromError(errors.Cause(err))
		if ok {
			switch st.Code() {
			case codes.Unavailable:
				log.Warn(ctx, 1).Msgf("%+v", err)
				return New(path, nil).Unavailable("%v", err)
			}
		}
		log.Error(ctx, 1).Msgf("%+v", err)
		if err != nil {
			return New(path, nil).InternalServerError(err.Error())
		}
		return New(path, nil).InternalServerError("")
	}
}
