package utils

import "github.com/google/uuid"

func GenerateUUID() (string, error) {

	uuidV4, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}
	return uuidV4.String(), nil
}
