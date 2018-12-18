package system

import (
	"github.com/lwlcom/mikrotik_exporter/collector"
	"github.com/lwlcom/mikrotik_exporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "mikrotik_system_"

var (
	versionDesc *prometheus.Desc
)

func init() {
	l := []string{"target", "name", "version"}
	versionDesc = prometheus.NewDesc(prefix+"version", "Current running version", l, nil)
}

// Collector collects interface metrics
type systemCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &systemCollector{}
}

// Describe describes the metrics
func (*systemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionDesc
}

// Collect collects metrics from mikrotik
func (c *systemCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	out, err := client.RunCommand("/system resource print")
	if err != nil {
		return err
	}
	resource, err := c.Parse(out)
	if err != nil {
		return err
	}
	l := append(labelValues, resource.Version)
	ch <- prometheus.MustNewConstMetric(versionDesc, prometheus.GaugeValue, 1, l...)
	return nil
}
