package apperror

import "errors"

type Code string

const (
	CodeInternal           Code = "internal"
	CodeInvalidArgument    Code = "invalid_argument"
	CodeNotFound           Code = "not_found"
	CodeFailedPrecondition Code = "failed_precondition"
	CodeAborted            Code = "aborted"
)

type Error struct {
	Code Code
	Err  error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return string(e.Code)
	}

	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func Wrap(code Code, err error) error {
	if err == nil {
		return nil
	}

	return &Error{Code: code, Err: err}
}

func CodeOf(err error) Code {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code
	}

	return CodeInternal
}
