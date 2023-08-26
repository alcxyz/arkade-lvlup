package tools

import (
	"io/ioutil"
	"os/user"
	"path/filepath"
)

// GetBinDir retrieves the path to the arkade/bin directory for the current user.
func GetBinDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".arkade", "bin"), nil
}

// ListToolsInBinDir lists all tools present in the arkade/bin directory.
func ListToolsInBinDir() ([]string, error) {
	binDir, err := GetBinDir()
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(binDir)
	if err != nil {
		return nil, err
	}

	var tools []string
	for _, file := range files {
		if !file.IsDir() {
			tools = append(tools, file.Name())
		}
	}

	return tools, nil
}
