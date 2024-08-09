package problems

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/goccha/http-constants/pkg/headers"
	"github.com/goccha/http-constants/pkg/mimetypes"
)

func TestNotFound(t *testing.T) {
	problem := New(Instance("/test"), Code("NotFound")).NotFound("Not Found")
	if code, ok := problem.(*DefaultProblem); !ok {
		t.Errorf("invalid struct. %v", problem)
	} else {
		if code.Instance != "/test" {
			t.Errorf("expect = /test, actual = %s", code.Instance)
		}
		if code.Detail != "Not Found" {
			t.Errorf("expect = Not Found, actual = %s", code.Detail)
		}
		if code.Type != DefaultType {
			t.Errorf("expect = %s, actual = %s", DefaultType, code.Type)
		}
		if code.Status != http.StatusNotFound {
			t.Errorf("expect = %d, actual = %d", http.StatusNotFound, code.Status)
		}
		expect := http.StatusText(http.StatusNotFound)
		if code.Title != expect {
			t.Errorf("expect = %s, actual = %s", expect, code.Title)
		}
		if code.Code != "NotFound" {
			t.Errorf("expect = NotFound, actual = %s", code.Code)
		}
	}
}

func TestNewBadRequest(t *testing.T) {
	p := New(InvalidParams(nil,
		InvalidParam{
			Name:   "X-Test-Key",
			Reason: "required",
		})).BadRequest("")

	bp := p.(*BadRequest)
	assert.Equal(t, 1, len(bp.InvalidParams))
}

func TestDefaultProblem_JSON(t *testing.T) {
	problem := New().Unavailable("unauthorized")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		problem.JSON(context.TODO(), w)
	}))
	defer server.Close()
	if res, err := http.Get(server.URL + "/"); err != nil {
		t.Errorf("%v", err)
	} else {
		actual := res.Header.Get(headers.ContentType)
		if actual != mimetypes.ProblemJson {
			t.Errorf("expect = %s, actual = %s", mimetypes.ProblemJson, actual)
		}
		if res.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("expect = %d, actual = %d", http.StatusServiceUnavailable, res.StatusCode)
		}
	}
}

func TestBuilder_BadRequest(t *testing.T) {
	type NestedObject struct {
		Name string `json:"name" validate:"required"`
	}
	type TestObject struct {
		Name string         `json:"name" validate:"required"`
		List []NestedObject `json:"nested" validate:"required,dive"`
	}
	obj := &TestObject{
		List: []NestedObject{
			{},
		},
	}
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		return name
	})
	err := validate.Struct(obj)
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/validate", nil)
	p := New(Path(req), ValidationErrors(err), InvalidParams(err)).BadRequest("bad request")
	if br, ok := p.(*BadRequest); ok {
		if br.Instance != "/validate" {
			t.Errorf("expect = /validate, actual = %s", br.Instance)
		}
		if br.Detail != "bad request" {
			t.Errorf("expect = bad request, actual = %s", br.Detail)
		}
		if br.Type != DefaultType {
			t.Errorf("expect = %s, actual = %s", DefaultType, br.Type)
		}
		if br.Status != http.StatusBadRequest {
			t.Errorf("expect = %d, actual = %d", http.StatusBadRequest, br.Status)
		}
		expect := http.StatusText(http.StatusBadRequest)
		if br.Title != expect {
			t.Errorf("expect = %s, actual = %s", expect, br.Title)
		}
		if len(br.InvalidParams) == 2 {
			param := br.InvalidParams[0]
			if param.Name != "name" {
				t.Errorf("expect = Name, actual = %s", param.Name)
			}
			if param.Reason != "required" {
				t.Errorf("expect = required, actual = %s", param.Reason)
			}
			param = br.InvalidParams[1]
			if param.Name != "name" {
				t.Errorf("expect = name, actual = %s", param.Name)
			}
			if param.Reason != "required" {
				t.Errorf("expect = required, actual = %s", param.Reason)
			}
		} else {
			t.Errorf("invalid length. %v", br.InvalidParams)
		}
		if len(br.Errors) == 2 {
			list := br.Errors
			if list[0].Detail != "required" {
				t.Errorf("expect = required is a required field, actual = %s", br.Errors[0].Detail)
			}
			if list[0].Pointer != "#/name" {
				t.Errorf("expect = #/name is a required field, actual = %s", br.Errors[0].Pointer)
			}
			if list[1].Detail != "required" {
				t.Errorf("expect = required is a required field, actual = %s", br.Errors[1].Detail)
			}
			if list[1].Pointer != "#/nested/0/name" {
				t.Errorf("expect = #/nested/0/name is a required field, actual = %s", br.Errors[1].Pointer)
			}
		} else {
			t.Errorf("invalid length. %v", br.Errors)

		}
	}
}

