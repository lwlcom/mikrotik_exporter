package system

import (
	"github.com/lwlcom/mikrotik_exporter/collector"
	"github.com/lwlcom/mikrotik_exporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "mikrotik_system_"

var (
	versionDesc *prometheus.Desc
	uptimeDesc  *prometheus.Desc
	cpuLoadDesc *prometheus.Desc

	voltageDesc          *prometheus.Desc
	currentDesc          *prometheus.Desc
	temperatureDesc      *prometheus.Desc
	cpuTemperaturDesc    *prometheus.Desc
	powerConsumptionDesc *prometheus.Desc
)

func init() {
	l := []string{"target", "name"}
	versionDesc = prometheus.NewDesc(prefix+"version", "Current running version", append(l, "version"), nil)
	uptimeDesc = prometheus.NewDesc(prefix+"uptime_seconds", "Seconds since boot", l, nil)
	cpuLoadDesc = prometheus.NewDesc(prefix+"cpu_load_percent", "Current CPU load in percent", l, nil)

	voltageDesc = prometheus.NewDesc(prefix+"voltage", "Supplied voltage", l, nil)
	currentDesc = prometheus.NewDesc(prefix+"current", "Supplied current in mA", l, nil)
	temperatureDesc = prometheus.NewDesc(prefix+"temp", "Current temperature in celsius", l, nil)
	cpuTemperaturDesc = prometheus.NewDesc(prefix+"cpu_temp", "Current CPU temperature in celsius", l, nil)
	powerConsumptionDesc = prometheus.NewDesc(prefix+"power_consumption", "Current power consumption in watt", l, nil)
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
	ch <- uptimeDesc
	ch <- cpuLoadDesc

	ch <- voltageDesc
	ch <- currentDesc
	ch <- temperatureDesc
	ch <- cpuTemperaturDesc
	ch <- powerConsumptionDesc
}

func (c *systemCollector) CollectResource(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	out, err := client.RunCommand("/system resource print")
	if err != nil {
		return err
	}
	resource, err := c.ParseResource(out)
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(versionDesc, prometheus.GaugeValue, 1, append(labelValues, resource.Version)...)
	ch <- prometheus.MustNewConstMetric(uptimeDesc, prometheus.GaugeValue, resource.Uptime, labelValues...)
	ch <- prometheus.MustNewConstMetric(cpuLoadDesc, prometheus.GaugeValue, resource.CPULoad, labelValues...)
	return nil
}

func (c *systemCollector) CollectHealth(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	out, err := client.RunCommand("/system health print")
	if err != nil {
		return err
	}
	health, err := c.ParseHealth(out)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(voltageDesc, prometheus.GaugeValue, health.Voltage, labelValues...)
	ch <- prometheus.MustNewConstMetric(currentDesc, prometheus.GaugeValue, health.Current, labelValues...)
	ch <- prometheus.MustNewConstMetric(temperatureDesc, prometheus.GaugeValue, health.Temperature, labelValues...)
	ch <- prometheus.MustNewConstMetric(cpuTemperaturDesc, prometheus.GaugeValue, health.CPUTemperature, labelValues...)
	ch <- prometheus.MustNewConstMetric(powerConsumptionDesc, prometheus.GaugeValue, health.PowerConsumption, labelValues...)
	return nil
}

// Collect collects metrics from mikrotik
func (c *systemCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	c.CollectResource(client, ch, labelValues)
	c.CollectHealth(client, ch, labelValues)

	return nil
}
