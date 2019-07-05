# mikrotik_exporter
Exporter for metrics from mikrotik devices (via SSH) https://prometheus.io/

###### example config
```yaml
targets:
  - name: router1510
    address: 192.168.0.1
    user: prom
    password: topsecret

features:
  optics: false
  system: false
  dhcp: false
```

# flags
Name     | Description | Default
---------|-------------|---------
version | Print version information. |
web.listen-address | Address on which to expose metrics and web interface. | :9362
web.telemetry-path | Path under which to expose metrics. | /metrics
config-file | Path to config file |
debug | Show verbose debug output | false


## Install
```bash
go get -u github.com/lwlcom/mikrotik_exporter
```

## Usage
```bash
./mikrotik_exporter -config-file config.yml
```

## Third Party Components
This software uses components of the following projects
* Prometheus Go client library (https://github.com/prometheus/client_golang)

## License
(c) Martin Poppen, 2018. Licensed under [MIT](LICENSE) license.

## Prometheus
see https://prometheus.io/
