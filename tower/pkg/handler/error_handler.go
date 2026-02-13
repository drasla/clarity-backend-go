package handler

import (
	"context"
	"errors"
	"log"
	"time"
	"tower/pkg/xerr"

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

func (h *ErrorHandler) Handle(ctx context.Context, err error) *xerr.AppError {
	var appErr *xerr.AppError
	ok := errors.As(err, &appErr)
	if !ok {
		appErr = xerr.NewInternalError(err, "Internal Server Error")
	}

	switch appErr.Action {
	case xerr.ActionPanic:
		h.logToDB(appErr)
		panic(appErr)

	case xerr.ActionLogAndReturn:
		h.logToDB(appErr)
		return appErr

	case xerr.ActionReturn:
		return appErr
	}

	return appErr
}

func (h *ErrorHandler) logToDB(appErr *xerr.AppError) {
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
