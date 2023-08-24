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
	"strings"
)

var (
	forceFlag       = flag.Bool("f", false, "Force sync.")
	passthroughFlag = flag.Bool("p", false, "Show arkade outputs.")
	syncFlag        = flag.Bool("sync", false, "Sync tools based on configuration. (Alias: -s)")
	getFlag         = flag.String("get", "", "Install or reinstall specified tools. (Alias: -g)")
	removeFlag      = flag.String("remove", "", "Remove specified tools. (Alias: -r)")
	configShellFlag = flag.Bool("config-shell", false, "Update the shell configuration to include arkade-lvlup in the PATH. (Alias: -c)")
)

var orderedFlags = []string{"passthrough", "forceSync", "configShell", "get", "remove"}

var flagHandlers = map[string]func(capturedFlags, string) error{
	"passthrough": func(c capturedFlags, configFilePath string) error {
		// Handle passthrough functionality
		return nil
	},
	"forceSync": func(c capturedFlags, configFilePath string) error {
		if c.force && !c.sync {
			return errors.New("-f or --force can only be used with --sync or -s.")
		}
		if c.sync {
			return handleSync(configFilePath, c.force, c.passthrough)
		}
		return nil
	},
	"configShell": func(c capturedFlags, configFilePath string) error {
		if c.configShell {
			return handleShellConfig()
		}
		return nil
	},
	"get": func(c capturedFlags, configFilePath string) error {
		if c.get != "" {
			return handleGet(configFilePath, c.get, c.force)
		}
		return nil
	},
	"remove": func(c capturedFlags, configFilePath string) error {
		if c.remove != "" {
			return handleRemove(configFilePath, c.remove)
		}
		return nil
	},
}

func init() {
	flag.BoolVar(syncFlag, "s", false, "Alias for --sync.")
	flag.StringVar(getFlag, "g", "", "Alias for --get.")
	flag.StringVar(removeFlag, "r", "", "Alias for --remove.")
	flag.BoolVar(configShellFlag, "c", false, "Alias for --config-shell.")
}

type ToolConfig struct {
	Tools []string `yaml:"tools"`
}

type capturedFlags struct {
	force       bool
	passthrough bool
	sync        bool
	get         string
	remove      string
	configShell bool
}

func captureAllFlags() capturedFlags {
	getTools := *getFlag
	removeTools := *removeFlag

	// Exclude flags from the tools list
	for _, flagVal := range []string{"-p", "-f", "-s", "-c", "-g", "-r"} {
		getTools = removeFlagFromToolList(getTools, flagVal)
		removeTools = removeFlagFromToolList(removeTools, flagVal)
	}

	return capturedFlags{
		force:       *forceFlag,
		passthrough: *passthroughFlag,
		sync:        *syncFlag,
		get:         getTools,
		remove:      removeTools,
		configShell: *configShellFlag,
	}
}

// Remove flag from tool list if accidentally passed as a tool name
func removeFlagFromToolList(toolList, flagVal string) string {
	return strings.ReplaceAll(toolList, flagVal, "")
}

func main() {
	flag.Parse()

	configDir, err := tools.GetConfigDir()
	if err != nil {
		log.Fatalf("Error getting config directory: %s\n", err)
	}
	configFilePath := filepath.Join(configDir, "lvlup.yaml")

	// Check and initialize the config if it doesn't exist
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		_, err := tools.InitializeConfigIfNotExists(configFilePath)
		if err != nil {
			log.Fatalf("Error initializing configuration: %s\n", err)
		}
	}

	c := captureAllFlags()
	if err := handleFlagsInOrder(c, configFilePath); err != nil {
		log.Fatalf("Error handling flags: %s", err)
	}
}

func handleFlagsInOrder(c capturedFlags, configFilePath string) error {
	for _, flagKey := range orderedFlags {
		handler, exists := flagHandlers[flagKey]
		if !exists {
			return fmt.Errorf("No handler found for flag: %s", flagKey)
		}

		err := handler(c, configFilePath)
		if err != nil {
			return err
		}
	}

	return handleDefaultState(configFilePath)
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
	tools.GetSyncState(configFilePath)
	return nil
}
