package xerr

import (
	"fmt"
	"net/http"
)

type ErrorAction int

const (
	ActionReturn ErrorAction = iota
	ActionLogAndReturn
	ActionPanic
)

type AppError struct {
	OriginalError error
	UserMessage   string
	Code          int
	Action        ErrorAction
}

func (e *AppError) Error() string {
	if e.OriginalError != nil {
		return fmt.Sprintf("%s: %v", e.UserMessage, e.OriginalError)
	}
	return e.UserMessage
}

func NewBadRequest(msg string) *AppError {
	return &AppError{
		UserMessage: msg,
		Code:        http.StatusBadRequest,
		Action:      ActionReturn,
	}
}

func NewInternalError(err error, msg string) *AppError {
	return &AppError{
		OriginalError: err,
		UserMessage:   msg,
		Code:          http.StatusInternalServerError,
		Action:        ActionLogAndReturn,
	}
}

func NewFatalError(err error, msg string) *AppError {
	return &AppError{
		OriginalError: err,
		UserMessage:   msg,
		Action:        ActionPanic,
	}
}
