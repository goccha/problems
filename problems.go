package problems

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/goccha/logging/log"
)

const (
	DefaultType = "about:blank"
)

type Problem interface {
	// ProblemStatus returns the status code of the problem.
	// Deprecated: Use StatusCode instead.
	ProblemStatus() int
	// StatusCode returns the status code of the problem.
	StatusCode() int
	Error() string
	Wrap() error
	String() string
	WithMessage(format string, args ...interface{}) Problem
	Message() string
	Renderer
}

type ProblemType interface {
	SetType(url string)
}
type ProblemTitle interface {
	SetTitle(title string)
}
type ProblemDetail interface {
	SetDetail(detail string)
}
type ProblemInstance interface {
	SetInstance(instance string)
}
type ProblemStatus interface {
	SetStatus(status int)
}
type ProblemExtension interface {
	Extension(key string, value interface{})
}

type Wrapper interface {
	WrapError(err error)
	Unwrap() error
}

type Extendable interface {
	Extended() bool
	Map() (map[string]interface{}, bool)
}

type ProblemDetails struct {
	XMLName       xml.Name          `json:"-" xml:"problem"`
	Xmlns         string            `json:"-" xml:"xmlns,attr"`
	Type          string            `json:"type" xml:"type"`
	Title         string            `json:"title" xml:"title,attr,omitempty"`
	Status        int               `json:"status,omitempty" xml:"status,attr,omitempty"`
	Detail        string            `json:"detail,omitempty" xml:"xsd:detail,attr,omitempty"`
	Instance      string            `json:"instance,omitempty" xml:"instance,attr,omitempty"`
	Code          string            `json:"code,omitempty" xml:"code,attr,omitempty"`
	InvalidParams []InvalidParam    `json:"invalid-params,omitempty" xml:"invalid-params,omitempty"`
	Errors        []ValidationError `json:"errors,omitempty" xml:"errors,omitempty"`
	extensions    map[string]interface{}
	err           error
}

// WithMessage sets the detail message of the problem.
func (p *ProblemDetails) WithMessage(format string, args ...interface{}) Problem {
	p.Detail = fmt.Sprintf(format, args...)
	return p
}

// Message returns the detail message of the problem.
func (p *ProblemDetails) Message() string {
	return p.Detail
}

func (p *ProblemDetails) WrapError(err error) {
	p.err = err
}

func (p *ProblemDetails) SetParams(url, detail string) {
	if p.Type == DefaultType {
		p.Type = url
	}
	if detail != "" {
		p.Detail = detail
	}
}
func (p *ProblemDetails) SetType(url string) {
	if url != "" {
		p.Type = url
	}
}
func (p *ProblemDetails) SetTitle(title string) {
	if title != "" {
		p.Title = title
	}
}
func (p *ProblemDetails) SetDetail(detail string) {
	p.Detail = detail
}
func (p *ProblemDetails) SetInstance(instance string) {
	p.Instance = instance
}
func (p *ProblemDetails) SetStatus(status int) {
	p.Status = status
}
func (p *ProblemDetails) SetCode(code string) {
	p.Code = code
}
func (p *ProblemDetails) SetValidationErrors(ve []ValidationError) {
	p.Errors = ve
}
func (p *ProblemDetails) SetInvalidParams(params []InvalidParam) {
	p.InvalidParams = params
}
func (p *ProblemDetails) Extension(key string, value interface{}) {
	if p.extensions == nil {
		p.extensions = make(map[string]interface{})
	}
	p.extensions[key] = value
}

// ProblemStatus returns the status code of the problem.
// Deprecated: Use StatusCode instead.
func (p *ProblemDetails) ProblemStatus() int {
	return p.Status
}

// StatusCode returns the status code of the problem.
func (p *ProblemDetails) StatusCode() int {
	return p.Status
}

// Extended checks if the ProblemDetails struct has any extensions.
func (p *ProblemDetails) Extended() bool {
	if p.extensions == nil {
		return false
	}
	return len(p.extensions) > 0
}

