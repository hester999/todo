package utils

import (
	"github.com/google/uuid"
	"todo/internal/apperr"
)

func GenerateUUID() (string, error) {

	uuidV4, err := uuid.NewRandom()

	if err != nil {
		return "", apperr.UUIDGenerationError
	}
	return uuidV4.String(), nil
}
