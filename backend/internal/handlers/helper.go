package handlers

import (
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"

	"github.com/gin-gonic/gin"
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
