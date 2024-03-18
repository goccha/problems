package problems

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/goccha/http-constants/pkg/mimetypes"
	"github.com/goccha/logging/log"
)

type GraphQLRenderer interface {
	GraphQL(ctx context.Context, w http.ResponseWriter)
}

type GraphQLDecoder interface {
	Decode(ext GraphQLExtension) Problem
}

type GraphQLResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors Errors      `json:"errors,omitempty"`
}

func (p *GraphQLResponse) Decode(ext GraphQLExtension) Problem {
	if len(p.Errors) > 0 {
		return ext.Decode(p.Errors[0])
	}
	if v, ok := ext.(Problem); ok {
		return v
	}
	return nil
}

func (p *GraphQLResponse) ProblemStatus() int {
	for _, err := range p.Errors {
		if v, ok := err.Extensions["status"]; ok {
			if status, ok := v.(int); ok {
				return status
			}
		}
	}
	return http.StatusInternalServerError
}
func (p *GraphQLResponse) Wrap() error {
	return &ProblemError{problem: p}
}
func (p *GraphQLResponse) String() string {
	bytes, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
func (p *GraphQLResponse) JSON(ctx context.Context, w http.ResponseWriter) {
	WriteJson(ctx, w, p.ProblemStatus(), p)
}
func (p *GraphQLResponse) XML(ctx context.Context, w http.ResponseWriter) {
	WriteXml(ctx, w, p.ProblemStatus(), p)
}

type GraphQLError struct {
	Message    string                 `json:"message"`
	Locations  []Location             `json:"locations"`
	Path       []interface{}          `json:"path"`
	Extensions map[string]interface{} `json:"extensions"`
}

func (err *GraphQLError) ProblemStatus() int {
	if v, ok := err.Extensions["status"]; ok {
		if status, ok := v.(int); ok {
			return status
		}
	}
	return http.StatusInternalServerError
}

func (err *GraphQLError) Decode(ext GraphQLExtension) Problem {
	return ext.Decode(*err)
}

func (err *GraphQLError) Problem() Problem {
	switch err.ProblemStatus() {
	case http.StatusBadRequest:
		return err.Decode(&BadRequest{})
	default:
		return err.Decode(&DefaultProblem{})
	}
}

type Location struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Errors []GraphQLError

type GraphQLExtension interface {
	Encode() GraphQLError
	Decode(err GraphQLError) Problem
}

func (p *DefaultProblem) Encode() GraphQLError {
	err := GraphQLError{
		Message: p.Detail,
		Extensions: map[string]interface{}{
			"type":     p.Type,
			"title":    p.Title,
			"status":   p.Status,
			"instance": p.Instance,
		},
	}
	if p.Code != "" {
		err.Extensions["code"] = p.Code
	}
	return err
}
func (p *DefaultProblem) Decode(err GraphQLError) Problem {
	p.Type = err.Extensions["type"].(string)
	p.Title = err.Extensions["title"].(string)
	p.Status = err.Extensions["status"].(int)
	p.Detail = err.Message
	p.Instance = err.Extensions["instance"].(string)
	p.Code = err.Extensions["code"].(string)
	return p
}
func (p *DefaultProblem) GraphQL(ctx context.Context, w http.ResponseWriter) {
	WriteGraphQL(ctx, w, p.ProblemStatus(), p)
}

func WriteGraphQL(ctx context.Context, w http.ResponseWriter, status int, v interface{}) {
	setHeader(ctx, w, status, mimetypes.JSON)
	res := &GraphQLResponse{}
	if encoder, ok := v.(GraphQLExtension); ok {
		err := encoder.Encode()
		res.Errors = []GraphQLError{err}
	} else {
		res.Data = v
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.EmbedObject(ctx, log.Warn(ctx).Err(err)).Send()
	}
}

func (p *BadRequest) Encode() GraphQLError {
	err := p.DefaultProblem.Encode()
	if p.InvalidParams != nil {
		params := make([]map[string]interface{}, 0, len(p.InvalidParams))
		for _, v := range p.InvalidParams {
			params = append(params, map[string]interface{}{
				"name":   v.Name,
				"reason": v.Reason,
			})
		}
		err.Extensions["invalid-params"] = params
	}
	if p.Errors != nil {
		errors := make([]map[string]interface{}, 0, len(p.Errors))
		for _, v := range p.Errors {
			errors = append(errors, map[string]interface{}{
				"detail":  v.Detail,
				"pointer": v.Pointer,
			})
		}
		err.Extensions["errors"] = errors
	}
	return err
}
func (p *BadRequest) Decode(err GraphQLError) Problem {
	p.DefaultProblem.Decode(err)
	if v, ok := err.Extensions["invalid-params"]; ok {
		params := v.([]map[string]interface{})
		for _, param := range params {
			p.InvalidParams = append(p.InvalidParams, InvalidParam{
				Name:   param["name"].(string),
				Reason: param["reason"].(string),
			})
		}
	}
	if v, ok := err.Extensions["errors"]; ok {
		params := v.([]map[string]interface{})
		for _, param := range params {
			p.Errors = append(p.Errors, ValidationError{
				Detail:  param["detail"].(string),
				Pointer: param["pointer"].(string),
			})
		}
	}
	return p
}
func (p *BadRequest) GraphQL(ctx context.Context, w http.ResponseWriter) {
	WriteGraphQL(ctx, w, p.ProblemStatus(), p)
}
