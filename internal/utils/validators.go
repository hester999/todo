package utils

import (
	"github.com/google/uuid"
	"todo/internal/apperr"
)

func UUIDValidator(id string) error {
	_, err := uuid.Parse(id)

	if err != nil {
		return apperr.BadUUIDError
	}
	return nil
}