// Map returns a map representation of the ProblemDetails struct.
func (p *ProblemDetails) Map() (map[string]interface{}, bool) {
	m := make(map[string]interface{}, len(p.extensions)+8)
	m["type"] = p.Type
	m["title"] = p.Title
	if p.Detail != "" {
		m["detail"] = p.Detail
	}
	if p.Instance != "" {
		m["instance"] = p.Instance
	}
	if p.Code != "" {
		m["code"] = p.Code
	}
	if len(p.InvalidParams) > 0 {
		m["invalid-params"] = p.InvalidParams
	}
	if len(p.Errors) > 0 {
		m["errors"] = p.Errors
	}
	if p.Status != 0 {
		m["status"] = p.Status
	}
	extended := false
	if p.extensions != nil {
		extended = true
		for k, v := range p.extensions {
			m[k] = v
		}
	}
	return m, extended
}

// JSON marshals the ProblemDetails struct into JSON format.
func (p *ProblemDetails) JSON(ctx context.Context, w http.ResponseWriter) {
	WriteJson(ctx, w, p.StatusCode(), p)
}

// MarshalXML marshals the ProblemDetails struct into XML format.
func (p *ProblemDetails) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	p.Xmlns = Ns9457
	start.Name.Local = "problem"
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	elements := []xml.Token{
		xml.StartElement{Name: xml.Name{Local: "type"}},
		xml.CharData(p.Type),
		xml.EndElement{Name: xml.Name{Local: "type"}},
		xml.StartElement{Name: xml.Name{Local: "title"}},
		xml.CharData(p.Title),
		xml.EndElement{Name: xml.Name{Local: "title"}},
		xml.StartElement{Name: xml.Name{Local: "detail"}},
		xml.CharData(p.Detail),
		xml.EndElement{Name: xml.Name{Local: "detail"}},
		xml.StartElement{Name: xml.Name{Local: "instance"}},
		xml.CharData(p.Instance),
		xml.EndElement{Name: xml.Name{Local: "instance"}},
	}
	if p.Status != 0 {
		elements = append(elements,
			xml.StartElement{Name: xml.Name{Local: "status"}},
			xml.CharData([]byte(fmt.Sprintf("%d", p.Status))),
			xml.EndElement{Name: xml.Name{Local: "status"}},
		)
	}
	if p.Code != "" {
		elements = append(elements,
			xml.StartElement{Name: xml.Name{Local: "code"}},
			xml.CharData([]byte(p.Code)),
			xml.EndElement{Name: xml.Name{Local: "code"}},
		)
	}
	for _, t := range elements {
		if err := e.EncodeToken(t); err != nil {
			return err
		}
	}
	if len(p.InvalidParams) > 0 {
		p.Xmlns = Ns7807
		se := xml.StartElement{Name: xml.Name{Local: "invalid-params"}}
		if err := e.EncodeToken(se); err != nil {
			return err
		}
		for _, param := range p.InvalidParams {
			if err := param.MarshalXML(e, se); err != nil {
				return err
			}
		}
		if err := e.EncodeToken(se.End()); err != nil {
			return err
		}
	}
	if len(p.Errors) > 0 {
		p.Xmlns = Ns9457
		se := xml.StartElement{Name: xml.Name{Local: "errors"}}
		if err := e.EncodeToken(se); err != nil {
			return err
		}
		for _, ve := range p.Errors {
			if err := ve.MarshalXML(e, start); err != nil {
				return err
			}
		}
		if err := e.EncodeToken(se.End()); err != nil {
			return err
		}
	}
	if p.extensions != nil {
		err := EncodeMap(e, p.extensions)
		if err != nil {
			return err
		}
	}
	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}
	return e.Flush()
}

// UnmarshalXML unmarshals the XML data into the ProblemDetails struct.
func (p *ProblemDetails) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	overlay := xmlProblemDetails{}
	if err := d.DecodeElement(&overlay, &start); err != nil {
		return err
	}
	p.Type = overlay.Type
	p.Title = overlay.Title
	p.Detail = overlay.Detail
	p.Instance = overlay.Instance
	p.Status = overlay.Status
	p.Code = overlay.Code
	p.InvalidParams = overlay.InvalidParams
	p.Errors = overlay.Errors
	if overlay.Extensions != nil {
		p.extensions = overlay.Extensions
	}
	return nil
}

// XML marshals the ProblemDetails struct into XML format.
func (p *ProblemDetails) XML(ctx context.Context, w http.ResponseWriter) {
	WriteXml(ctx, w, p.StatusCode(), p)
}

