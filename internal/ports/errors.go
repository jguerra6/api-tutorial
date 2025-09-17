package ports

import (
	"errors"
	"fmt"
)

type Code string

const (
	CodeValidation      Code = "validation"
	CodeUnauthenticated Code = "unauthenticated"
	CodeUnauthorized    Code = "unauthorized"
	CodeNotFound        Code = "not_found"
	CodeConflict        Code = "conflict"
	CodeTooManyRequests Code = "too_many_requests"
	CodeExternal        Code = "external"
	CodeInternal        Code = "internal"
	CodeRateLimited     Code = "rate_limited"
	CodeForbidden       Code = "forbidden"
	CodeInvalid         Code = "invalid"
)

type AppError struct {
	Code   Code
	Op     string
	Msg    string
	Err    error
	Fields map[string]string
}

func (e *AppError) Error() string {
	switch {
	case e.Op != "" && e.Msg != "":
		return fmt.Sprintf("%s: %s", e.Op, e.Msg)
	case e.Msg != "":
		return e.Msg
	case e.Op != "":
		return e.Op
	default:
		return string(e.Code)
	}
}

func (e *AppError) Unwrap() error { return e.Err }

func NewAppError(code Code, msg string) *AppError {
	return &AppError{
		Code: code,
		Msg:  msg,
	}
}
func Wrap(code Code, msg string, err error) *AppError {
	return &AppError{Code: code, Msg: msg, Err: err}
}
func WithOp(err error, op string) *AppError {
	var ae *AppError
	if errors.As(err, &ae) {
		ae.Op = op
		return ae
	}
	return &AppError{Code: CodeInternal, Op: op, Msg: err.Error(), Err: err}
}
func IsCode(err error, code Code) bool {
	var ae *AppError
	return errors.As(err, &ae) && ae.Code == code
}
func NewValidationError(msg string) *AppError      { return NewAppError(CodeValidation, msg) }
func NewUnauthenticatedError(msg string) *AppError { return NewAppError(CodeUnauthenticated, msg) }
func NewUnauthorizedError(msg string) *AppError    { return NewAppError(CodeUnauthorized, msg) }
func NewNotFoundError(msg string) *AppError        { return NewAppError(CodeNotFound, msg) }
func NewConflictError(msg string) *AppError        { return NewAppError(CodeConflict, msg) }
func NewTooManyRequestsError(msg string) *AppError { return NewAppError(CodeTooManyRequests, msg) }
func NewExternalError(msg string) *AppError        { return NewAppError(CodeExternal, msg) }
func NewInternalError(msg string) *AppError        { return NewAppError(CodeInternal, msg) }
