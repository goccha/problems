package problems

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"github.com/goccha/http-constants/pkg/headers"
	"github.com/goccha/http-constants/pkg/mimetypes"
)

func TestNotFound(t *testing.T) {
	problem := New(Instance("/test"), Code("NotFound")).NotFound("Not Found")
	if code, ok := problem.(*CodeProblem); !ok {
		t.Errorf("invalid struct. %v", code)
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
	p := New(Instance("/validate"), ValidationErrors(err), InvalidParams(err)).BadRequest("bad request")
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
	if cp, ok := p.(*CodeProblem); ok {
		if cp.Type != "http://localhost:8080/test?code=E001" {
			t.Errorf("expect = http://localhost:8080/test?code=E001, actual = %s", cp.Type)
		}
		if cp.Code != "E001" {
			t.Errorf("expect = E001, actual = %s", cp.Code)
		}
	} else {
		t.Errorf("expect = CodeProblem, actual=%v", cp)
	}
}
