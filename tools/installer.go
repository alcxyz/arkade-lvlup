package tools

import (
	"arkade-lvlup/config"
	"fmt"
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
