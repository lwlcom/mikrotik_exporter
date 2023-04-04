package optics

import (
	"github.com/lwlcom/mikrotik_exporter/collector"
	"github.com/lwlcom/mikrotik_exporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "mikrotik_optics_"

var (
	receivePowerDesc  *prometheus.Desc
	transmitPowerDesc *prometheus.Desc
)

func init() {
	l := []string{"target", "name", "interface"}
	receivePowerDesc = prometheus.NewDesc(prefix+"rx_power", "Transceiver Rx power", l, nil)
	transmitPowerDesc = prometheus.NewDesc(prefix+"tx_power", "Transceiver Tx power", l, nil)
}

// Collector collects interface metrics
type opticsCollector struct {
       monitorOpticsWithNoLink	bool
}

// NewCollector creates a new collector
func NewCollector(monitorOpticsWithNoLink bool) collector.RPCCollector {
	return &opticsCollector{monitorOpticsWithNoLink}
}

// Describe describes the metrics
func (*opticsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- receivePowerDesc
	ch <- transmitPowerDesc
}

// Collect collects metrics from mikrotik
func (c *opticsCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	out, err := client.RunCommand("/interface ethernet print terse")
	if err != nil {
		return err
	}
	sfps, err := c.ParseInterfaces(out)
	if err != nil {
		return err
	}
	for _, sfp := range sfps {
		out, err := client.RunCommand("/interface ethernet monitor " + sfp + " once")
		if err != nil {
			return err
		}
		item, err := c.ParseSfp(out)
		if err != nil {
			continue
		}
		l := append(labelValues, sfp)
		ch <- prometheus.MustNewConstMetric(receivePowerDesc, prometheus.GaugeValue, item.RxPower, l...)
		ch <- prometheus.MustNewConstMetric(transmitPowerDesc, prometheus.GaugeValue, item.TxPower, l...)
	}
	return nil
}
