package problems

import (
	"fmt"
	"net/http"
)

type Builder struct {
	url     string
	options []Option
}

type Option func(p Problem) Problem

func New(opt ...Option) *Builder {
	b := &Builder{
		url: DefaultType,
	}
	b.options = opt
	return b
}

// Type sets the type of the problem
func (b *Builder) Type(format string, args ...interface{}) *Builder {
	b.url = fmt.Sprintf(format, args...)
	return b
}
func (b *Builder) build(status int, detail string) (p Problem) {
	p = NewDetails(status)
	if dp, ok := p.(ProblemType); ok {
		dp.SetType(b.url)
	}
	for _, f := range b.options {
		p = f(p)
	}
	if detail != "" {
		if dp, ok := p.(ProblemDetail); ok {
			dp.SetDetail(detail)
		}
	}
	return p
}
func (b *Builder) Build(status int, format string, args ...interface{}) Problem {
	return b.build(status, fmt.Sprintf(format, args...))
}
func (b *Builder) BadRequest(format string, args ...interface{}) Problem {
	return b.Build(http.StatusBadRequest, format, args...)
}
func (b *Builder) Unauthorized(format string, args ...interface{}) Problem {
	return b.Build(http.StatusUnauthorized, format, args...)
}
func (b *Builder) PaymentRequired(format string, args ...interface{}) Problem {
	return b.Build(http.StatusPaymentRequired, format, args...)
}
func (b *Builder) Forbidden(format string, args ...interface{}) Problem {
	return b.Build(http.StatusForbidden, format, args...)
}
func (b *Builder) NotFound(format string, args ...interface{}) Problem {
	return b.Build(http.StatusNotFound, format, args...)
}
func (b *Builder) MethodNotAllowed(format string, args ...interface{}) Problem {
	return b.Build(http.StatusMethodNotAllowed, format, args...)
}
func (b *Builder) NotAcceptable(format string, args ...interface{}) Problem {
	return b.Build(http.StatusNotAcceptable, format, args...)
}
func (b *Builder) ProxyAuthRequired(format string, args ...interface{}) Problem {
	return b.Build(http.StatusProxyAuthRequired, format, args...)
}
func (b *Builder) RequestTimeout(format string, args ...interface{}) Problem {
	return b.Build(http.StatusRequestTimeout, format, args...)
}
func (b *Builder) Conflict(format string, args ...interface{}) Problem {
	return b.Build(http.StatusConflict, format, args...)
}
func (b *Builder) Gone(format string, args ...interface{}) Problem {
	return b.Build(http.StatusGone, format, args...)
}
func (b *Builder) LengthRequired(format string, args ...interface{}) Problem {
	return b.Build(http.StatusLengthRequired, format, args...)
}
func (b *Builder) PreconditionFailed(format string, args ...interface{}) Problem {
	return b.Build(http.StatusPreconditionFailed, format, args...)
}
func (b *Builder) RequestEntityTooLarge(format string, args ...interface{}) Problem {
	return b.Build(http.StatusRequestEntityTooLarge, format, args...)
}
func (b *Builder) RequestURITooLong(format string, args ...interface{}) Problem {
	return b.Build(http.StatusRequestURITooLong, format, args...)
}
func (b *Builder) UnsupportedMediaType(format string, args ...interface{}) Problem {
	return b.Build(http.StatusUnsupportedMediaType, format, args...)
}
func (b *Builder) RequestedRangeNotSatisfiable(format string, args ...interface{}) Problem {
	return b.Build(http.StatusRequestedRangeNotSatisfiable, format, args...)
}
func (b *Builder) ExpectationFailed(format string, args ...interface{}) Problem {
	return b.Build(http.StatusExpectationFailed, format, args...)
}
func (b *Builder) Teapot(format string, args ...interface{}) Problem {
	return b.Build(http.StatusTeapot, format, args...)
}
func (b *Builder) MisdirectedRequest(format string, args ...interface{}) Problem {
	return b.Build(http.StatusMisdirectedRequest, format, args...)
}
func (b *Builder) UnprocessableEntity(format string, args ...interface{}) Problem {
	return b.Build(http.StatusUnprocessableEntity, format, args...)
}
func (b *Builder) Locked(format string, args ...interface{}) Problem {
	return b.Build(http.StatusLocked, format, args...)
}
func (b *Builder) FailedDependency(format string, args ...interface{}) Problem {
	return b.Build(http.StatusFailedDependency, format, args...)
}
func (b *Builder) TooEarly(format string, args ...interface{}) Problem {
	return b.Build(http.StatusTooEarly, format, args...)
}
func (b *Builder) UpgradeRequired(format string, args ...interface{}) Problem {
	return b.Build(http.StatusUpgradeRequired, format, args...)
}
func (b *Builder) PreconditionRequired(format string, args ...interface{}) Problem {
	return b.Build(http.StatusPreconditionRequired, format, args...)
}
func (b *Builder) TooManyRequests(format string, args ...interface{}) Problem {
	return b.Build(http.StatusTooManyRequests, format, args...)
}
func (b *Builder) RequestHeaderFieldsTooLarge(format string, args ...interface{}) Problem {
	return b.Build(http.StatusRequestHeaderFieldsTooLarge, format, args...)
}
func (b *Builder) UnavailableForLegalReasons(format string, args ...interface{}) Problem {
	return b.Build(http.StatusUnavailableForLegalReasons, format, args...)
}
func (b *Builder) InternalServerError(format string, args ...interface{}) Problem {
	return b.Build(http.StatusInternalServerError, format, args...)
}
func (b *Builder) NotImplemented(format string, args ...interface{}) Problem {
	return b.Build(http.StatusNotImplemented, format, args...)
}
func (b *Builder) BadGateway(format string, args ...interface{}) Problem {
	return b.Build(http.StatusBadGateway, format, args...)
}

// Unavailable is an alias for ServiceUnavailable
// Deprecated: Use ServiceUnavailable instead
func (b *Builder) Unavailable(format string, args ...interface{}) Problem {
	return b.ServiceUnavailable(format, args...)
}
func (b *Builder) ServiceUnavailable(format string, args ...interface{}) Problem {
	return b.Build(http.StatusServiceUnavailable, format, args...)
}
func (b *Builder) GatewayTimeout(format string, args ...interface{}) Problem {
	return b.Build(http.StatusGatewayTimeout, format, args...)
}
func (b *Builder) HTTPVersionNotSupported(format string, args ...interface{}) Problem {
	return b.Build(http.StatusHTTPVersionNotSupported, format, args...)
}
func (b *Builder) VariantAlsoNegotiates(format string, args ...interface{}) Problem {
	return b.Build(http.StatusVariantAlsoNegotiates, format, args...)
}
func (b *Builder) InsufficientStorage(format string, args ...interface{}) Problem {
	return b.Build(http.StatusInsufficientStorage, format, args...)
}
func (b *Builder) LoopDetected(format string, args ...interface{}) Problem {
	return b.Build(http.StatusLoopDetected, format, args...)
}
func (b *Builder) NotExtended(format string, args ...interface{}) Problem {
	return b.Build(http.StatusNotExtended, format, args...)
}
func (b *Builder) NetworkAuthenticationRequired(format string, args ...interface{}) Problem {
	return b.Build(http.StatusNetworkAuthenticationRequired, format, args...)
}
