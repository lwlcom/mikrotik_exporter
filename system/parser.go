package system

import (
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

// Parse parses cli output and parses system resource data
func (c *systemCollector) ParseResource(output string) (Resource, error) {
	resource := Resource{}

	versionRegexp := regexp.MustCompile(`version: (\S*)`)
	match := versionRegexp.FindStringSubmatch(output)
	if len(match) > 0 {
		resource.Version = match[1]
	}

	uptimeRegexp := regexp.MustCompile(`uptime: (\S*)`)
	match = uptimeRegexp.FindStringSubmatch(output)
	if len(match) > 0 {
		resource.Uptime = uptimeToSeconds(match[1])
	}

	cpuLoadRegexp := regexp.MustCompile(`cpu-load: (\d*)%`)
	match = cpuLoadRegexp.FindStringSubmatch(output)
	if len(match) < 0 {
		i, err := strconv.Atoi(match[1])
		if err == nil {
			resource.CPULoad = float64(i)
		}
	}
	return resource, nil
}

// Parse parses cli output and parses system health data
func (c *systemCollector) ParseHealth(output string) (Health, error) {
	health := Health{}

	voltageRegexp := regexp.MustCompile(`voltage: ([\d\.]*)V`)
	currentRegexp := regexp.MustCompile(`current: (\d*)mA`)
	tempRegexp := regexp.MustCompile(`temperature: (\d*)C`)
	cpuTempRegexp := regexp.MustCompile(`cpu-temperature: (\d*)C`)
	powerRegexp := regexp.MustCompile(`power-consumption: ([\d\.]*)W`)

	match := voltageRegexp.FindStringSubmatch(output)
	if len(match) > 0 {
		if i, err := strconv.ParseFloat(match[1], 64); err == nil {
			health.Voltage = i
		}
	}

	match = currentRegexp.FindStringSubmatch(output)
	if len(match) > 0 {
		if i, err := strconv.ParseFloat(match[1], 64); err == nil {
			health.Current = i
		}
	}

	match = tempRegexp.FindStringSubmatch(output)
	if len(match) > 0 {
		if i, err := strconv.ParseFloat(match[1], 64); err == nil {
			health.Temperature = i
		}
	}

	match = cpuTempRegexp.FindStringSubmatch(output)
	if len(match) > 0 {
		if i, err := strconv.ParseFloat(match[1], 64); err == nil {
			health.CPUTemperature = i
		}
	}

	match = powerRegexp.FindStringSubmatch(output)
	if len(match) > 0 {
		if i, err := strconv.ParseFloat(match[1], 64); err == nil {
			health.PowerConsumption = i
		}
	}

	return health, nil
}
