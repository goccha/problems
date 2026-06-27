package problems

type CodeProblem interface {
	SetCode(code string)
}

func Code(code string) Option {
	return func(p Problem) Problem {
		if cp, ok := p.(CodeProblem); ok {
			cp.SetCode(code)
		}
		return p
	}
}
