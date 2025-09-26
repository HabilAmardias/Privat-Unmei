package middlewares

import (
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware(tokenUtil *utils.JWTUtil, usedFor int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := ctx.Cookie(constants.AUTH_COOKIE_KEY)
		if err != nil{
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

		ctx.Set(constants.CTX_AUTH_PAYLOAD_KEY, payload)
		ctx.Set(constants.CTX_AUTH_TOKEN_KEY, accessToken)
		ctx.Next()
	}
}

func WSAuthenticationMiddleware(tokenUtil *utils.JWTUtil, usedFor int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.AuthenticationReq
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		payload, err := tokenUtil.VerifyJWT(req.Token, usedFor)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set(constants.CTX_AUTH_PAYLOAD_KEY, payload)
		ctx.Set(constants.CTX_AUTH_TOKEN_KEY, req.Token)
		ctx.Next()
	}
}