func TestOf(t *testing.T) {
	err := errors.New("test error")
	p := Of(context.TODO(), "/problems", err)
	if dp, ok := p.(*DefaultProblem); ok {
		if dp.Instance != "/problems" {
			t.Errorf("expect = /problems, actual = %s", dp.Instance)
		}
		if dp.Detail != "test error" {
			t.Errorf("expect = test error, actual = %s", dp.Detail)
		}
		if dp.Type != DefaultType {
			t.Errorf("expect = %s, actual = %s", DefaultType, dp.Type)
		}
		if dp.Status != http.StatusInternalServerError {
			t.Errorf("expect = %d, actual = %d", http.StatusInternalServerError, dp.Status)
		}
		expect := http.StatusText(http.StatusInternalServerError)
		if dp.Title != expect {
			t.Errorf("expect = %s, actual = %s", expect, dp.Title)
		}
	}
}

func TestOfNil(t *testing.T) {
	p := Of(context.TODO(), "/problems", nil)
	if dp, ok := p.(*DefaultProblem); ok {
		if dp.Instance != "/problems" {
			t.Errorf("expect = /problems, actual = %s", dp.Instance)
		}
		if dp.Detail != "" {
			t.Errorf("expect = '', actual = '%s'", dp.Detail)
		}
		if dp.Type != DefaultType {
			t.Errorf("expect = %s, actual = %s", DefaultType, dp.Type)
		}
		if dp.Status != http.StatusInternalServerError {
			t.Errorf("expect = %d, actual = %d", http.StatusInternalServerError, dp.Status)
		}
		expect := http.StatusText(http.StatusInternalServerError)
		if dp.Title != expect {
			t.Errorf("expect = %s, actual = %s", expect, dp.Title)
		}
	}
}

func TestCodeProblem(t *testing.T) {
	p := New(Instance("/problems"), Code("E001"), Type("http://localhost:8080/test?code=E001")).Unavailable("")
	if dp, ok := p.(*DefaultProblem); ok {
		if dp.Type != "http://localhost:8080/test?code=E001" {
			t.Errorf("expect = http://localhost:8080/test?code=E001, actual = %s", dp.Type)
		}
		if dp.Code != "E001" {
			t.Errorf("expect = E001, actual = %s", dp.Code)
		}
	} else {
		t.Errorf("expect = DefaultProblem, actual=%v", p)
	}
}

func TestBind(t *testing.T) {
	p := New(Instance("/problems")).BadRequest("bad request")
	bin, err := json.Marshal(p)
	if err != nil {
		t.Errorf("%v", err)
	}
	if bp, err := Bind(context.TODO(), http.StatusBadRequest, bin); err != nil {
		t.Errorf("%v", err)
	} else {
		if bp.ProblemStatus() != http.StatusBadRequest {
			t.Errorf("expect = %d, actual = %d", http.StatusBadRequest, bp.ProblemStatus())
		}
		req, ok := bp.(*BadRequest)
		if !ok {
			t.Errorf("invalid struct. %v", bp)
		}
		if req.Instance != "/problems" {
			t.Errorf("expect = /problems, actual = %s", req.Instance)
		}
		if req.Detail != "bad request" {
			t.Errorf("expect = bad request, actual = %s", req.Detail)
		}
	}
}

func TestDecode(t *testing.T) {
	p := New(Instance("/problems")).BadRequest("bad request")
	bin, err := json.Marshal(p)
	if err != nil {
		t.Errorf("%v", err)
	}
	if bp, err := Decode(context.TODO(), http.StatusBadRequest, bytes.NewBuffer(bin)); err != nil {
		t.Errorf("%v", err)
	} else {
		if bp.ProblemStatus() != http.StatusBadRequest {
			t.Errorf("expect = %d, actual = %d", http.StatusBadRequest, bp.ProblemStatus())
		}
		req, ok := bp.(*BadRequest)
		if !ok {
			t.Errorf("invalid struct. %v", bp)
		}
		if req.Instance != "/problems" {
			t.Errorf("expect = /problems, actual = %s", req.Instance)
		}
		if req.Detail != "bad request" {
			t.Errorf("expect = bad request, actual = %s", req.Detail)
		}
	}
}

func TestNamespace(t *testing.T) {
	type Test struct {
		Value string `json:"value" validate:"required"`
	}

	tv := &Test{}
	validation := validator.New(validator.WithRequiredStructEnabled())
	validation.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	if err := validation.Struct(tv); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, v := range ve {
				jp := convertNamespaceToJsonPointer(v.Namespace())
				if jp != "#/value" {
					t.Errorf("expect = #/value, actual = %s", jp)
				}
			}
		}
	}
}

