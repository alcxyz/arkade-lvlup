package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config is the structure of our yaml configuration file
type Config struct {
	Tools []string `yaml:"tools"`
}

// containsElement checks if a slice contains a given element.
func containsElement(slice []string, elem string) bool {
	for _, item := range slice {
		if item == elem {
			return true
		}
	}
	return false
}

// populateArray splits a space-separated string into a slice of strings.
func populateArray(str string) []string {
	return strings.Fields(str)
}

// readConfig reads the YAML configuration file and unmarshals it into a Config struct.
func readConfig(filePath string) (Config, error) {
	data, err := ioutil.ReadFile(filePath)
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

// writeConfig marshals a Config struct and writes it to a YAML file.
func writeConfig(filePath string, config Config) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}

// forceSyncTools forcefully syncs tools with the configuration.
func forceSyncTools(configFilePath string) {
	fmt.Println("Synchronizing tools...")

	config, err := readConfig(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
		return
	}

	setupWithArkadeIdempotent(configFilePath, config.Tools)
	uninstallExtraneousTools(configFilePath)
}

// setupWithArkadeIdempotent checks if the tool exists in the config,
// if not, adds it and installs/reinstalls it using arkade.
func setupWithArkadeIdempotent(configFilePath string, toolsToProcess []string) {
	for _, tool := range toolsToProcess {
		tool = strings.TrimSpace(tool)
		if tool == "" {
			continue
		}

		config, err := readConfig(configFilePath)
		if err != nil {
			fmt.Printf("Error reading config: %s\n", err)
			return
		}

		if !containsElement(config.Tools, tool) {
			config.Tools = append(config.Tools, tool)
			err := writeConfig(configFilePath, config)
			if err != nil {
				fmt.Printf("Error writing to config: %s\n", err)
				return
			}
			fmt.Printf("Added %s to the configuration file.\n", tool)
		}

		fmt.Printf("Installing tool: %s...\n", tool)
		cmd := exec.Command("arkade", "get", tool)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error: Failed to install %s via arkade. Moving on to the next tool.\n", tool)
		} else {
			fmt.Printf("Successfully installed %s.\n", tool)
		}
	}
}

// removeWithArkade removes tools from the config and uninstalls them using arkade.
func removeWithArkade(configFilePath string, toolsToRemove []string) {
	for _, tool := range toolsToRemove {
		config, err := readConfig(configFilePath)
		if err != nil {
			fmt.Printf("Error reading config: %s\n", err)
			return
		}

		// Find and remove the tool from the Tools slice
		for i, t := range config.Tools {
			if t == tool {
				config.Tools = append(config.Tools[:i], config.Tools[i+1:]...)
				break
			}
		}

		err = writeConfig(configFilePath, config)
		if err != nil {
			fmt.Printf("Error writing to config: %s\n", err)
			return
		}

		fmt.Printf("Marked %s for removal.\n", tool)
	}

	uninstallExtraneousTools(configFilePath)
}

// uninstallExtraneousTools removes tools that are present in the arkade directory but not in the config.
func uninstallExtraneousTools(configFilePath string) {
	config, err := readConfig(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
		return
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error fetching user details: %s\n", err)
		return
	}
	binDir := filepath.Join(usr.HomeDir, ".arkade", "bin")
	files, err := ioutil.ReadDir(binDir)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		return
	}

	removedTools := 0
	for _, file := range files {
		if !containsElement(config.Tools, file.Name()) {
			fmt.Printf("Found extraneous tool: %s. Removing...\n", file.Name())

			toolPath := filepath.Join(binDir, file.Name())
			err := os.Remove(toolPath)
			if err != nil {
				fmt.Printf("Failed to remove tool %s: %s\n", file.Name(), err)
				continue
			} else {
				removedTools++
			}
		}
	}

	if removedTools == 0 {
		fmt.Println("No extraneous tools found. Everything is in sync!")
	}
}

// initializeConfigIfNotExists checks if the config exists. If not, it initializes a new config with tools found in the arkade directory.
func initializeConfigIfNotExists(filePath string) ([]string, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("Config file not found. Initializing a new configuration file...")

		usr, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("Error fetching user details: %s", err)
		}
		binDir := filepath.Join(usr.HomeDir, ".arkade", "bin")
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

		config := Config{
			Tools: tools,
		}

		err = writeConfig(filePath, config)
		if err != nil {
			return nil, err
		}

		fmt.Println("Configuration file has been created at:", filePath)
		fmt.Println("Tools added to the configuration:", strings.Join(tools, ", "))

		return tools, nil
	}
	return nil, err
}

func main() {
	// Parsing command-line flags
	syncFlag := flag.Bool("sync", false, "Sync tools based on configuration.")
	forceFlag := flag.Bool("f", false, "Force sync (only valid with -sync).")
	getFlag := flag.String("get", "", "Install or reinstall specified tools.")
	removeFlag := flag.String("remove", "", "Remove specified tools.")
	flag.Parse()

	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error getting the current user: %s\n", err)
		return
	}
	configFilePath := filepath.Join(usr.HomeDir, ".arkade", "lvlup.yaml")

	// Validate flags
	if *forceFlag && !*syncFlag {
		fmt.Println("Error: -f or --force can only be used with -sync")
		return
	}

	// Initialize the config if it does not exist
	tools, err := initializeConfigIfNotExists(configFilePath)
	if err != nil {
		fmt.Printf("Error initializing config: %s\n", err)
		return
	}

	if tools != nil {
		fmt.Printf("Discovered tools: %s\n", strings.Join(tools, ", "))
	}

	// Handle flags
	if *syncFlag {
		if *forceFlag {
			forceSyncTools(configFilePath)
		} else {
			setupWithArkadeIdempotent(configFilePath, populateArray(*getFlag))
		}
	}

	if *getFlag != "" {
		tools := populateArray(*getFlag)
		setupWithArkadeIdempotent(configFilePath, tools)
	}

	if *removeFlag != "" {
		tools := populateArray(*removeFlag)
		removeWithArkade(configFilePath, tools)
	}
}
