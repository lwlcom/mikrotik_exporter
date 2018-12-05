package interfaces

import (
	"github.com/lwlcom/mikrotik_exporter/collector"
	"github.com/lwlcom/mikrotik_exporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "mikrotik_interface_"

var (
	receiveBytesDesc    *prometheus.Desc
	receivePacketsDesc  *prometheus.Desc
	receiveErrorsDesc   *prometheus.Desc
	receiveDropsDesc    *prometheus.Desc
	transmitBytesDesc   *prometheus.Desc
	transmitPacketsDesc *prometheus.Desc
	transmitErrorsDesc  *prometheus.Desc
	transmitDropsDesc   *prometheus.Desc
)

func init() {
	l := []string{"target", "name", "interface", "comment"}
	receiveBytesDesc = prometheus.NewDesc(prefix+"receive_bytes", "Received data in bytes", l, nil)
	receivePacketsDesc = prometheus.NewDesc(prefix+"receive_packets_total", "Received packets", l, nil)
	receiveErrorsDesc = prometheus.NewDesc(prefix+"receive_errors", "Number of errors caused by incoming packets", l, nil)
	receiveDropsDesc = prometheus.NewDesc(prefix+"receive_drops", "Number of dropped incoming packets", l, nil)
	transmitBytesDesc = prometheus.NewDesc(prefix+"transmit_bytes", "Transmitted packets", l, nil)
	transmitPacketsDesc = prometheus.NewDesc(prefix+"transmit_packets_total", "Transmitted data in bytes", l, nil)
	transmitErrorsDesc = prometheus.NewDesc(prefix+"transmit_errors", "Number of errors caused by outgoing packets", l, nil)
	transmitDropsDesc = prometheus.NewDesc(prefix+"transmit_drops", "Number of dropped outgoing packets", l, nil)
}

// Collector collects interface metrics
type interfaceCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &interfaceCollector{}
}

// Describe describes the metrics
func (*interfaceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- receiveBytesDesc
	ch <- receivePacketsDesc
	ch <- receiveErrorsDesc
	ch <- receiveDropsDesc
	ch <- transmitBytesDesc
	ch <- transmitPacketsDesc
	ch <- transmitErrorsDesc
	ch <- transmitDropsDesc
}

// Collect collects metrics from mikrotik
func (c *interfaceCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	out, err := client.RunCommand("/interface print stats-detail")
	if err != nil {
		return err
	}
	items, err := c.Parse(out)

	for _, item := range items {
		l := append(labelValues, item.Name, item.Comment)

		ch <- prometheus.MustNewConstMetric(receiveBytesDesc, prometheus.GaugeValue, item.RxByte, l...)
		ch <- prometheus.MustNewConstMetric(receivePacketsDesc, prometheus.GaugeValue, item.TxPacket, l...)
		ch <- prometheus.MustNewConstMetric(receiveErrorsDesc, prometheus.GaugeValue, item.RxError, l...)
		ch <- prometheus.MustNewConstMetric(receiveDropsDesc, prometheus.GaugeValue, item.RxDrop, l...)
		ch <- prometheus.MustNewConstMetric(transmitBytesDesc, prometheus.GaugeValue, item.TxByte, l...)
		ch <- prometheus.MustNewConstMetric(transmitPacketsDesc, prometheus.GaugeValue, item.TxPacket, l...)
		ch <- prometheus.MustNewConstMetric(transmitErrorsDesc, prometheus.GaugeValue, item.TxError, l...)
		ch <- prometheus.MustNewConstMetric(transmitDropsDesc, prometheus.GaugeValue, item.TxDrop, l...)
	}

	return nil
}
