package config

import (
	"github.com/spf13/pflag"
)

type CmdFlags struct {
	Force       bool
	Passthrough bool
}

func RegisterCommonFlags(fs *pflag.FlagSet, flags *CmdFlags) {
	fs.BoolVarP(&flags.Passthrough, "passthrough", "p", false, "Show arkade outputs.")
	fs.BoolVarP(&flags.Force, "force", "f", false, "Force sync.")
	// ... other flag bindings if needed in the future
}
