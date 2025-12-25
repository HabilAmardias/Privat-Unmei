package middlewares

import (
	"context"
	"errors"
	"privat-unmei/internal/customerrors"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

func RateLimiterMiddleware(limiter *redis_rate.Limiter) gin.HandlerFunc {
	requestRate := 10
	return func(ctx *gin.Context) {
		res, err := limiter.Allow(context.Background(), ctx.ClientIP(), redis_rate.PerMinute(requestRate))
		if err != nil {
			ctx.Error(customerrors.NewError(
				"something went wrong",
				err,
				customerrors.CommonErr,
			))
			ctx.Abort()
			return
		}
		if res.Remaining <= 0 {
			ctx.Error(customerrors.NewError(
				"too many requests going, please try again",
				errors.New("too many request"),
				customerrors.TooManyRequest,
			))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
