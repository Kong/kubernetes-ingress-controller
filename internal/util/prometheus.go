package util

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	proxyProcessers = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "proxyconfig_failures_total",
			Help: "Number of configuration proccessed w/ failure.",
		},
	)
	proxyProcesserFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "proxyconfig_success_total",
			Help: "Number of successful configuration processed",
		},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(proxyProcessers, proxyProcesserFailures)
}
