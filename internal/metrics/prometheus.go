package prometheusmetrics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: prometheus.DefBuckets, // Use default buckets
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(requestDuration)
}

func RegisterMetrics(reg prometheus.Registerer) {
	// Register other metrics
	registerCPUMetrics(reg)
	registerHDMetrics(reg)
}

func registerCPUMetrics(reg prometheus.Registerer) {
	cpuTemp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})
	reg.MustRegister(cpuTemp)
}

func registerHDMetrics(reg prometheus.Registerer) {
	hdFailures := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk goAuthException.",
		},
		[]string{"device"},
	)
	reg.MustRegister(hdFailures)
}

// Middleware for tracking request duration
func InstrumentHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := fmt.Sprintf("%d", c.Writer.Status())
		method := c.Request.Method
		path := c.FullPath()

		requestDuration.WithLabelValues(method, path, status).Observe(duration)
	}
}

// Handler for exposing Prometheus metrics
func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
