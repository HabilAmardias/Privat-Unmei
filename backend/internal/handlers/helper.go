package handlers

import (
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"regexp"

	"github.com/gin-gonic/gin"
)

var (
	degreelist = []string{"bachelor", "diploma", "high school", "master", "professor"}
)

func getAuthenticationPayload(ctx *gin.Context) (*entity.CustomClaim, error) {

	claims, ok := ctx.Get(constants.CTX_AUTH_PAYLOAD_KEY)
	if !ok {
		return nil, customerrors.NewError(
			"user credential identification failed",
			errors.New("cannot find authentication claim"),
			customerrors.CommonErr,
		)
	}

	customClaims, ok := claims.(*entity.CustomClaim)
	if !ok {
		return nil, customerrors.NewError(
			"user credential identification failed",
			errors.New("cannot parse authentication claim"),
			customerrors.CommonErr,
		)
	}
	return customClaims, nil
}

func ValidateDegree(degree string) bool {
	for _, item := range degreelist {
		if degree == item {
			return true
		}
	}
	return false
}

func ValidatePhoneNumber(phoneNumber string) bool {

	pattern := `^0\d{9,12}$`

	regex := regexp.MustCompile(pattern)

	return regex.MatchString(phoneNumber)
}
