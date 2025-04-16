package apperr

import (
	"errors"
)

var (
	NotFoundError       = errors.New("entity not found")
	DatabaseError       = errors.New("database operation failed")
	UUIDGenerationError = errors.New("uuid generation failed")
	BadUUIDError        = errors.New("invalid uuid")
)
