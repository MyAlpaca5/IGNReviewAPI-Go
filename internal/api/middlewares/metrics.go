package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

func MonitorMetrics(registry *prometheus.Registry) gin.HandlerFunc {
	// create monitoring
	version := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "version",
		Help: "Version information about this binary",
		ConstLabels: map[string]string{
			"version": viper.GetString("general.version"),
		},
	})

	httpRequestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"method", "path"})

	httpRequestDuration := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "http_request_duration_seconds",
		Help:       "Duration of all HTTP requests",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.00: 0.001},
	}, []string{"method", "path"})

	// register monitoring
	registry.MustRegister(version)
	registry.MustRegister(httpRequestsTotal)
	registry.MustRegister(httpRequestDuration)

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Since(start).Milliseconds()
		httpRequestDuration.With(prometheus.Labels{"method": c.Request.Method, "path": c.Request.URL.Path}).Observe(float64(elapsed))
		httpRequestsTotal.With(prometheus.Labels{"method": c.Request.Method, "path": c.Request.URL.Path}).Inc()
	}
}
