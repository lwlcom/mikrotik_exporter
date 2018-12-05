package interfaces

import (
	"regexp"

	"github.com/lwlcom/mikrotik_exporter/util"
)

// Parse parses cli output and tries to find interfaces with related traffic stats
func (c *interfaceCollector) Parse(output string) ([]Interface, error) {
	items := []Interface{}
	headerRegexp := regexp.MustCompile(`(Flags.*)`)
	newlinesRegexp := regexp.MustCompile(`\r\n\s+\n`)
	ethernetRegexp := regexp.MustCompile(`(?ms)^\s*\d+\s+[DXRS]*(\s+;;; (.*)$)?\s+name="(.*)".*\s+rx-byte=([\d ]+).*\s+tx-byte=([\d ]+).*\s+rx-packet=([\d ]+).*\s+tx-packet=([\d ]+).*\s+rx-drop=([\d ]+).*\s+tx-drop=([\d ]+).*\s+rx-error=([\d ]+).*\s+tx-error=([\d ]+).*\s+`)

	output = headerRegexp.ReplaceAllString(output, "")
	interfaces := newlinesRegexp.Split(output, -1)

	for _, data := range interfaces {
		if match := ethernetRegexp.FindStringSubmatch(data); match != nil {
			item := Interface{
				Comment:  util.Normalize(match[2]),
				Name:     match[3],
				RxByte:   util.Str2float64(match[4]),
				TxByte:   util.Str2float64(match[5]),
				RxPacket: util.Str2float64(match[6]),
				TxPacket: util.Str2float64(match[7]),
				RxDrop:   util.Str2float64(match[8]),
				TxDrop:   util.Str2float64(match[9]),
				RxError:  util.Str2float64(match[10]),
				TxError:  util.Str2float64(match[11]),
			}
			items = append(items, item)
		}
	}
	return items, nil
}
