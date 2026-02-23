package fnError

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type ErrorLog struct {
	ID        uint `gorm:"primaryKey"`
	Message   string
	Detail    string
	CreatedAt time.Time
}

type ErrorHandler struct {
	DB *gorm.DB
}

func NewErrorHandler(db *gorm.DB) *ErrorHandler {
	return &ErrorHandler{DB: db}
}

func (h *ErrorHandler) Handle(_ context.Context, err error) *AppError {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	if !ok {
		appErr = NewInternalError(err, "Internal Server Error")
	}

	switch appErr.Action {
	case ActionPanic:
		h.logToDB(appErr)
		panic(appErr)

	case ActionLogAndReturn:
		h.logToDB(appErr)
		return appErr

	case ActionReturn:
		return appErr
	}

	return appErr
}

func (h *ErrorHandler) logToDB(appErr *AppError) {
	go func() {
		errMsg := appErr.UserMessage
		errDetail := ""
		if appErr.OriginalError != nil {
			errDetail = appErr.OriginalError.Error()
		}

		if err := h.DB.Create(&ErrorLog{
			Message: errMsg,
			Detail:  errDetail,
		}).Error; err != nil {
			log.Printf("Failed to write error log to DB: %v", err)
		}
	}()
}
