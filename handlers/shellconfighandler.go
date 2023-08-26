package handlers

import (
	"arkade-lvlup/shell"
	"log"
)

func HandleShellConfig() {
	err := shell.UpdateShellConfig()
	if err != nil {
		log.Fatal(err)
	}
}
