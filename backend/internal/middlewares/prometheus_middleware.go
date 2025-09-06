package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware(httpRequestsTotal *prometheus.CounterVec, httpRequestDuration *prometheus.HistogramVec) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		duration := time.Since(start).Seconds()

		status := ctx.Writer.Status()
		httpRequestsTotal.WithLabelValues(ctx.Request.Method, ctx.FullPath(),
			http.StatusText(status)).Inc()
		httpRequestDuration.WithLabelValues(ctx.Request.Method, ctx.FullPath()).Observe(duration)

	}
}
