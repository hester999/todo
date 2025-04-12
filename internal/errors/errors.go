package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	// TODO Нах тебе тут код и Message. Есть для такого врапинг ошибок
	Code    int
	Message string
	Err     error
}

func (err *AppError) Error() string {
	if err.Err == nil {
		return fmt.Sprintf("code=%d, message=%s", err.Code, err.Message)
	}
	return fmt.Sprintf("code=%d, message=%s, cause=%s", err.Code, err.Message, err.Err.Error())
}

func NewInternal(message string, cause error) *AppError {
	if message == "" {
		message = "internal server error"
	}
	return &AppError{Code: http.StatusInternalServerError, Message: message, Err: cause}
}

func NewBadRequest(message string, cause error) *AppError {
	if message == "" {
		message = "bad request"
	}
	return &AppError{Code: http.StatusBadRequest, Message: message, Err: cause}
}

func NewNotFound(message string, cause error) *AppError {
	if message == "" {
		message = "resource not found"
	}
	return &AppError{Code: http.StatusNotFound, Message: message, Err: cause}
}

var NotFoundError = errors.New("entity not found")
