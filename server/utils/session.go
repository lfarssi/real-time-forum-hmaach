package utils

import "github.com/gofrs/uuid"

func GenerateToken() (string, error) {
	tokenID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return tokenID.String(), nil
}
