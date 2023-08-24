package tools

import (
	"arkade-lvlup/config"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SetupWithArkadeIdempotent checks if the tool exists in the config,
// if not, adds it and installs/reinstalls it using arkade.
func SetupWithArkadeIdempotent(configFilePath string, toolsToProcess []string, forceFlag bool, passthroughFlag bool) error {
	binDir, err := GetBinDir()
	if err != nil {
		return fmt.Errorf("error fetching bin directory: %w", err)
	}

	for _, tool := range toolsToProcess {
		tool = strings.TrimSpace(tool)
		if tool == "" {
			continue
		}

		cfg, err := config.ReadConfig(configFilePath)
		if err != nil {
			fmt.Printf("Error reading config: %s\n", err)
			return err
		}

		toolPath := filepath.Join(binDir, tool)
		if _, err := os.Stat(toolPath); err == nil && !forceFlag {
			fmt.Printf("Tool: %s already installed. Skipping...\n", tool)
			continue
		}

		if !ContainsElement(cfg.Tools, tool) {
			cfg.Tools = append(cfg.Tools, tool)
			err := config.WriteConfig(configFilePath, cfg)
			if err != nil {
				fmt.Printf("Error writing to config: %s\n", err)
				return err
			}
			fmt.Printf("Added %s to the configuration file.\n", tool)
		}

		fmt.Printf("Installing tool: %s...\n", tool)
		cmd := exec.Command("arkade", "get", tool)
		if passthroughFlag {
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
	return nil
}

// ForceSyncTools forcefully syncs tools with the configuration.
func ForceSyncTools(configFilePath string, forceFlag bool, passthroughFlag bool) error {
	fmt.Println("Synchronizing tools...")

	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
		return
	}

	// removedTools should be calculated and returned by uninstallExtraneousTools
	removedTools := uninstallExtraneousTools(configFilePath)
	if removedTools == 0 && !forceFlag {
		fmt.Println("Everything is in sync!")
	}

	SetupWithArkadeIdempotent(configFilePath, cfg.Tools, forceFlag, passthroughFlag)
	return nil
}

// removeWithArkade removes tools from the config and uninstalls them using arkade.
func RemoveWithArkade(configFilePath string, toolsToRemove []string) error {
	for _, tool := range toolsToRemove {
		cfg, err := config.ReadConfig(configFilePath)
		if err != nil {
			fmt.Printf("Error reading config: %s\n", err)
			return
		}

		// Find and remove the tool from the Tools slice
		for i, t := range cfg.Tools {
			if t == tool {
				config.Tools = append(config.Tools[:i], config.Tools[i+1:]...)
				break
			}
		}

		err = config.WriteConfig(configFilePath, cfg)
		if err != nil {
			fmt.Printf("Error writing to config: %s\n", err)
			return err
		}

		fmt.Printf("Marked %s for removal.\n", tool)
	}

	uninstallExtraneousTools(configFilePath)
	return nil
}

// uninstallExtraneousTools removes tools that are present in the arkade directory but not in the config.
func uninstallExtraneousTools(configFilePath string) int {
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
		if !ContainsElement(config.Tools, file.Name()) {
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

// GetSyncState provides an overview of the current sync state.
func GetSyncState(configFilePath string) {
	binDir, err := GetBinDir()
	if err != nil {
		fmt.Printf("Error fetching bin directory: %s\n", err)
		return
	}

	installedTools, err := ListToolsInBinDir(binDir)
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
}
