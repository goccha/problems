package problems

import (
	"fmt"
	"net/http"
)

func Type(format string, args ...interface{}) Option {
	return func(p Problem) Problem {
		if dp, ok := p.(ProblemType); ok {
			dp.SetType(fmt.Sprintf(format, args...))
		}
		return p
	}
}

func Tag(tag string) Option {
	return func(p Problem) Problem {
		if dp, ok := p.(ProblemType); ok {
			dp.SetType(fmt.Sprintf("tag:<%s>", tag))
		}
		return p
	}
}

func Title(title string) Option {
	return func(p Problem) Problem {
		if dp, ok := p.(ProblemTitle); ok {
			dp.SetTitle(title)
		}
		return p
	}
}

func Detail(detail string) Option {
	return func(p Problem) Problem {
		if dp, ok := p.(ProblemDetail); ok {
			dp.SetDetail(detail)
		}
		return p
	}
}

func Instance(instance string) Option {
	return func(p Problem) Problem {
		if dp, ok := p.(ProblemInstance); ok {
			dp.SetInstance(instance)
		}
		return p
	}
}

func Path(req *http.Request) Option {
	if req != nil {
		return Instance(req.URL.Path)
	}
	return func(p Problem) Problem {
		return p
	}
}

func Status(status int) Option {
	return func(p Problem) Problem {
		if s, ok := p.(ProblemStatus); ok {
			s.SetStatus(status)
		}
		return p
	}
}

func Extension(key string, value interface{}) Option {
	return func(p Problem) Problem {
		if cp, ok := p.(ProblemExtension); ok {
			cp.Extension(key, value)
		}
		return p
	}
}

func Error(err error) Option {
	return func(p Problem) Problem {
		if wrap, ok := p.(Wrapper); ok {
			wrap.WrapError(err)
		}
		return p
	}
}
