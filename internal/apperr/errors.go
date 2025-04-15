package apperr

import (
	"errors"
)

//type AppError struct {
//	// TODO Нах тебе тут код и Message. Есть для такого врапинг ошибок
//	Code    int
//	Message string
//	Err     error
//}
//
//func (apperr *AppError) Error() string {
//	if apperr.Err == nil {
//		return fmt.Sprintf("code=%d, message=%s", apperr.Code, apperr.Message)
//	}
//	return fmt.Sprintf("code=%d, message=%s, cause=%s", apperr.Code, apperr.Message, apperr.Err.Error())
//}
//
//func NewInternal(message string, cause error) *AppError {
//	if message == "" {
//		message = "internal server error"
//	}
//	return &AppError{Code: http.StatusInternalServerError, Message: message, Err: cause}
//}
//
//func NewBadRequest(message string, cause error) *AppError {
//	if message == "" {
//		message = "bad request"
//	}
//	return &AppError{Code: http.StatusBadRequest, Message: message, Err: cause}
//}
//
//func NewNotFound(message string, cause error) *AppError {
//	if message == "" {
//		message = "resource not found"
//	}
//	return &AppError{Code: http.StatusNotFound, Message: message, Err: cause}
//}

var (
	NotFoundError              = errors.New("entity not found")
	TitleBadRequestError       = errors.New("invalid title")
	DescriptionBadRequestError = errors.New("invalid description")
	DuplicateIdError           = errors.New("task with this id already exists")
	DatabaseError              = errors.New("database operation failed")
	InvalidStatusError         = errors.New("invalid status")
	UUIDGenerationError        = errors.New("uuid generation failed")
	BadUUIDError               = errors.New("invalid uuid")
	TitleTooLongError          = errors.New("title too long")
	DescriptionTooLongError    = errors.New("description too short")
)
