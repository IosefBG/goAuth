// prometheusmetrics/metrics.go
package prometheusmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// RegisterMetrics registers all Prometheus metrics
func RegisterMetrics(reg prometheus.Registerer) {
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
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)
	reg.MustRegister(hdFailures)
}
