package apperr

import (
	"errors"
)

var (
	NotFoundError       = errors.New("entity not found")
	UUIDGenerationError = errors.New("uuid generation failed")
	BadUUIDError        = errors.New("invalid uuid")
)
