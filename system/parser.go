package system

import (
	"errors"
	"regexp"
)

// Parse parses cli output and tries to find the version number
func (c *systemCollector) Parse(output string) (Resource, error) {
	versionRegexp := regexp.MustCompile(`version: (\S*)`)

	match := versionRegexp.FindStringSubmatch(output)
	if len(match) > 0 {
		return Resource{
			Version: match[1],
		}, nil
	}
	return Resource{}, errors.New("version not found")
}
