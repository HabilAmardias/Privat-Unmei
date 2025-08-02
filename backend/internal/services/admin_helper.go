package services

import (
	"crypto/rand"
	"math/big"
)

const (
	letters  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits   = "0123456789"
	specials = "!@#"
	allChars = letters + digits + specials
	length   = 10
)

func generateRandomPassword() (string, error) {
	password := make([]byte, length)
	specialIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(specials))))
	if err != nil {
		return "", err
	}
	password[0] = specials[specialIndex.Int64()]

	digitIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
	if err != nil {
		return "", err
	}
	password[1] = digits[digitIndex.Int64()]

	for i := 2; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[randomIndex.Int64()]
	}

	return string(password), nil
}
