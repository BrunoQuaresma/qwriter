package cli

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	path     string
	Profiles []Profile `yaml:"profiles"`
}

type Profile struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func NewConfig(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{path: path}
	err = yaml.Unmarshal(raw, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
