package config

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the configuration for the exporter
type Config struct {
	Targets  []Target `yaml:"targets"`
	Features struct {
		Optics bool `yaml:"optics,omitempty"`
		System bool `yaml:"system,omitempty"`
		Dhcp   bool `yaml:"dhcp,omitempty"`
	} `yaml:"features"`
}

type Target struct {
	Name     string `yaml:"name,omitempty"`
	Address  string `yaml:"address,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
}

func New() *Config {
	c := &Config{}
	return c
}

// Load loads a config from reader
func Load(reader io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	c := New()
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
