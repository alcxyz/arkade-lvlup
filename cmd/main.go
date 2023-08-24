package main

import (
	"arkade-lvlup/config"
	"arkade-lvlup/shell"
	"arkade-lvlup/tools"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	forceFlag       = flag.Bool("f", false, "Force sync.")
	passthroughFlag = flag.Bool("p", false, "Show arkade outputs.")
	syncFlag        = flag.Bool("sync", false, "Sync tools based on configuration. (Alias: -s)")
	getFlag         = flag.String("get", "", "Install or reinstall specified tools. (Alias: -g)")
	removeFlag      = flag.String("remove", "", "Remove specified tools. (Alias: -r)")
	configShellFlag = flag.Bool("config-shell", false, "Update the shell configuration to include arkade-lvlup in the PATH. (Alias: -c)")
)

func init() {
	flag.BoolVar(syncFlag, "s", false, "Alias for --sync.")
	flag.StringVar(getFlag, "g", "", "Alias for --get.")
	flag.StringVar(removeFlag, "r", "", "Alias for --remove.")
	flag.BoolVar(configShellFlag, "c", false, "Alias for --config-shell.")
}

type ToolConfig struct {
	Tools []string `yaml:"tools"`
}

func main() {
	flag.Parse()

	binDir, err := tools.GetBinDir()
	if err != nil {
		log.Fatalf("Error getting bin directory: %s\n", err)
	}

	configFilePath := filepath.Join(binDir, "lvlup.yaml")

	// Check and initialize the config if it doesn't exist
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		_, err := tools.InitializeConfigIfNotExists(configFilePath)
		if err != nil {
			log.Fatalf("Error initializing configuration: %s\n", err)
		}
	}

	if err := handleFlags(configFilePath); err != nil {
		log.Fatalf("Error handling flags: %s", err)
	}
}

func handleFlags(configFilePath string) error {
	if *forceFlag && !*syncFlag {
		return errors.New("-f or --force can only be used with --sync or -s.")
	}

	switch {
	case *syncFlag:
		return handleSync(configFilePath, *forceFlag, *passthroughFlag)
	case *getFlag != "":
		return handleGet(configFilePath, *getFlag, *forceFlag)
	case *removeFlag != "":
		return handleRemove(configFilePath, *removeFlag)
	case *configShellFlag:
		return handleShellConfig()
	default:
		return handleDefaultState(configFilePath)
	}
}

func handleSync(configFilePath string, force bool, passthrough bool) error {
	if force {
		return tools.ForceSyncTools(configFilePath, passthrough)
	}

	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	return tools.SetupWithArkadeIdempotent(configFilePath, cfg.Tools, force, passthrough)
}

func handleGet(configFilePath string, getFlagValue string, force bool) error {
	toolNames := tools.PopulateArray(getFlagValue)
	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}
	for _, tool := range toolNames {
		if !tools.ContainsElement(cfg.Tools, tool) {
			cfg.Tools = append(cfg.Tools, tool)
			err = config.WriteConfig(configFilePath, cfg)
			if err != nil {
				return fmt.Errorf("error writing to config: %w", err)
			}
			fmt.Printf("Added %s to the configuration file.\n", tool)
		}
	}

	err = tools.SetupWithArkadeIdempotent(configFilePath, toolNames, force, *passthroughFlag)
	if err != nil {
		return err
	}
	return nil
}

func handleRemove(configFilePath string, removeFlagValue string) error {
	toolsToRemove := tools.PopulateArray(removeFlagValue)
	err := tools.RemoveWithArkade(configFilePath, toolsToRemove)
	if err != nil {
		return err
	}

	binDir, _ := tools.GetBinDir()
	for _, tool := range toolsToRemove {
		toolPath := filepath.Join(binDir, tool)
		os.Remove(toolPath)
	}

	return nil
}

func handleShellConfig() error {
	shell.UpdateShellConfig()
	return nil // For now, we are not handling any errors from UpdateShellConfig
}

func handleDefaultState(configFilePath string) error {
	// TODO: Implement this function
	return nil
}
