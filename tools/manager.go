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

// runArkadeCommand executes the arkade command to fetch and install the given tool.
func runArkadeCommand(tool string, passthroughFlag bool) error {
	cmd := exec.Command("arkade", "get", tool)
	if passthroughFlag {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

// toolInstalled checks if a given tool is already installed.
func toolInstalled(tool string) (bool, error) {
	binDir, err := GetBinDir()
	if err != nil {
		return false, fmt.Errorf("error fetching bin directory: %w", err)
	}
	toolPath := filepath.Join(binDir, tool)
	if _, err := os.Stat(toolPath); err == nil {
		return true, nil
	}
	return false, nil
}

// InstallToolsIdempotently installs tools using arkade. If the tool is already installed and
// forceFlag is not set, it skips the installation. If the tool isn't in the configuration file,
// it adds it there.
func InstallToolsIdempotently(configFilePath string, toolsToProcess []string, forceFlag bool, passthroughFlag bool) error {
	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	for _, tool := range toolsToProcess {
		tool = strings.TrimSpace(tool)
		if tool == "" {
			continue
		}

		installed, err := toolInstalled(tool)
		if err != nil {
			return err
		}

		if installed && !forceFlag {
			fmt.Printf("Tool: %s already installed. Skipping...\n", tool)
			continue
		}

		if !ContainsElement(cfg.Tools, tool) {
			cfg.Tools = append(cfg.Tools, tool)
			err := config.WriteConfig(configFilePath, cfg)
			if err != nil {
				return fmt.Errorf("error writing to config: %w", err)
			}
			fmt.Printf("Added %s to the configuration file.\n", tool)
		}

		err = runArkadeCommand(tool, passthroughFlag)
		if err != nil {
			fmt.Printf("Error: Failed to install %s via arkade. Moving on to the next tool.\n", tool)
		} else {
			fmt.Printf("Successfully installed %s.\n", tool)
		}
	}
	return nil
}

// SyncToolsForcefully ensures all tools in the configuration file are installed by force. It will also remove
// any tools found in the arkade directory but not in the configuration.
func SyncToolsForcefully(configFilePath string, passthroughFlag bool) error {
	fmt.Println("Synchronizing tools...")

	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	SyncFileSystemWithConfig(configFilePath)
	InstallToolsIdempotently(configFilePath, cfg.Tools, true, passthroughFlag)
	return nil
}

func RemoveWithArkade(configFilePath string, toolsToRemove []string) error {
	for _, tool := range toolsToRemove {
		cfg, err := config.ReadConfig(configFilePath)
		if err != nil {
			fmt.Printf("Error reading config: %s\n", err)
			return err
		}

		// Find and remove the tool from the Tools slice
		for i, t := range cfg.Tools {
			if t == tool {
				cfg.Tools = append(cfg.Tools[:i], cfg.Tools[i+1:]...)
				break
			}
		}

		err = config.WriteConfig(configFilePath, cfg)
		if err != nil {
			fmt.Printf("Error writing to config: %s\n", err)
			return err
		}

		fmt.Printf("Marked %s for removal from configuration.\n", tool)
	}
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
		fmt.Println("No extraneous tools found. Everything is in sync!")
	}

	return removedTools
}

func RemoveToolsFromFS(toolsToRemove []string) error {
	binDir, err := GetBinDir()
	if err != nil {
		return err
	}

	for _, tool := range toolsToRemove {
		toolPath := filepath.Join(binDir, tool)
		err = os.Remove(toolPath)
		if err != nil {
			fmt.Printf("Failed to remove tool %s from file system: %s\n", tool, err)
			// Deciding to just print the error rather than halting the entire process.
			// If you want to stop, you can return the error here.
		}
	}
	return nil
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
