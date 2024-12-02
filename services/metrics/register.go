package metrics

import "github.com/prometheus/client_golang/prometheus"

func NewRegister(values []prometheus.Collector) *prometheus.Registry {
	reg := prometheus.NewRegistry()
	for _, value := range values {
		reg.MustRegister(value)
	}
	return reg
}
