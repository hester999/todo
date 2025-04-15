package utils

import (
	"github.com/google/uuid"
	"todo/internal/apperr"
)

func ValidateTitle(title string, maxLen int) error {
	if maxLen == 0 {
		maxLen = 100
	}
	if len(title) == 0 {
		return apperr.TitleBadRequestError
	}
	if len(title) > maxLen {
		return apperr.TitleTooLongError
	}
	return nil
}

func ValidateDescription(description string, maxLen int) error {
	if maxLen == 0 {
		maxLen = 100
	}

	if len(description) == 0 {
		return apperr.DescriptionBadRequestError
	}

	if len(description) > maxLen {
		return apperr.DescriptionTooLongError

	}
	return nil
}

func UUIDValidator(id string) error {
	_, err := uuid.Parse(id)

	if err != nil {
		return apperr.BadUUIDError
	}
	return nil
}
