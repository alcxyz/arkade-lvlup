package main

import (
	"arkade-lvlup/tools"
	"flag"
	"fmt"
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

type Config struct {
	Tools []string `yaml:"tools"`
}

func main() {
	flag.Parse()

	binDir, err := tools.GetBinDir()
	if err != nil {
		fmt.Printf("Error getting bin directory: %s\n", err)
		return
	}

	configFilePath := filepath.Join(binDir, "lvlup.yaml")

	// Check and initialize the config if it doesn't exist
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		_, err := tools.InitializeConfigIfNotExists(configFilePath)
		if err != nil {
			fmt.Printf("Error initializing configuration: %s\n", err)
			return
		}
	}
	handleFlags(configFilePath)
}

func handleFlags(configFilePath string) {
	if *forceFlag && !*syncFlag {
		fmt.Println("Error: -f or --force can only be used with --sync or -s.")
		return
	}

	if *syncFlag {
		handleSync(configFilePath, *forceFlag, *passthroughFlag)
	} else if *getFlag != "" {
		handleGet(configFilePath, *getFlag, *forceFlag)
	} else if *removeFlag != "" {
		handleRemove(configFilePath, *removeFlag)
	} else if *configShellFlag {
		handleShellConfig()
	} else {
		handleDefaultState(configFilePath)
	}
}

// Additional refactored functions such as handleSync, handleGet, handleRemove, etc.

func handleSync(configFilePath string, force bool, passthrough bool) {
	if force {
		tools.ForceSyncTools(configFilePath, force)
	} else {
		config, err := readConfig(configFilePath)
		if err != nil {
			fmt.Printf("Error reading config: %s\n", err)
			return
		}
		tools.SetupWithArkadeIdempotent(configFilePath, config.Tools, force, passthrough)
	}
}

func handleGet(configFilePath string, getFlagValue string, force bool) {
	tools := tools.PopulateArray(getFlagValue)
	config, err := readConfig(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
		return
	}
	for _, tool := range tools {
		if !tools.ContainsElement(config.Tools, tool) {
			config.Tools = append(config.Tools, tool)
			err := writeConfig(configFilePath, config)
			if err != nil {
				fmt.Printf("Error writing to config: %s\n", err)
				return
			}
			fmt.Printf("Added %s to the configuration file.\n", tool)
		}
	}

	if force {
		tools.SetupWithArkadeIdempotent(configFilePath, tools, force, passthrough)
	} else {
		binDir, _ := tools.GetBinDir()
		var newTools []string
		for _, tool := range tools {
			if _, err := os.Stat(filepath.Join(binDir, tool)); os.IsNotExist(err) {
				newTools = append(newTools, tool)
			}
		}
		tools.SetupWithArkadeIdempotent(configFilePath, newTools, force, passthrough)
	}
}

func handleRemove(configFilePath string, removeFlagValue string) {
	toolsToRemove := tools.PopulateArray(removeFlagValue)
	removeWithArkade(configFilePath, toolsToRemove)
}

func handleDefaultState(configFilePath string) {
	getSyncState(configFilePath)
}
