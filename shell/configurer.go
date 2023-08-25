package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func UpdateShellConfig() error {
	// Determine the directory containing the arkade-lvlup executable
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("Error determining the arkade-lvlup location: %w", err)
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
		return fmt.Errorf("Unsupported shell.")
	}

	configPath := filepath.Join(os.Getenv("HOME"), configFile)

	// Check if the block already exists
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Error reading %s: %w", configFile, err)
	}

	relativePath, err := filepath.Rel(os.Getenv("HOME"), arkadeLvlupDir)
	if err != nil {
		return fmt.Errorf("Error computing relative path: %w", err)
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
		return nil
	}

	// Confirm with the user
	fmt.Printf("Do you want to update %s to include arkade-lvlup in the PATH? (y/n): ", configFile)
	var response string
	fmt.Scanln(&response)
	if response != "y" {
		fmt.Println("Aborted.")
		return nil
	}

	// Append the block to the file
	err = os.WriteFile(configPath, append(content, []byte("\n"+block)...), 0644)
	if err != nil {
		return fmt.Errorf("Error updating %s: %w", configFile, err)
	}

	fmt.Println("Configuration updated successfully!")
	return nil
}
