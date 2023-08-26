package tools

import (
	"arkade-lvlup/config"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// GetConfigDir retrieves the path to the arkade directory for the configuration.
func GetConfigDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".arkade"), nil
}

func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", fmt.Errorf("Error getting config directory: %s", err)
	}
	return filepath.Join(configDir, "lvlup.yaml"), nil
}

// configExists checks if the configuration file exists.
func configExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// InitializeConfigIfNotExists checks if the config exists. If not, it initializes a new config with tools found in the arkade directory.
func InitializeConfigIfNotExists(filePath string) ([]string, error) {
	if configExists(filePath) {
		return nil, nil
	}

	fmt.Println("Config file not found. Initializing a new configuration file...")

	tools, err := ListToolsInBinDir()
	if err != nil {
		return nil, err
	}

	cfg := config.ArkadeTools{Tools: tools} // Updated this line

	err = config.WriteConfig(filePath, cfg)
	if err != nil {
		return nil, err
	}

	fmt.Println("Configuration file has been created at:", filePath)
	fmt.Println("Tools added to the configuration:", strings.Join(tools, ", "))

	return tools, nil
}
