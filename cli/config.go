package cli

import (
	"os"

	"github.com/BrunoQuaresma/openwritter/pkg/qwriter"
	"gopkg.in/yaml.v3"
)

type Config struct {
	path     string
	Profiles []qwriter.Profile `yaml:"profiles"`
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
