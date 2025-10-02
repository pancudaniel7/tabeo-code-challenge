package apperr

import (
	"errors"
	"fmt"
)

type Code string

const (
	InternalError        Code = "INTERNAL_ERROR"
	NotFound             Code = "NOT_FOUND"
	AlreadyExists        Code = "ALREADY_EXISTS"
	BadGatewayError      Code = "BAD_GATEWAY"
	InvalidArgumentError Code = "INVALID_ARGUMENT"
)

type Error struct {
	code Code
	msg  string
	err  error
}

func New(code Code, msg string, cause error) *Error  { return &Error{code: code, msg: msg, err: cause} }
func Internal(msg string, cause error) *Error        { return New(InternalError, msg, cause) }
func NotFoundErr(msg string, cause error) *Error     { return New(NotFound, msg, cause) }
func Exists(msg string, cause error) *Error          { return New(AlreadyExists, msg, cause) }
func BadGateway(msg string, cause error) *Error      { return New(BadGatewayError, msg, cause) }
func InvalidArgument(msg string, cause error) *Error { return New(InvalidArgumentError, msg, cause) }

func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.code, e.msg, e.err)
	}
	return fmt.Sprintf("[%s] %s", e.code, e.msg)
}

func (e *Error) Unwrap() error   { return e.err }
func (e *Error) Code() Code      { return e.code }
func (e *Error) Message() string { return e.msg }

func IsCode(err error, code Code) bool {
	var ce *Error
	if errors.As(err, &ce) {
		return ce.code == code
	}
	return false
}

func IsInvalid(err error) bool         { return IsCode(err, InvalidArgumentError) }
func IsInternal(err error) bool        { return IsCode(err, InternalError) }
func IsNotFound(err error) bool        { return IsCode(err, NotFound) }
func IsExists(err error) bool          { return IsCode(err, AlreadyExists) }
func IsBadGateway(err error) bool      { return IsCode(err, BadGatewayError) }
func IsInvalidArgument(err error) bool { return IsCode(err, InvalidArgumentError) }