func TestValidationErrors(t *testing.T) {
	verr := validator.ValidationErrors{}
	verr = append(verr, &fieldError{tag: "a", ns: "a"})
	problem := &DefaultProblem{}
	p := ValidationErrors(verr)(problem)
	if v, ok := p.(*BadRequest); !ok {
		t.Errorf("expect = BadRequest, actual = %v", p)
	} else {
		if len(v.Errors) != 1 {
			t.Errorf("expect = error, actual = %v", v)
		}
	}
	var err error
	_, err = strconv.ParseInt("a", 10, 32)
	problem = &DefaultProblem{}
	p = ValidationErrors(err)(problem)
	if v, ok := p.(*BadRequest); !ok {
		t.Errorf("expect = BadRequest, actual = %v", p)
	} else {
		if len(v.Errors) != 1 {
			t.Errorf("expect = error, actual = %v", v)
		}
	}
	err = &json.UnmarshalTypeError{
		Value:  "",
		Type:   reflect.TypeOf(map[string]interface{}{}),
		Offset: 0,
		Struct: "",
		Field:  "test",
	}
	problem = &DefaultProblem{}
	p = ValidationErrors(err)(problem)
	if v, ok := p.(*BadRequest); !ok {
		t.Errorf("expect = BadRequest, actual = %v", p)
	} else {
		if len(v.Errors) != 1 {
			t.Errorf("expect = error, actual = %v", v)
		}
	}

}

func TestInvalidParams(t *testing.T) {
	verr := validator.ValidationErrors{}
	verr = append(verr, &fieldError{tag: "a", ns: "a"})
	problem := &DefaultProblem{}
	p := InvalidParams(verr)(problem)
	if v, ok := p.(*BadRequest); !ok {
		t.Errorf("expect = BadRequest, actual = %v", p)
	} else {
		if len(v.InvalidParams) != 1 {
			t.Errorf("expect = error, actual = %v", v)
		}
	}
	var err error
	_, err = strconv.ParseInt("a", 10, 32)
	problem = &DefaultProblem{}
	p = InvalidParams(err)(problem)
	if v, ok := p.(*BadRequest); !ok {
		t.Errorf("expect = BadRequest, actual = %v", p)
	} else {
		if len(v.InvalidParams) != 1 {
			t.Errorf("expect = error, actual = %v", v)
		}
	}
	err = &json.UnmarshalTypeError{
		Value:  "",
		Type:   reflect.TypeOf(map[string]interface{}{}),
		Offset: 0,
		Struct: "",
		Field:  "test",
	}
	problem = &DefaultProblem{}
	p = InvalidParams(err)(problem)
	if v, ok := p.(*BadRequest); !ok {
		t.Errorf("expect = BadRequest, actual = %v", p)
	} else {
		if len(v.InvalidParams) != 1 {
			t.Errorf("expect = error, actual = %v", v)
		}
	}

}

type fieldError struct {
	tag            string
	actualTag      string
	ns             string
	structNs       string
	fieldLen       uint8
	structfieldLen uint8
	value          interface{}
	param          string
	kind           reflect.Kind
	typ            reflect.Type
}

// Tag returns the validation tag that failed.
func (fe *fieldError) Tag() string {
	return fe.tag
}

// ActualTag returns the validation tag that failed, even if an
// alias the actual tag within the alias will be returned.
func (fe *fieldError) ActualTag() string {
	return fe.actualTag
}

// Namespace returns the namespace for the field error, with the tag
// name taking precedence over the field's actual name.
func (fe *fieldError) Namespace() string {
	return fe.ns
}

// StructNamespace returns the namespace for the field error, with the field's
// actual name.
func (fe *fieldError) StructNamespace() string {
	return fe.structNs
}

// Field returns the field's name with the tag name taking precedence over the
// field's actual name.
func (fe *fieldError) Field() string {

	return fe.ns[len(fe.ns)-int(fe.fieldLen):]
	// // return fe.field
	// fld := fe.ns[len(fe.ns)-int(fe.fieldLen):]

	// log.Println("FLD:", fld)

	// if len(fld) > 0 && fld[:1] == "." {
	// 	return fld[1:]
	// }

	// return fld
}

// StructField returns the field's actual name from the struct, when able to determine.
func (fe *fieldError) StructField() string {
	// return fe.structField
	return fe.structNs[len(fe.structNs)-int(fe.structfieldLen):]
}

// Value returns the actual field's value in case needed for creating the error
// message
func (fe *fieldError) Value() interface{} {
	return fe.value
}

// Param returns the param value, in string form for comparison; this will
// also help with generating an error message
func (fe *fieldError) Param() string {
	return fe.param
}

// Kind returns the Field's reflect Kind
func (fe *fieldError) Kind() reflect.Kind {
	return fe.kind
}

// Type returns the Field's reflect Type
func (fe *fieldError) Type() reflect.Type {
	return fe.typ
}

// Error returns the fieldError's error message
func (fe *fieldError) Error() string {
	return ""
}

// Translate returns the FieldError's translated error
// from the provided 'ut.Translator' and registered 'TranslationFunc'
//
// NOTE: if no registered translation can be found, it returns the original
// untranslated error message.
func (fe *fieldError) Translate(ut ut.Translator) string {
	return ""
}
