package cmd

import (
	"arkade-lvlup/config"
	"arkade-lvlup/handlers"

	"github.com/spf13/cobra"
)

var syncFlags = &config.CmdFlags{}

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync tools based on configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		handlers.HandleSync(args, syncFlags)
	},
}

func init() {
	config.RegisterCommonFlags(SyncCmd.Flags(), syncFlags)
}
