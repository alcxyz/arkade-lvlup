package cmd

import (
	"arkade-lvlup/config"
	"arkade-lvlup/handlers"

	"github.com/spf13/cobra"
)

var removeFlags = &config.CmdFlags{}

var RemoveCmd = &cobra.Command{
	Use:   "remove [tools]",
	Short: "Remove specified tools.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handlers.HandleRemove(args)
	},
}

func init() {
	config.RegisterCommonFlags(RemoveCmd.Flags(), removeFlags)
}
