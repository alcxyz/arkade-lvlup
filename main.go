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

var (
	forceFlag       = flag.Bool("f", false, "Force sync (only valid with -sync).")
	passthroughFlag = flag.Bool("passthrough", false, "Show arkade outputs.")
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

// populateArray splits a comma-separated string into a slice of strings.
func populateArray(str string) []string {
	return strings.Split(strings.TrimSpace(str), ",")
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

	// removedTools should be calculated and returned by uninstallExtraneousTools
	removedTools := uninstallExtraneousTools(configFilePath)
	if removedTools == 0 && !*forceFlag {
		fmt.Println("Everything is in sync!")
	}

	setupWithArkadeIdempotent(configFilePath, config.Tools)
}

// setupWithArkadeIdempotent checks if the tool exists in the config,
// if not, adds it and installs/reinstalls it using arkade.
func setupWithArkadeIdempotent(configFilePath string, toolsToProcess []string) {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error fetching user details: %s\n", err)
		return
	}
	binDir := filepath.Join(usr.HomeDir, ".arkade", "bin")
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

		toolPath := filepath.Join(binDir, tool)
		if _, err := os.Stat(toolPath); err == nil && !*forceFlag {
			fmt.Printf("Tool: %s already installed. Skipping...\n", tool)
			continue
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
		if *passthroughFlag {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error: Failed to install %s via arkade. Moving on to the next tool.\n", tool)
		} else {
			fmt.Printf("Successfully installed %s.\n", tool)
		}
	}
	fmt.Println("Everything is in sync!")
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
func uninstallExtraneousTools(configFilePath string) int {
	config, err := readConfig(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
		return -1
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error fetching user details: %s\n", err)
		return -1
	}
	binDir := filepath.Join(usr.HomeDir, ".arkade", "bin")
	files, err := ioutil.ReadDir(binDir)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		return -1
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

	return removedTools
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

func updateShellConfig() {
	// Determine the directory containing the arkade-lvlup executable
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error determining the arkade-lvlup location:", err)
		return
	}
	arkadeLvlupDir := filepath.Dir(exePath)

	shell := os.Getenv("SHELL")
	var configFile string

	// Determine which config file to update based on the shell
	switch {
	case strings.Contains(shell, "zsh"):
		configFile = ".zshrc"
	case strings.Contains(shell, "bash"):
		configFile = ".bashrc"
	default:
		fmt.Println("Unsupported shell.")
		return
	}

	configPath := filepath.Join(os.Getenv("HOME"), configFile)

	// Check if the block already exists
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading %s: %s\n", configFile, err)
		return
	}

	relativePath, err := filepath.Rel(os.Getenv("HOME"), arkadeLvlupDir)
	if err != nil {
		fmt.Println("Error computing relative path:", err)
		return
	}

	block := fmt.Sprintf(`# Check for arkade and arkade-lvlup
if command -v arkade &> /dev/null; then
    # Add arkade-lvlup to PATH if it exists
    if [[ -f "$HOME/%s/arkade-lvlup" ]]; then
        export PATH="$HOME/%s:$PATH"
    fi
fi`, relativePath, relativePath)

	if strings.Contains(string(content), block) {
		fmt.Println("Configuration already set up.")
		return
	}

	// Confirm with the user
	fmt.Printf("Do you want to update %s to include arkade-lvlup in the PATH? (y/n): ", configFile)
	var response string
	fmt.Scanln(&response)
	if response != "y" {
		fmt.Println("Aborted.")
		return
	}

	// Append the block to the file
	err = ioutil.WriteFile(configPath, append(content, []byte("\n"+block)...), 0644)
	if err != nil {
		fmt.Printf("Error updating %s: %s\n", configFile, err)
		return
	}

	fmt.Println("Configuration updated successfully!")
}

func main() {
	// Parsing command-line flags
	syncFlag := flag.Bool("sync", false, "Sync tools based on configuration.")
	// forceFlag := flag.Bool("f", false, "Force sync (only valid with -sync).")
	getFlag := flag.String("get", "", "Install or reinstall specified tools.")
	removeFlag := flag.String("remove", "", "Remove specified tools.")
	configShellFlag := flag.Bool("config-shell", false, "Update the shell configuration to include arkade-lvlup in the PATH.")
	// passthroughFlag := flag.Bool("passthrough", false, "Show arkade outputs.")

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
			config, err := readConfig(configFilePath)
			if err != nil {
				fmt.Printf("Error reading config: %s\n", err)
				return
			}
			setupWithArkadeIdempotent(configFilePath, config.Tools)
		}
	}

	if *getFlag != "" {
		tools := populateArray(*getFlag)

		// Update the configuration with the tools provided in getFlag
		config, err := readConfig(configFilePath)
		if err != nil {
			fmt.Printf("Error reading config: %s\n", err)
			return
		}
		for _, tool := range tools {
			if !containsElement(config.Tools, tool) {
				config.Tools = append(config.Tools, tool)
				err := writeConfig(configFilePath, config)
				if err != nil {
					fmt.Printf("Error writing to config: %s\n", err)
					return
				}
				fmt.Printf("Added %s to the configuration file.\n", tool)
			}
		}

		if *forceFlag {
			setupWithArkadeIdempotent(configFilePath, tools)
		} else {
			// Only install tools that are not present in arkade/bin
			usr, err := user.Current()
			if err != nil {
				fmt.Printf("Error fetching user details: %s\n", err)
				return
			}
			binDir := filepath.Join(usr.HomeDir, ".arkade", "bin")

			var newTools []string
			for _, tool := range tools {
				if _, err := os.Stat(filepath.Join(binDir, tool)); os.IsNotExist(err) {
					newTools = append(newTools, tool)
				}
			}
			setupWithArkadeIdempotent(configFilePath, newTools)
		}
	}

	if *removeFlag != "" {
		tools := populateArray(*removeFlag)
		removeWithArkade(configFilePath, tools)
	}

	if *configShellFlag {
		updateShellConfig()
	}

}
