package problems

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type BadRequest interface {
	SetValidationErrors(ve []ValidationError)
	SetInvalidParams(params []InvalidParam)
}

type InvalidParam struct {
	XMLName xml.Name `json:"-" xml:"invalid-param"`
	Name    string   `json:"name" xml:"name"`
	Reason  string   `json:"reason" xml:"reason"`
}

func (ip InvalidParam) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	//return e.EncodeElement(ip, start)
	start.Name.Local = "invalid-param"
	root := []xml.Token{start}
	root = append(root,
		xml.StartElement{Name: xml.Name{Local: "name"}},
		xml.CharData(ip.Name),
		xml.EndElement{Name: xml.Name{Local: "name"}},
		xml.StartElement{Name: xml.Name{Local: "reason"}},
		xml.CharData(ip.Reason),
		xml.EndElement{Name: xml.Name{Local: "reason"}},
	)
	for _, t := range root {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

type ValidationError struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Detail  string   `json:"detail" xml:"detail"`
	Pointer string   `json:"pointer" xml:"pointer"`
}

func (ve ValidationError) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "error"
	root := []xml.Token{start}
	root = append(root,
		xml.StartElement{Name: xml.Name{Local: "detail"}},
		xml.CharData(ve.Detail),
		xml.EndElement{Name: xml.Name{Local: "detail"}},
		xml.StartElement{Name: xml.Name{Local: "pointer"}},
		xml.CharData(ve.Pointer),
		xml.EndElement{Name: xml.Name{Local: "pointer"}},
	)
	for _, t := range root {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}
func (ve ValidationError) Error() string {
	return fmt.Sprintf("Property '%s' does not match the schema", ve.Pointer)
}

// InvalidParams Create RFC7807-style validation error messages
func InvalidParams(err error, params ...InvalidParam) Option {
	var fields []InvalidParam
	var ve validator.ValidationErrors
	var ne *strconv.NumError
	var ute *json.UnmarshalTypeError
	if errors.As(err, &ve) {
		fields = make([]InvalidParam, 0, len(ve))
		for _, v := range ve {
			p := InvalidParam{Name: v.Field(), Reason: v.Tag()}
			fields = append(fields, p)
		}
	} else if errors.As(err, &ne) {
		fields = []InvalidParam{
			{Name: ne.Func, Reason: ne.Num},
		}
	} else if errors.As(err, &ute) {
		fields = []InvalidParam{
			{Name: ute.Field, Reason: "Illegal value type"},
		}
	}
	fields = append(fields, params...)
	return func(p Problem) Problem {
		if err != nil {
			if dp, ok := p.(ProblemDetail); ok {
				dp.SetDetail(err.Error())
			}
		}
		if bp, ok := p.(BadRequest); ok {
			bp.SetInvalidParams(fields)
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
			if len(names) == 1 {
				buf.WriteRune('/')
				buf.WriteString(n)
			}
		} else if len(n) > 0 && n != "-" {
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
func ValidationErrors(err error, verrs ...ValidationError) Option {
	var fields []ValidationError
	var ve validator.ValidationErrors
	var ne *strconv.NumError
	var ute *json.UnmarshalTypeError
	if errors.As(err, &ve) {
		fields = make([]ValidationError, 0, len(ve))
		for _, v := range ve {
			p := ValidationError{Detail: v.Tag(), Pointer: convertNamespaceToJsonPointer(v.Namespace())}
			fields = append(fields, p)
		}
	} else if errors.As(err, &ne) {
		fields = []ValidationError{
			{Detail: ne.Num, Pointer: ne.Func},
		}
	} else if errors.As(err, &ute) {
		fields = []ValidationError{
			{Detail: "Illegal value type", Pointer: convertNamespaceToJsonPointer(ute.Field)},
		}
	}
	fields = append(fields, verrs...)
	return func(p Problem) Problem {
		if err != nil {
			if dp, ok := p.(ProblemDetail); ok {
				dp.SetDetail(err.Error())
			}
		}
		if bp, ok := p.(BadRequest); ok {
			bp.SetValidationErrors(fields)
		}
		return p
	}
}
