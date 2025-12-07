package middlewares

import (
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/utils"

	"github.com/gin-gonic/gin"
)

func RefreshAuthMiddleware(tokenUtil *utils.JWTUtil) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		refreshToken, err := ctx.Cookie(constants.REFRESH_COOKIE_KEY)
		if err != nil {
			ctx.Error(customerrors.NewError(
				"credentials does not found",
				err,
				customerrors.Unauthenticate,
			))
			ctx.Abort()
			return
		}
		payload, err := tokenUtil.VerifyJWT(refreshToken, constants.ForRefresh)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set(constants.CTX_REFRESH_PAYLOAD_KEY, payload)
		ctx.Next()
	}
}

func AuthenticationMiddleware(tokenUtil *utils.JWTUtil, usedFor int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := ctx.Cookie(constants.AUTH_COOKIE_KEY)
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

		ctx.Set(constants.CTX_AUTH_PAYLOAD_KEY, payload)
		ctx.Set(constants.CTX_AUTH_TOKEN_KEY, accessToken)
		ctx.Next()
	}
}

// func WSAuthenticationMiddleware(tokenUtil *utils.JWTUtil, usedFor int) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		token := ctx.Param("token")
// 		if len(token) == 0 {
// 			ctx.Error(customerrors.NewError(
// 				"unauthorized",
// 				errors.New("no token found"),
// 				customerrors.Unauthenticate,
// 			))
// 			ctx.Abort()
// 			return
// 		}
// 		payload, err := tokenUtil.VerifyJWT(token, usedFor)
// 		if err != nil {
// 			ctx.Error(err)
// 			ctx.Abort()
// 			return
// 		}

// 		ctx.Set(constants.CTX_AUTH_PAYLOAD_KEY, payload)
// 		ctx.Set(constants.CTX_AUTH_TOKEN_KEY, token)
// 		ctx.Next()
// 	}
// }
