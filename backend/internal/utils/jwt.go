package utils

import (
	"fmt"
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

func (ju *JWTUtil) VerifyJWT(tokenStr string, usedFor int) (*entity.CustomClaim, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	custClaim := new(entity.CustomClaim)

	token, err := jwt.ParseWithClaims(tokenStr, custClaim, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, customerrors.NewError(
				"Failed to authorize",
				fmt.Errorf("failed to parse jwt"),
				customerrors.Unauthenticate,
			)
		}
		return []byte(jwtSecret), nil
	},
		jwt.WithIssuer("privat-unmei"),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return nil, customerrors.NewError(
			"Failed to authorize",
			err,
			customerrors.Unauthenticate,
		)
	}
	if !token.Valid {
		return nil, customerrors.NewError(
			"Failed to authorize",
			fmt.Errorf("jwt token is not valid"),
			customerrors.Unauthenticate,
		)
	}

	claim, ok := token.Claims.(*entity.CustomClaim)
	if !ok {
		return nil, customerrors.NewError(
			"Failed to authorize",
			fmt.Errorf("failed to get jwt claim"),
			customerrors.Unauthenticate,
		)
	}
	exp, expErr := claim.GetExpirationTime()
	if expErr != nil {
		return nil, customerrors.NewError(
			"Failed to authorize",
			expErr,
			customerrors.Unauthenticate,
		)
	}
	if float64(time.Now().Unix()) > float64(exp.Time.Unix()) {
		return nil, customerrors.NewError(
			"session expired",
			fmt.Errorf("jwt is expired"),
			customerrors.Unauthenticate,
		)
	}
	if claim.For != usedFor {
		return nil, customerrors.NewError(
			"Failed to authorize",
			fmt.Errorf("wrong token usage"),
			customerrors.Unauthenticate,
		)
	}
	return claim, nil
}

func (ju *JWTUtil) GenerateJWT(id string, role int, usedFor int, status string, age time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claim := entity.CustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			Issuer:    "privat-unmei",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(age)),
		},
		Role:   role,
		For:    usedFor,
		Status: status,
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
