package server

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_requests_seconds_count",
			Help: "Count of HTTP requests processed, labeled by status code, method, and path.",
		},
		[]string{"status", "method", "path"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_server_requests_duration_seconds",
			Help:    "Histogram of response latency (seconds) of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status", "method", "path"},
	)
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		statusString := strconv.Itoa(status)

		path := c.FullPath()
		method := c.Request.Method

		// 메트릭 업데이트
		httpRequestsTotal.WithLabelValues(statusString, method, path).Inc()
		httpRequestDuration.WithLabelValues(statusString, method, path).Observe(duration)
	}
}

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}
