package collector

import (
	"github.com/lwlcom/mikrotik_exporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

// RPCCollector collects metrics from mikrotik using rpc.Client
type RPCCollector interface {

	// Describe describes the metrics
	Describe(ch chan<- *prometheus.Desc)

	// Collect collects metrics from mikrotik
	Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error
}
