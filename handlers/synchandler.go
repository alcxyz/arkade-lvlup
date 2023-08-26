package handlers

import (
	"arkade-lvlup/config"
	"arkade-lvlup/tools"
	"log"
)

func HandleSync(args []string, flags *config.CmdFlags) {
	configFilePath, err := tools.GetConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}

	globalTools, err := config.ReadConfig(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if flags.Force {
		err = tools.SyncToolsForcefully(configFilePath, flags.Passthrough)
	} else {
		err = tools.InstallToolsIdempotently(configFilePath, globalTools.Tools, flags.Force, flags.Passthrough)
	}

	if err != nil {
		log.Fatal(err)
	}
}
