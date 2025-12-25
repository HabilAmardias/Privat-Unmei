package middlewares

import (
	"context"
	"errors"
	"os"
	"privat-unmei/internal/customerrors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

func RateLimiterMiddleware(limiter *redis_rate.Limiter) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		requestRate, ok := os.LookupEnv("REQUEST_RATE")
		if !ok {
			ctx.Error(customerrors.NewError(
				"something went wrong",
				errors.New("request rate config does not exist"),
				customerrors.CommonErr,
			))
			ctx.Abort()
			return
		}
		requestRateParsed, err := strconv.Atoi(requestRate)
		if err != nil {
			ctx.Error(customerrors.NewError(
				"something went wrong",
				err,
				customerrors.CommonErr,
			))
			ctx.Abort()
			return
		}
		res, err := limiter.Allow(context.Background(), ctx.ClientIP(), redis_rate.PerMinute(requestRateParsed))
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
