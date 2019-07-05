package system

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func uptimeToSeconds(time string) float64 {
	seconds := 0
	scale := map[string]int{
		"w": 604800,
		"d": 86400,
		"h": 3600,
		"m": 60,
		"s": 1,
	}
	units := []string{"w", "d", "h", "m", "s"}
	for _, unit := range units {
		s := strings.Split(time, unit)
		if len(s) == 1 {
			continue
		}
		i, err := strconv.Atoi(s[0])
		if err != nil {
			fmt.Println(err)
			return 0
		}
		seconds = seconds + (i * scale[unit])
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
