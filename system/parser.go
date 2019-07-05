package system

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func uptimeToSeconds(time string) float64 {
	seconds := 0
	units := map[string]int{
		"w": 604800,
		"d": 86400,
		"h": 3600,
		"m": 60,
		"s": 1,
	}

	for unit, v := range units {
		s := strings.Split(time, unit)
		if len(s) == 1 {
			continue
		}
		i, err := strconv.Atoi(s[0])
		if err != nil {
			return 0
		}
		seconds = seconds + (i * v)
		time = s[1]
	}
	return float64(seconds)
}

// Parse parses cli output and tries to find the version number
func (c *systemCollector) Parse(output string) (Resource, error) {
	resource := Resource{}

	versionRegexp := regexp.MustCompile(`version: (\S*)`)
	match := versionRegexp.FindStringSubmatch(output)
	if len(match) > 0 {
		resource.Version = match[1]
	} else {
		return Resource{}, errors.New("version not found")
	}

	uptimeRegexp := regexp.MustCompile(`uptime: (\S*)`)
	match = uptimeRegexp.FindStringSubmatch(output)
	if len(match) > 0 {
		resource.Uptime = uptimeToSeconds(match[1])
	} else {
		return Resource{}, errors.New("uptime not found")
	}

	return resource, nil
}
