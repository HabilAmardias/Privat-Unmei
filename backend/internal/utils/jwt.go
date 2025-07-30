package utils

import (
	"os"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct{}

func CreateJWTUtil() *JWTUtil {
	return &JWTUtil{}
}

func (ju *JWTUtil) GenerateJWT(id string, role int, usedFor int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claim := entity.CustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			Issuer:    "privat-unmei",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		Role: role,
		For:  usedFor,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenstr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", customerrors.NewError(
			"failed to authorize",
			err,
			customerrors.CommonErr,
		)
	}
	return tokenstr, nil
}
