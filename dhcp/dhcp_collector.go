package dhcp

import (
	"github.com/lwlcom/mikrotik_exporter/collector"
	"github.com/lwlcom/mikrotik_exporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "mikrotik_dhcp_"

var (
	leaseCountDesc *prometheus.Desc
)

func init() {
	l := []string{"target", "name", "server"}
	leaseCountDesc = prometheus.NewDesc(prefix+"lease_count", "Current number of leases per server", l, nil)
}

// Collector collects dhcp metrics
type dhcpCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &dhcpCollector{}
}

// Describe describes the metrics
func (*dhcpCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- leaseCountDesc
}

// Collect collects metrics from mikrotik
func (c *dhcpCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	out, err := client.RunCommand("/ip dhcp-server lease print detail")
	if err != nil {
		return err
	}
	items, err := c.Parse(out)

	for server, count := range items {
		l := append(labelValues, server)

		ch <- prometheus.MustNewConstMetric(leaseCountDesc, prometheus.GaugeValue, float64(count), l...)
	}

	return nil
}
