package handlers

import (
	"arkade-lvlup/tools"
	"log"
)

func HandleRemove(args []string) {
	configFilePath, err := tools.GetConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}

	toolsToRemove := tools.PopulateArray(args[0])

	err = tools.RemoveWithArkade(configFilePath, toolsToRemove)
	if err != nil {
		log.Fatal(err)
	}

	err = tools.RemoveToolsFromFS(toolsToRemove)
	if err != nil {
		log.Fatal(err)
	}

	tools.SyncFileSystemWithConfig(configFilePath)
}
