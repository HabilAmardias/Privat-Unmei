package middlewares

import (
	"privat-unmei/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(lg logger.CustomLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		elapsed := time.Since(start)
		lg.Infoln(ctx.Request.Method, ctx.Request.URL.Path, ctx.Writer.Status(), elapsed.String())
	}
}
