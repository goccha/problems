package problems

import (
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// DefaultProblem
// Deprecated: use ProblemDetail instead
type DefaultProblem = ProblemDetail

// ProblemError
// Deprecated: use ProblemDetails instead
type ProblemError = ProblemDetails

// Problem 旧バージョン互換用
// Deprecated: use ProblemDetails instead
func (p *ProblemDetails) Problem() Problem {
	if p.err != nil {
		return New(Instance(p.Instance)).InternalServerError("%s", p.err.Error())
	}
	return p
}

// MsgFunc 詳細メッセージ設定用
// Deprecated: use ProblemDetails instead
type MsgFunc func() string

// selectMsg エラーメッセージ
// Deprecated
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

// ServerProblemOf 旧バージョン互換用
// Deprecated: use Of instead
func ServerProblemOf(ctx context.Context, path string, err error, f ...MsgFunc) Problem {
	msg := selectMsg(err, f...)
	return Of(ctx, path, err).WithMessage("%s", msg())
}

// NewBadRequest
// Deprecated: use ValidationErrors instead
func NewBadRequest(err error, params ...InvalidParam) func(p *ProblemDetails) Problem {
	var fields []InvalidParam
	switch e := err.(type) {
	case validator.ValidationErrors:
		fields = make([]InvalidParam, 0, len(e))
		for _, v := range e {
			p := InvalidParam{Name: v.Field(), Reason: v.Tag()}
			fields = append(fields, p)
		}
	case *strconv.NumError:
		fields = []InvalidParam{
			{Name: e.Func, Reason: e.Num},
		}
	case *json.UnmarshalTypeError:
		fields = []InvalidParam{
			{Name: e.Field, Reason: "Illegal value type"},
		}
	}
	fields = append(fields, params...)
	return func(p *ProblemDetails) Problem {
		if p.Detail == "" && err != nil {
			p.Detail = err.Error()
		}
		p.InvalidParams = fields
		return p
	}
}

// NewCodeProblem
// Deprecated: use Code instead
func NewCodeProblem(code string, typ ...string) func(p *ProblemDetails) Problem {
	return func(p *ProblemDetails) Problem {
		if len(typ) > 0 {
			p.Type = typ[0]
		}
		p.Code = code
		return p
	}
}

// Bind is a convenience function to bind the request body to a Problem instance.
// Deprecated: Use UnmarshalJson instead.
func Bind(ctx context.Context, status int, body []byte, f ...func(status int) Problem) (problem Problem, err error) {
	return UnmarshalJson(ctx, status, body, f...)
}

// Decode is a convenience function to decode the request body into a Problem instance.
// Deprecated: Use DecodeJson instead.
func Decode(ctx context.Context, status int, body io.Reader, f ...func(status int) Problem) (problem Problem, err error) {
	return DecodeJson(ctx, status, body, f...)
}
