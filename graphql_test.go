package problems

import "testing"

func TestDefaultProblemEncode(t *testing.T) {
	problem := &DefaultProblem{
		Type:     "type",
		Title:    "title",
		Status:   401,
		Detail:   "Unauthorized",
		Instance: "/api/test",
		Code:     "error-code",
	}
	err := problem.Encode()
	if err.Extensions["type"] != "type" {
		t.Errorf("expect = type, actual = %s", err.Extensions["type"])
	}
	if err.Extensions["title"] != "title" {
		t.Errorf("expect = title, actual = %s", err.Extensions["title"])
	}
	if err.Extensions["status"] != 401 {
		t.Errorf("expect = 401, actual = %v", err.Extensions["status"])
	}
	if err.Message != "Unauthorized" {
		t.Errorf("expect = Unauthorized, actual = %s", err.Message)
	}
	if err.Extensions["instance"] != "/api/test" {
		t.Errorf("expect = /api/test, actual = %s", err.Extensions["instance"])
	}
	if err.Extensions["code"] != "error-code" {
		t.Errorf("expect = error-code, actual = %s", err.Extensions["code"])
	}
}

func TestDefaultProblemDecode(t *testing.T) {
	problem := &DefaultProblem{}
	err := GraphQLError{
		Message: "Unauthorized",
		Extensions: map[string]interface{}{
			"type":     "type",
			"title":    "title",
			"status":   401,
			"instance": "/api/test",
			"code":     "error-code",
		},
	}
	problem.Decode(err)
	if problem.Type != "type" {
		t.Errorf("expect = type, actual = %s", problem.Type)
	}
	if problem.Title != "title" {
		t.Errorf("expect = title, actual = %s", problem.Title)
	}
	if problem.Status != 401 {
		t.Errorf("expect = 401, actual = %v", problem.Status)
	}
	if problem.Detail != "Unauthorized" {
		t.Errorf("expect = Unauthorized, actual = %s", problem.Detail)
	}
	if problem.Instance != "/api/test" {
		t.Errorf("expect = /api/test, actual = %s", problem.Instance)
	}
	if problem.Code != "error-code" {
		t.Errorf("expect = error-code, actual = %s", problem.Code)
	}
}
