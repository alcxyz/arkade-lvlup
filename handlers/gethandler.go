package handlers

import (
	"fmt"
	"log"
	"strings"

	"arkade-lvlup/config"
	"arkade-lvlup/tools"
)

func HandleGet(args []string, flags *config.CmdFlags) {

	configFilePath, err := tools.GetConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}

	// Read the config
	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error reading config: %s", err)
	}

	// Convert tools to get from arguments to a slice
	toolsToGet := tools.PopulateArray(strings.Join(args, ","))

	for _, tool := range toolsToGet {
		if !tools.ContainsElement(cfg.Tools, tool) {
			cfg.Tools = append(cfg.Tools, tool)
			err = config.WriteConfig(configFilePath, cfg)
			if err != nil {
				log.Fatalf("Error writing to config: %s", err)
			}
			fmt.Printf("Added %s to the configuration file.\n", tool)
		}
	}

	err = tools.InstallToolsIdempotently(configFilePath, toolsToGet, flags.Force, flags.Passthrough)
	if err != nil {
		log.Fatalf("Failed to install tools: %s", err)
	}
}
