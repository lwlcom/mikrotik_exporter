package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/lwlcom/mikrotik_exporter/config"
	"github.com/lwlcom/mikrotik_exporter/connector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const version string = "0.1"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9436", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	configFile    = flag.String("config-file", "config.yml", "Path to config file")
	debug         = flag.Bool("debug", false, "Show verbose debug output in log")
	cfg           *config.Config
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: mikrotik_exporter [ ... ]\n\nParameters:")
		fmt.Println()
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	c, err := loadConfig()
	if err != nil {
		log.Fatalf("could not load config file. %v", err)
	}
	cfg = c

	if cfg.Identity != "" {
		err = connector.ReadIdentity(cfg.Identity)
		if err != nil {
			log.Fatalf("could not load identity file. %v", err)
		}
	}

	startServer()
}

func printVersion() {
	fmt.Println("mikrotik_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Martin Poppen")
	fmt.Println("Metric exporter for mikrotik devices")
}

func loadConfig() (*config.Config, error) {
	log.Infoln("Loading config from", *configFile)
	b, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	return config.Load(bytes.NewReader(b))
}

func startServer() {
	log.Infof("Starting mikrotik exporter (Version: %s)\n", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>mikrotik Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>mikrotik Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			<h2>More information:</h2>
			<p><a href="https://github.com/lwlcom/mikrotik_exporter">github.com/lwlcom/mikrotik_exporter</a></p>
			</body>
			</html>`))
	})
	http.HandleFunc(*metricsPath, handleMetricsRequest)

	log.Infof("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	reg := prometheus.NewRegistry()

	c := newMikrotikCollector(cfg)
	reg.MustRegister(c)

	l := log.New()
	l.Level = log.ErrorLevel

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorLog:      l,
		ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(w, r)
}
