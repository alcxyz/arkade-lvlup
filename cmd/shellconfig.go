package cmd

import (
	"arkade-lvlup/handlers"

	"github.com/spf13/cobra"
)

var ShellCmd = &cobra.Command{
	Use:   "config-shell",
	Short: "Update the shell configuration to include arkade-lvlup in the PATH.",
	Run: func(cmd *cobra.Command, args []string) {
		handlers.HandleShellConfig()
	},
}
