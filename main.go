package main

import (
	"arkade-lvlup/cmd"
	"arkade-lvlup/config"
	"arkade-lvlup/handlers"
	"arkade-lvlup/tools"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	configShell    bool
	configFilePath string
	flags          config.CmdFlags
	globalTools    config.ArkadeTools
)

// Root command represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arkade-lvlup",
	Short: "lvlup - arkade CLI tool manager",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		configFilePath, err = tools.GetConfigFilePath()
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

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync tools based on configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		handlers.HandleSync(args, &flags)
	},
}

var shellCmd = &cobra.Command{
	Use:   "config-shell",
	Short: "Update the shell configuration to include arkade-lvlup in the PATH.",
	Run: func(cmd *cobra.Command, args []string) {
		handlers.HandleShellConfig()
	},
}

func init() {
	config.RegisterCommonFlags(rootCmd.PersistentFlags(), &flags)
	rootCmd.AddCommand(syncCmd, cmd.GetCmd, cmd.RemoveCmd, shellCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleDefaultState(configFilePath string) error {
	tools.GetSyncState(configFilePath)
	return nil
}
