package handlers

import (
	"arkade-lvlup/tools"
	"log"
	"os"
	"path/filepath"
)

func HandleRemove(args []string) {
	configFilePath, err := tools.GetConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}

	// Commented out code that reads the config but never uses it.
	// If you need this later, you can uncomment it.
	// globalTools, err := config.ReadConfig(configFilePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	toolsToRemove := tools.PopulateArray(args[0]) // Assuming the arguments are a list of tools to remove

	err = tools.RemoveWithArkade(configFilePath, toolsToRemove)
	if err != nil {
		log.Fatal(err)
	}

	binDir, _ := tools.GetBinDir()
	for _, tool := range toolsToRemove {
		toolPath := filepath.Join(binDir, tool)

		err = os.Remove(toolPath)
		if err != nil {
			log.Fatalf("Failed to remove tool %s: %s", tool, err)
		}
	}
}
