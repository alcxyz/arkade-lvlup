package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// ReadConfig reads the YAML configuration file and unmarshals it into a Config struct.
func ReadConfig(filePath string) (Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// WriteConfig marshals a Config struct and writes it to a YAML file.
func WriteConfig(filePath string, config Config) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}