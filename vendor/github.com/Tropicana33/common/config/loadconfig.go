package config

import (
	"os"
)

func LoadConfigByFile(file string) (*Config, error) {
	fl, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	c := NewConfig()

	err = c.UpdateConfig(fl)
	fl.Close()
	if err != nil {
		return nil, err
	}

	return c, nil
}
