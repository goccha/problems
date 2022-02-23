package problems

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/goccha/http-constants/pkg/headers"
	"github.com/goccha/http-constants/pkg/mimetypes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFound(t *testing.T) {
	problem := New("/test", NewCodeProblem("NotFound")).NotFound("Not Found")
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

func TestDefaultProblem_JSON(t *testing.T) {
	problem := New("").Unavailable("unauthorized")
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
	type TestObject struct {
		Name string `json:"name" validate:"required"`
	}
	obj := &TestObject{}
	validate := validator.New()
	err := validate.Struct(obj)
	p := New("/validate", NewBadRequest(err)).BadRequest("bad request")
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
		if len(br.InvalidParams) == 1 {
			params := br.InvalidParams[0]
			if params.Name != "Name" {
				t.Errorf("expect = Name, actual = %s", params.Name)
			}
			if params.Reason != "required" {
				t.Errorf("expect = required, actual = %s", params.Reason)
			}
		}
	}
}

func TestServerProblemOf(t *testing.T) {
	err := errors.New("test error")
	p := ServerProblemOf(context.TODO(), "/problems", err)
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

func TestServerProblemOfNil(t *testing.T) {
	p := ServerProblemOf(context.TODO(), "/problems", nil)
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
