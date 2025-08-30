package middlewares

import (
	"errors"
	"fmt"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func AuthenticationMiddleware(tokenUtil *utils.JWTUtil, usedFor int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")

			ctx.Error(customerrors.NewError(
				"user credential does not exist",
				err,
				customerrors.Unauthenticate,
			))
			ctx.Abort()
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := errors.New("invalid token format")

			ctx.Error(customerrors.NewError(
				"invalid credential",
				err,
				customerrors.Unauthenticate,
			))
			ctx.Abort()
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)

			ctx.Error(customerrors.NewError(
				"invalid credential",
				err,
				customerrors.Unauthenticate,
			))
			ctx.Abort()
			return
		}

		accessToken := fields[1]
		payload, err := tokenUtil.VerifyJWT(accessToken, usedFor)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set(constants.CTX_AUTH_PAYLOAD_KEY, payload)
		ctx.Set(constants.CTX_AUTH_TOKEN_KEY, accessToken)
		ctx.Next()
	}
}
