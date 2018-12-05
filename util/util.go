package util

import (
	"io/ioutil"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// Str2float64 converts a string to float64
func Str2float64(str string) float64 {
	value, err := strconv.ParseFloat(strings.Replace(str, " ", "", -1), 64)
	if err != nil {
		return -1
	}
	return value
}

// Normalize converts a string to valid utf8
func Normalize(str string) string {
	if len(str) <= 1 {
		return ""
	}
	str = strings.TrimSuffix(str, "\r")
	inUTF8 := transform.NewReader(strings.NewReader(str), charmap.Windows1252.NewDecoder())
	decBytes, _ := ioutil.ReadAll(inUTF8)
	return string(decBytes)
}
