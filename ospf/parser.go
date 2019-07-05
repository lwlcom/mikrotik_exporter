package ospf

import (
	"regexp"
	"strconv"
)

// Parse parses cli output and parses leases
func (c *ospfCollector) Parse(output string) ([]Ospf, error) {
	items := []Ospf{}

	leaseRegexp := regexp.MustCompile(`address=([\d\.]+)\s+state="(\w+)" state-changes=(\d+)`)
	matches := leaseRegexp.FindAllStringSubmatch(output, -1)

	for _, match := range matches {
		i, err := strconv.Atoi(match[3])
		if err != nil {
			i = 0
		}

		item := Ospf{
			Address:      match[1],
			State:        match[2],
			StateChanges: float64(i),
		}
		items = append(items, item)
	}
	return items, nil
}
