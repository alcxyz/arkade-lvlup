package tools

import (
	"arkade-lvlup/config"
	"fmt"
	"os"
	"path/filepath"
)

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
