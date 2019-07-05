package dhcp

import (
	"regexp"
	"strings"
)

// Parse parses cli output and parses leases
func (c *dhcpCollector) Parse(output string) (map[string]int, error) {
	items := make(map[string]int)
	leaseRegexp := regexp.MustCompile(`active-server=(.*?)(\r\n| host-name)`)
	matches := leaseRegexp.FindAllStringSubmatch(output, -1)
	for _, match := range matches {
		server := strings.Trim(match[1], " ")
		if _, ok := items[server]; ok {
			items[server] = items[server] + 1
		} else {
			items[server] = 1
		}
	}
	return items, nil
}
