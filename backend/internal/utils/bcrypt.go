package utils

import (
	"privat-unmei/internal/customerrors"

	"golang.org/x/crypto/bcrypt"
)

type BcryptUtil struct{}

func (bu *BcryptUtil) HashPassword(plainPass string) (string, error) {
	strBytes := []byte(plainPass)
	hashed, err := bcrypt.GenerateFromPassword(strBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", customerrors.NewError(
			"failed to save password",
			err,
			customerrors.CommonErr,
		)
	}
	return string(hashed), nil
}

func (bu *BcryptUtil) ComparePassword(plainPass string, hashedPass string) bool {
	hashedBytes := []byte(hashedPass)
	plainBytes := []byte(plainPass)
	return bcrypt.CompareHashAndPassword(hashedBytes, plainBytes) == nil
}
