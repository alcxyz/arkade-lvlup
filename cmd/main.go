package main

import (
	"arkade-lvlup/config"
	"arkade-lvlup/shell"
	"arkade-lvlup/tools"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Global Variables for capturing flag values
var (
	force           bool
	passthroughFlag bool
	configShell     bool
)

// Root command represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arkade-lvlup",
	Short: "Arkade LvlUp CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior when no flag is provided
		configFilePath, err := getConfigFilePath()
		if err != nil {
			log.Fatal(err)
		}
		handleDefaultState(configFilePath)
	},
}

// Define sub-commands
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync tools based on configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath, err := getConfigFilePath()
		if err != nil {
			log.Fatal(err)
		}
		err = handleSync(configFilePath)
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
		configFilePath, err := getConfigFilePath()
		if err != nil {
			log.Fatal(err)
		}
		toolsToGet := strings.Join(args, ",")              // Joining all arguments with commas
		err = handleGet(configFilePath, toolsToGet, force) // Removed passthroughFlag
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
		configFilePath, err := getConfigFilePath()
		if err != nil {
			log.Fatal(err)
		}
		toolsToRemove := strings.Join(args, ",")          // Joining all arguments with commas
		err = handleRemove(configFilePath, toolsToRemove) // Removed passthroughFlag
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

func handleSync(configFilePath string) error {
	if force {
		return tools.ForceSyncTools(configFilePath, passthroughFlag)
	}

	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	return tools.SetupWithArkadeIdempotent(configFilePath, cfg.Tools, force, passthroughFlag)
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

	err = tools.SetupWithArkadeIdempotent(configFilePath, toolNames, force, passthroughFlag)
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
