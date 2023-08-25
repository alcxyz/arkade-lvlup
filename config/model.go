package config

// ArkadeTools represents the tools configuration structure derived from the YAML file.
type ArkadeTools struct {
	Tools []string `yaml:"tools"`
}

// CmdFlags captures the flags passed from the command line.
type CmdFlags struct {
	Force       bool
	Passthrough bool
	Sync        bool
	Get         string
	Remove      string
	ConfigShell bool
}
