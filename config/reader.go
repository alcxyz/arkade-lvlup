package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// ReadConfig reads the YAML configuration file and unmarshals it into an ArkadeTools struct.
func ReadConfig(filePath string) (ArkadeTools, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ArkadeTools{}, err
	}

	var config ArkadeTools
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return ArkadeTools{}, err
	}

	return config, nil
}

// WriteConfig marshals an ArkadeTools struct and writes it to a YAML file.
func WriteConfig(filePath string, config ArkadeTools) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
