package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware(httpRequestsTotal *prometheus.CounterVec, httpRequestDuration *prometheus.HistogramVec) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(c.Writer.Status())
		endpoint := c.Request.URL.Path

		httpRequestsTotal.WithLabelValues(c.Request.Method, endpoint, statusCode).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, endpoint, statusCode).Observe(duration)
	}
}
