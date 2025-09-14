package middlewares

import (
	"errors"
	"fmt"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/logger"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/repositories/cache"

	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware(permission, resource int, rbacr *repositories.RBACRepository, rbacCache *cache.RBACCacheRepository, lg logger.CustomLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const cacheDuration = 604800 // 1 week
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

		cacheHasAccess, err := rbacCache.CheckRoleAccess(ctx, role, permission, resource)
		if err != nil {
			lg.Errorln(err.Error())
		}
		unauthorizedErr := customerrors.NewError(
			"unauthorize",
			fmt.Errorf("rbac for role: %d, permission: %d, resource: %d is not found in rbac records", role, permission, resource),
			customerrors.Unauthenticate,
		)
		if cacheHasAccess != nil {
			if !*cacheHasAccess {
				ctx.Error(unauthorizedErr)
				ctx.Abort()
				return
			}

			ctx.Next()
			return
		}

		rbac := new(entity.Rbac)
		if err := rbacr.CheckRoleAccess(ctx, rbac, role, permission, resource); err != nil {
			var parsedErr *customerrors.CustomError
			if errors.As(err, &parsedErr) {
				if parsedErr.ErrCode == customerrors.Unauthenticate {
					if err := rbacCache.SetCheckRoleAccess(ctx, role, permission, resource, cacheDuration, false); err != nil {
						lg.Errorln(err.Error())
					}
				}
			}
			ctx.Error(err)
			ctx.Abort()
		}

		if err := rbacCache.SetCheckRoleAccess(ctx, role, permission, resource, cacheDuration, true); err != nil {
			lg.Errorln(err.Error())
		}

		ctx.Next()
	}
}
