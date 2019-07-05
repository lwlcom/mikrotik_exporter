package ospf

import (
	"strings"

	"github.com/lwlcom/mikrotik_exporter/collector"
	"github.com/lwlcom/mikrotik_exporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "mikrotik_ospf_"

var (
	stateChangeCountDesc *prometheus.Desc
	neighborStateDesc    *prometheus.Desc
)

func init() {
	l := []string{"target", "name", "address"}
	stateChangeCountDesc = prometheus.NewDesc(prefix+"state_change_count", "Number of state changes", l, nil)
	l = append(l, "state")
	neighborStateDesc = prometheus.NewDesc(prefix+"neighbor_state", "Current neighbor state (0=down,1=attempt,2=init,3=2-way,4=ExStart,5=Exchange,6=Loading,7=full)", l, nil)
}

// Collector collects ospf metrics
type ospfCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &ospfCollector{}
}

// Describe describes the metrics
func (*ospfCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- stateChangeCountDesc
	ch <- neighborStateDesc
}

func valueForState(state string) float64 {
	switch strings.ToLower(state) {
	case "attempt":
		return 1
	case "init":
		return 2
	case "2-way":
		return 3
	case "exstart":
		return 4
	case "exchange":
		return 5
	case "loading":
		return 6
	case "full":
		return 7
	default:
		return 0
	}
}

// Collect collects metrics from mikrotik
func (c *ospfCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	out, err := client.RunCommand("/routing ospf neighbor print")
	if err != nil {
		return err
	}
	items, err := c.Parse(out)

	for _, item := range items {
		l := append(labelValues, item.Address)
		ch <- prometheus.MustNewConstMetric(stateChangeCountDesc, prometheus.GaugeValue, item.StateChanges, l...)
		l = append(l, item.State)
		ch <- prometheus.MustNewConstMetric(neighborStateDesc, prometheus.GaugeValue, valueForState(item.State), l...)
	}

	return nil
}
