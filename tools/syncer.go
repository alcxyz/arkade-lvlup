package tools

import (
	"arkade-lvlup/config"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// SyncToolsForcefully ensures all tools in the configuration file are installed by force. It will also remove
// any tools found in the arkade directory but not in the configuration.
func SyncToolsForcefully(configFilePath string, passthroughFlag bool) error {
	fmt.Println("Force synchronization initiated...")

	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	installedTools, err := ListToolsInBinDir()
	if err != nil {
		return fmt.Errorf("error listing tools in bin directory: %w", err)
	}

	fmt.Println("Reinstalling all tools in the bin directory...")
	InstallToolsIdempotently(configFilePath, installedTools, true, passthroughFlag)

	fmt.Println("Installing tools from the configuration that are not in the bin directory...")
	InstallToolsIdempotently(configFilePath, cfg.Tools, false, passthroughFlag)

	fmt.Println("Force synchronization completed!")
	return nil
}

// SyncFilesSystemWithConfig identifies and removes tools that are present in the arkade directory
// but not in the configuration file.
func SyncFileSystemWithConfig(configFilePath string) int {

	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
		return -1
	}

	binDir, err := GetBinDir()
	if err != nil {
		fmt.Printf("Error fetching bin directory: %s\n", err)
		return -1
	}

	files, err := ioutil.ReadDir(binDir)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		return -1
	}

	removedTools := 0
	for _, file := range files {
		if !ContainsElement(cfg.Tools, file.Name()) {
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
		fmt.Println("No extraneous tools found. Bin directory is in sync with configuration!")
	} else {
		fmt.Printf("%d extraneous tools removed from bin directory.\n", removedTools)
	}

	return removedTools
}

// GetSyncState prints the state of tools synchronization. It checks which tools are installed but not
// in the config, and which tools are in the config but not installed. It then prints a summary of the results.
func GetSyncState(configFilePath string) {
	installedTools, err := ListToolsInBinDir()
	if err != nil {
		fmt.Printf("Error listing tools: %s\n", err)
		return
	}

	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
		return
	}

	inConfigNotInstalled := make([]string, 0)
	for _, tool := range cfg.Tools {
		if !ContainsElement(installedTools, tool) {
			inConfigNotInstalled = append(inConfigNotInstalled, tool)
		}
	}

	installedNotInConfig := make([]string, 0)
	for _, tool := range installedTools {
		if !ContainsElement(cfg.Tools, tool) {
			installedNotInConfig = append(installedNotInConfig, tool)
		}
	}

	fmt.Println("----- Sync State -----")
	if len(inConfigNotInstalled) > 0 {
		fmt.Printf("Tools in config but not installed: %s\n", strings.Join(inConfigNotInstalled, ", "))
	} else {
		fmt.Println("All tools in the config are installed!")
	}

	if len(installedNotInConfig) > 0 {
		fmt.Printf("Tools installed but not in config: %s\n", strings.Join(installedNotInConfig, ", "))
	} else {
		fmt.Println("All installed tools are in the config!")
	}
	fmt.Println("----------------------")

	// Displaying the tools managed by arkade-lvlup
	if len(cfg.Tools) == 0 {
		fmt.Println("No tools are currently managed by arkade-lvlup.")
	} else {
		fmt.Println("Tools managed by arkade-lvlup:")
		for _, tool := range cfg.Tools {
			fmt.Println("-", tool)
		}
	}
}
