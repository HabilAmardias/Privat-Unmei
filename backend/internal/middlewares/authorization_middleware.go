package middlewares

import (
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"

	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware(permission, resource int, rbacr *repositories.RBACRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayloadCtx := ctx.Value(constants.CTX_AUTH_PAYLOAD_KEY)

		if authPayloadCtx == nil {
			ctx.Error(
				customerrors.NewError(
					"cannot identified user",
					errors.New("cannot get auth payload from auth payload context"),
					customerrors.Unauthenticate,
				))
			ctx.Abort()
			return
		}
		authPayload := authPayloadCtx.(*entity.CustomClaim)
		role := authPayload.Role
		rbac := new(entity.Rbac)
		checkRoleErr := rbacr.CheckRoleAccess(ctx, rbac, role, permission, resource)
		if checkRoleErr != nil {
			ctx.Error(checkRoleErr)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
