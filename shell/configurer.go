package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func UpdateShellConfig() {
	// Determine the directory containing the arkade-lvlup executable
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error determining the arkade-lvlup location:", err)
		return
	}
	arkadeLvlupDir := filepath.Dir(exePath)

	shell := os.Getenv("SHELL")
	var configFile string

	// Determine which config file to update based on the shell
	switch {
	case strings.Contains(shell, "zsh"):
		configFile = ".zshrc"
	case strings.Contains(shell, "bash"):
		configFile = ".bashrc"
	default:
		fmt.Println("Unsupported shell.")
		return
	}

	configPath := filepath.Join(os.Getenv("HOME"), configFile)

	// Check if the block already exists
	content, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading %s: %s\n", configFile, err)
		return
	}

	relativePath, err := filepath.Rel(os.Getenv("HOME"), arkadeLvlupDir)
	if err != nil {
		fmt.Println("Error computing relative path:", err)
		return
	}

	block := fmt.Sprintf(`# Check for arkade and arkade-lvlup
if command -v arkade &> /dev/null; then
    # Add arkade-lvlup to PATH if it exists
    if [[ -f "$HOME/%s/arkade-lvlup" ]]; then
        export PATH="$HOME/%s:$PATH"
    fi
fi`, relativePath, relativePath)

	if strings.Contains(string(content), block) {
		fmt.Println("Configuration already set up.")
		return
	}

	// Confirm with the user
	fmt.Printf("Do you want to update %s to include arkade-lvlup in the PATH? (y/n): ", configFile)
	var response string
	fmt.Scanln(&response)
	if response != "y" {
		fmt.Println("Aborted.")
		return
	}

	// Append the block to the file
	err = os.WriteFile(configPath, append(content, []byte("\n"+block)...), 0644)
	if err != nil {
		fmt.Printf("Error updating %s: %s\n", configFile, err)
		return
	}

	fmt.Println("Configuration updated successfully!")
}
