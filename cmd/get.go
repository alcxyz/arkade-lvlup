package cmd

import (
	"arkade-lvlup/config"
	"arkade-lvlup/handlers"

	"github.com/spf13/cobra"
)

var getFlags = &config.CmdFlags{}

var GetCmd = &cobra.Command{
	Use:   "get [tools]",
	Short: "Install or reinstall specified tools.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handlers.HandleGet(args, getFlags)
	},
}

func init() {
	config.RegisterCommonFlags(GetCmd.Flags(), getFlags)
}
