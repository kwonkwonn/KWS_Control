package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	VmInternalSubnets []string `yaml:"vm_internal_subnets"`
	Cores             []string `yaml:"cores"`
}

func ReadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}
