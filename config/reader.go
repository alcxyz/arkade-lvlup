package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// ReadConfig reads the YAML configuration file and unmarshals it into a ToolConfig struct.
func ReadConfig(filePath string) (ToolConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ToolConfig{}, err
	}

	var config ToolConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return ToolConfig{}, err
	}

	return config, nil
}

// WriteConfig marshals a ToolConfig struct and writes it to a YAML file.
func WriteConfig(filePath string, config ToolConfig) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
