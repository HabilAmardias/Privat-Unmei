package handlers

import (
	"errors"
	"fmt"
	"privat-unmei/internal/customerrors"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func GetJWT(ctx *gin.Context) *string {
	authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

	if len(authorizationHeader) == 0 {
		err := errors.New("authorization header is not provided")

		ctx.Error(customerrors.NewError(
			"authorization header not found",
			err,
			customerrors.Unauthenticate,
		))
		ctx.Abort()
		return nil
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) != 2 {
		err := errors.New("invalid token format")

		ctx.Error(customerrors.NewError(
			"invalid authorization header format",
			err,
			customerrors.Unauthenticate,
		))
		ctx.Abort()
		return nil
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		err := fmt.Errorf("unsupported authorization type %s", authorizationType)

		ctx.Error(customerrors.NewError(
			"unsupported authorization type",
			err,
			customerrors.Unauthenticate,
		))
		ctx.Abort()
		return nil
	}

	accessToken := fields[1]
	return &accessToken
}
