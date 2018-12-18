package optics

import (
	"errors"
	"regexp"
	"strings"

	"github.com/lwlcom/mikrotik_exporter/util"
)

// ParseInterfaces parses cli output and tries to find interface names
func (c *opticsCollector) ParseInterfaces(output string) ([]string, error) {
	items := []string{}
	ifNameRegexp := regexp.MustCompile(`name=(.*) default-name`)

	matches := ifNameRegexp.FindAllStringSubmatch(output, -1)
	for _, match := range matches {
		items = append(items, match[1])
	}
	return items, nil
}

// ParseSfp parses cli output and tries to find optical power data
func (c *opticsCollector) ParseSfp(output string) (Optic, error) {
	noLinkRegexp := regexp.MustCompile(`status: no-link`)
	if noLinkRegexp.MatchString(output) {
		return Optic{}, errors.New("no-link")
	}
	if !strings.Contains(output, "sfp") {
		return Optic{}, errors.New("not an sfp port")
	}

	powerRegexp := regexp.MustCompile(`(?ms)sfp-tx-power: (.+)dBm.*sfp-rx-power: (.+)dBm`)
	match := powerRegexp.FindStringSubmatch(output)

	if match != nil {
		return Optic{
			TxPower: util.Str2float64(match[1]),
			RxPower: util.Str2float64(match[2]),
		}, nil
	}
	return Optic{}, errors.New("no-power")
}
