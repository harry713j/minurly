package utils

import (
	"crypto/rand"
	"encoding/base64"
)


func GenerateRandomStrings(length int) (string, error){
	byteLength := (length * 6 + 7) / 8
	bytes := make([]byte, byteLength)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return  base64.RawURLEncoding.EncodeToString(bytes)[:length] , nil
}