// String returns the JSON representation of the ProblemDetails struct.
func (p *ProblemDetails) String() string {
	bytes, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

// Error returns the JSON representation of the ProblemDetails struct.
func (p *ProblemDetails) Error() string {
	return p.String()
}

// Wrap wraps the ProblemDetails struct with an error.
func (p *ProblemDetails) Wrap() error {
	return p
}

// Unwrap returns the underlying error wrapped by the ProblemDetails struct.
func (p *ProblemDetails) Unwrap() error {
	return p.err
}

// NewDetails creates a new ProblemDetails instance with the given status code.
func NewDetails(status int) *ProblemDetails {
	p := &ProblemDetails{Type: DefaultType}
	p.Title = http.StatusText(status)
	p.Status = status
	return p
}

// As checks if the error is of type Problem and matches the given status codes.
func As(err error, status ...int) (Problem, bool) {
	var pe Problem
	if errors.As(err, &pe) {
		if len(status) > 0 {
			for _, s := range status {
				if pe.StatusCode() == s {
					return pe, true
				}
			}
			return nil, false
		}
		return pe, true
	}
	return nil, false
}

type FromError func(err error) (Problem, bool)

type option struct {
	fromError FromError
	message   string
}

type Arg func(*option)

func WithFromError(f FromError) Arg {
	return func(arg *option) {
		arg.fromError = f
	}
}

// WithMessage sets a custom message for the Problem.
func WithMessage(msg string) Arg {
	return func(arg *option) {
		arg.message = msg
	}
}

func Of(ctx context.Context, path string, err error, args ...Arg) Problem {
	var pe Problem
	if errors.As(err, &pe) {
		if pi, ok := pe.(ProblemInstance); ok {
			pi.SetInstance(path)
		}
		return pe
	}
	opt := &option{}
	for _, arg := range args {
		arg(opt)
	}
	if opt.fromError != nil {
		if p, ok := opt.fromError(err); ok {
			return p
		}
	}
	log.EmbedObject(ctx, log.Error(ctx, 1)).Stack().Err(err).Msgf("%+v", err)
	msg := ""
	if opt.message != "" {
		msg = opt.message
	} else if err != nil {
		msg = err.Error()
	}
	return New(Instance(path), Error(err)).InternalServerError("%s", msg)
}

// UnmarshalJson is a convenience function to unmarshal the json body into a Problem instance.
func UnmarshalJson(ctx context.Context, status int, body []byte, f ...func(status int) Problem) (problem Problem, err error) {
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

// UnmarshalXml is a convenience function to unmarshal the xml body into a Problem instance.
func UnmarshalXml(ctx context.Context, status int, body []byte, f ...func(status int) Problem) (problem Problem, err error) {
	problem = newProblem(status, f...)
	if len(body) <= 0 {
		return
	}
	if err = xml.Unmarshal(body, problem); err != nil {
		return nil, err
	}
	return
}

func newProblem(status int, f ...func(status int) Problem) (problem Problem) {
	if len(f) > 0 {
		problem = f[0](status)
	}
	if problem == nil {
		problem = &ProblemDetails{
			extensions: make(map[string]interface{}),
		}
	}
	return
}

// DecodeJson is a convenience function to decode the io.Reader into a Problem instance.
func DecodeJson(ctx context.Context, status int, body io.Reader, f ...func(status int) Problem) (problem Problem, err error) {
	problem = newProblem(status, f...)
	if body == nil {
		return
	}
	if err = json.NewDecoder(body).Decode(&problem); err != nil {
		_, _ = io.Copy(io.Discard, body) // Discard the body
		return problem, fmt.Errorf("%w", err)
	}
	return
}

// DecodeXml is a convenience function to decode the io.Reader into a Problem instance.
func DecodeXml(ctx context.Context, status int, body io.Reader, f ...func(status int) Problem) (problem Problem, err error) {
	problem = newProblem(status, f...)
	if body == nil {
		return
	}
	if err = xml.NewDecoder(body).Decode(&problem); err != nil {
		_, _ = io.Copy(io.Discard, body) // Discard the body
		return problem, fmt.Errorf("%w", err)
	}
	return
}
