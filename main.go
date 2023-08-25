package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"arkade-lvlup/config"
	"arkade-lvlup/shell"
	"arkade-lvlup/tools"

	"github.com/spf13/cobra"
)

// Global Variables for capturing flag values
var (
	force           bool
	passthroughFlag bool
	configShell     bool
	configFilePath  string
)

var globalTools config.ArkadeTools

// Root command represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arkade-lvlup",
	Short: "lvlup - arkade CLI tool manager",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		configFilePath, err = getConfigFilePath()
		if err != nil {
			log.Fatal(err)
		}

		globalTools, err = config.ReadConfig(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		handleDefaultState(configFilePath)
	},
}

// Define sub-commands
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync tools based on configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		err := handleSync(globalTools, configFilePath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "get [tools]",
	Short: "Install or reinstall specified tools.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		toolsToGet := strings.Join(args, ",")
		err := handleGet(globalTools, toolsToGet, force, configFilePath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [tools]",
	Short: "Remove specified tools.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		toolsToRemove := strings.Join(args, ",")
		err := handleRemove(globalTools, toolsToRemove, configFilePath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var shellCmd = &cobra.Command{
	Use:   "config-shell",
	Short: "Update the shell configuration to include arkade-lvlup in the PATH.",
	Run: func(cmd *cobra.Command, args []string) {
		err := handleShellConfig()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// Add flags to commands
	syncCmd.Flags().BoolVarP(&force, "force", "f", false, "Force sync.")
	syncCmd.Flags().BoolVarP(&passthroughFlag, "passthrough", "p", false, "Show arkade outputs.")
	getCmd.Flags().BoolVarP(&passthroughFlag, "passthrough", "p", false, "Show arkade outputs.")
	removeCmd.Flags().BoolVarP(&passthroughFlag, "passthrough", "p", false, "Show arkade outputs.")

	// Add sub-commands to root command
	rootCmd.AddCommand(syncCmd, getCmd, removeCmd, shellCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleSync(cfg config.ArkadeTools, configFilePath string) error {
	if force {
		return tools.SyncToolsForcefully(configFilePath, passthroughFlag)
	}
	return tools.InstallToolsIdempotently(configFilePath, cfg.Tools, force, passthroughFlag)
}

func handleGet(cfg config.ArkadeTools, getFlagValue string, force bool, configFilePath string) error {
	toolNames := tools.PopulateArray(getFlagValue)

	// Moved the reading of the config outside the loop.
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

	err = tools.InstallToolsIdempotently(configFilePath, toolNames, force, passthroughFlag)
	if err != nil {
		return fmt.Errorf("failed to install tools: %w", err)
	}
	return nil
}

func handleRemove(cfg config.ArkadeTools, removeFlagValue string, configFilePath string) error {
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
	return shell.UpdateShellConfig()
}

func handleDefaultState(configFilePath string) error {
	tools.GetSyncState(configFilePath)
	return nil
}

func getConfigFilePath() (string, error) {
	configDir, err := tools.GetConfigDir()
	if err != nil {
		return "", fmt.Errorf("Error getting config directory: %s", err)
	}
	return filepath.Join(configDir, "lvlup.yaml"), nil
}
