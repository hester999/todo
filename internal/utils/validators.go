package utils

import (
	"errors"
	"github.com/google/uuid"
)

func ValidateTitle(title string, maxLen int) error {
	if maxLen == 0 {
		maxLen = 100
	}
	if len(title) == 0 {
		return errors.New("invalid title")
	}
	if len(title) > maxLen {
		return errors.New("title too long")
	}
	return nil
}

func ValidateDescription(description string, maxLen int) error {
	if maxLen == 0 {
		maxLen = 100
	}

	if len(description) == 0 {
		return errors.New("invalid description")
	}

	if len(description) > maxLen {
		return errors.New("description too long")

	}
	return nil
}

func UUIDValidator(id string) error {
	_, err := uuid.Parse(id)

	if err != nil {
		return errors.New("invalid uuid")
	}
	return nil
}
