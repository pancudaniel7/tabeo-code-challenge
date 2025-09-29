package apperr

type BaseError interface {
	error
	Code() string
	Message() string
	Cause() error
}
