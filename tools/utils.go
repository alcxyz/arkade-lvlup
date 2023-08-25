package tools

import (
	"arkade-lvlup/config"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// GetBinDir retrieves the path to the arkade/bin directory for the current user.
func GetBinDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".arkade", "bin"), nil
}

// GetConfigDir retrieves the path to the arkade directory for the configuration.
func GetConfigDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".arkade"), nil
}

// ListToolsInBinDir lists all tools present in the arkade/bin directory.
func ListToolsInBinDir() ([]string, error) {
	binDir, err := GetBinDir()
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(binDir)
	if err != nil {
		return nil, err
	}

	var tools []string
	for _, file := range files {
		if !file.IsDir() {
			tools = append(tools, file.Name())
		}
	}

	return tools, nil
}

// ContainsElement checks if a slice contains a given element.
func ContainsElement(slice []string, elem string) bool {
	for _, item := range slice {
		if item == elem {
			return true
		}
	}
	return false
}

// PopulateArray splits a comma-separated string into a slice of strings.
func PopulateArray(input string) []string {
	return strings.Split(strings.TrimSpace(input), ",")
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

// configExists checks if the configuration file exists.
func configExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
