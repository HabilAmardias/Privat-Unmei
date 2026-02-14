package middlewares

import (
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware(tokenUtil *utils.JWTUtil, usedFor int, cookieKey string, payloadKey string, tokenKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := ctx.Cookie(cookieKey)
		if err != nil {
			ctx.Error(customerrors.NewError(
				"credentials does not found",
				err,
				customerrors.Unauthenticate,
			))
			ctx.Abort()
			return
		}
		payload, err := tokenUtil.VerifyJWT(accessToken, usedFor)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set(payloadKey, payload)
		ctx.Set(tokenKey, accessToken)
		ctx.Next()
	}
}
