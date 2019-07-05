package main

import (
	"time"

	"sync"

	"github.com/lwlcom/mikrotik_exporter/collector"
	"github.com/lwlcom/mikrotik_exporter/config"
	"github.com/lwlcom/mikrotik_exporter/connector"
	"github.com/lwlcom/mikrotik_exporter/dhcp"
	"github.com/lwlcom/mikrotik_exporter/interfaces"
	"github.com/lwlcom/mikrotik_exporter/optics"
	"github.com/lwlcom/mikrotik_exporter/ospf"
	"github.com/lwlcom/mikrotik_exporter/rpc"
	"github.com/lwlcom/mikrotik_exporter/system"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const prefix = "mikrotik_"

var (
	scrapeDurationDesc *prometheus.Desc
	upDesc             *prometheus.Desc
)

func init() {
	upDesc = prometheus.NewDesc(prefix+"up", "Scrape of target was successful", []string{"target", "name"}, nil)
	scrapeDurationDesc = prometheus.NewDesc(prefix+"collector_duration_seconds", "Duration of a collector scrape for one target", []string{"target", "name"}, nil)
}

type mikrotikCollector struct {
	cfg        *config.Config
	collectors map[string]collector.RPCCollector
}

func newMikrotikCollector(cfg *config.Config) *mikrotikCollector {
	collectors := collectors()
	return &mikrotikCollector{cfg, collectors}
}

func collectors() map[string]collector.RPCCollector {
	m := map[string]collector.RPCCollector{
		"interfaces": interfaces.NewCollector(),
	}

	f := &cfg.Features

	if f.Optics {
		m["optics"] = optics.NewCollector()
	}

	if f.System {
		m["system"] = system.NewCollector()
	}

	if f.Dhcp {
		m["dhcp"] = dhcp.NewCollector()
	}

	if f.Ospf {
		m["ospf"] = ospf.NewCollector()
	}

	return m
}

// Describe implements prometheus.Collector interface
func (c *mikrotikCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upDesc
	ch <- scrapeDurationDesc

	for _, col := range c.collectors {
		col.Describe(ch)
	}
}

// Collect implements prometheus.Collector interface
func (c *mikrotikCollector) Collect(ch chan<- prometheus.Metric) {
	hosts := c.cfg.Targets
	wg := &sync.WaitGroup{}

	wg.Add(len(hosts))
	for _, h := range hosts {
		go c.collectForHost(h.Name, h.Address, h.User, h.Password, ch, wg)
	}

	wg.Wait()
}

func (c *mikrotikCollector) collectForHost(name, address, user, password string, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()

	l := []string{address, name}

	t := time.Now()
	defer func() {
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(t).Seconds(), l...)
	}()

	conn, err := connector.NewSSSHConnection(address, user, password)
	if err != nil {
		log.Errorln(err)
		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 0, l...)
		return
	}
	defer conn.Close()

	ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 1, l...)

	rpc := rpc.NewClient(conn, *debug)
	for k, col := range c.collectors {
		err = col.Collect(rpc, ch, l)
		if err != nil && err.Error() != "EOF" {
			log.Errorln(k + ": " + err.Error())
		}
	}
}
