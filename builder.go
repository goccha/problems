package problems

import (
	"fmt"
	"net/http"
)

type Builder struct {
	url  string
	path string
	f    func(p *DefaultProblem) Problem
}

func New(path string, f ...func(p *DefaultProblem) Problem) *Builder {
	b := &Builder{
		url:  DefaultType,
		path: path,
	}
	if len(f) > 0 {
		b.f = f[0]
	}
	return b
}

func (b *Builder) Type(format string, args ...interface{}) *Builder {
	b.url = fmt.Sprintf(format, args...)
	return b
}
func (b *Builder) build(status int, detail string, f func(p *DefaultProblem) Problem) (sp Problem) {
	if f != nil {
		sp = b.f(NewProblem(status))
	} else {
		sp = NewProblem(status)
	}
	if dp, ok := sp.(DefaultParams); ok {
		dp.SetParams(b.url, b.path, detail)
	}
	return sp
}
func (b *Builder) BadRequest(format string, args ...interface{}) Problem {
	return b.build(http.StatusBadRequest, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Unauthorized(format string, args ...interface{}) Problem {
	return b.build(http.StatusUnauthorized, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) PaymentRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusPaymentRequired, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Forbidden(format string, args ...interface{}) Problem {
	return b.build(http.StatusForbidden, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) NotFound(format string, args ...interface{}) Problem {
	return b.build(http.StatusNotFound, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) MethodNotAllowed(format string, args ...interface{}) Problem {
	return b.build(http.StatusMethodNotAllowed, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) NotAcceptable(format string, args ...interface{}) Problem {
	return b.build(http.StatusNotAcceptable, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) ProxyAuthRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusProxyAuthRequired, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) RequestTimeout(format string, args ...interface{}) Problem {
	return b.build(http.StatusRequestTimeout, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Conflict(format string, args ...interface{}) Problem {
	return b.build(http.StatusConflict, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Gone(format string, args ...interface{}) Problem {
	return b.build(http.StatusGone, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) LengthRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusLengthRequired, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) PreconditionFailed(format string, args ...interface{}) Problem {
	return b.build(http.StatusPreconditionFailed, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) RequestEntityTooLarge(format string, args ...interface{}) Problem {
	return b.build(http.StatusRequestEntityTooLarge, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) RequestURITooLong(format string, args ...interface{}) Problem {
	return b.build(http.StatusRequestURITooLong, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) UnsupportedMediaType(format string, args ...interface{}) Problem {
	return b.build(http.StatusUnsupportedMediaType, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) RequestedRangeNotSatisfiable(format string, args ...interface{}) Problem {
	return b.build(http.StatusRequestedRangeNotSatisfiable, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) ExpectationFailed(format string, args ...interface{}) Problem {
	return b.build(http.StatusExpectationFailed, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Teapot(format string, args ...interface{}) Problem {
	return b.build(http.StatusTeapot, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) MisdirectedRequest(format string, args ...interface{}) Problem {
	return b.build(http.StatusMisdirectedRequest, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) UnprocessableEntity(format string, args ...interface{}) Problem {
	return b.build(http.StatusUnprocessableEntity, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Locked(format string, args ...interface{}) Problem {
	return b.build(http.StatusLocked, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) FailedDependency(format string, args ...interface{}) Problem {
	return b.build(http.StatusFailedDependency, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) TooEarly(format string, args ...interface{}) Problem {
	return b.build(http.StatusTooEarly, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) UpgradeRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusUpgradeRequired, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) PreconditionRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusPreconditionRequired, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) TooManyRequests(format string, args ...interface{}) Problem {
	return b.build(http.StatusTooManyRequests, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) RequestHeaderFieldsTooLarge(format string, args ...interface{}) Problem {
	return b.build(http.StatusRequestHeaderFieldsTooLarge, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) UnavailableForLegalReasons(format string, args ...interface{}) Problem {
	return b.build(http.StatusUnavailableForLegalReasons, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) InternalServerError(format string, args ...interface{}) Problem {
	return b.build(http.StatusInternalServerError, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) NotImplemented(format string, args ...interface{}) Problem {
	return b.build(http.StatusNotImplemented, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) BadGateway(format string, args ...interface{}) Problem {
	return b.build(http.StatusBadGateway, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) Unavailable(format string, args ...interface{}) Problem {
	return b.build(http.StatusServiceUnavailable, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) GatewayTimeout(format string, args ...interface{}) Problem {
	return b.build(http.StatusGatewayTimeout, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) HTTPVersionNotSupported(format string, args ...interface{}) Problem {
	return b.build(http.StatusHTTPVersionNotSupported, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) VariantAlsoNegotiates(format string, args ...interface{}) Problem {
	return b.build(http.StatusVariantAlsoNegotiates, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) InsufficientStorage(format string, args ...interface{}) Problem {
	return b.build(http.StatusInsufficientStorage, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) LoopDetected(format string, args ...interface{}) Problem {
	return b.build(http.StatusLoopDetected, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) NotExtended(format string, args ...interface{}) Problem {
	return b.build(http.StatusNotExtended, fmt.Sprintf(format, args...), b.f)
}
func (b *Builder) NetworkAuthenticationRequired(format string, args ...interface{}) Problem {
	return b.build(http.StatusNetworkAuthenticationRequired, fmt.Sprintf(format, args...), b.f)
}
