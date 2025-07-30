package entity

import "github.com/golang-jwt/jwt/v5"

type (
	CustomClaim struct {
		jwt.RegisteredClaims
		Role int
		For  int
	}
)
