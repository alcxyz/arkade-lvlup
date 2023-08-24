package config

type ToolConfig struct {
	Tools []string `yaml:"tools"`
}

type capturedFlags struct {
	force       bool
	passthrough bool
	sync        bool
	get         string
	remove      string
	configShell bool
}
