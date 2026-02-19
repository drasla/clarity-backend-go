package cError

type CustomError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *CustomError) Error() string {
	return e.Message
}

func New(code ErrorCode, message string) *CustomError {
	return &CustomError{Code: code, Message: message}
}

func Wrap(err error, code ErrorCode, msg string) *CustomError {
	return &CustomError{Code: code, Message: msg, Err: err}
}
