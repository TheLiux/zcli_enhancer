package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: unable to get current directory.")
		return
	}

	envFilePath := filepath.Join(currentDir, ".env")
	manifestFilePath := filepath.Join(currentDir, "manifest.json")

	if _, err := os.Stat(manifestFilePath); os.IsNotExist(err) {
		fmt.Println("Error: manifest.json does not exist in the current directory.")
		return
	}

	errorChannel := make(chan error)
	go checkEnv(envFilePath, errorChannel)

	if err := <-errorChannel; err != nil {
		fmt.Println(err)
		return
	}

	if err := updateVersion(manifestFilePath); err != nil {
		fmt.Println("Error updating version:", err)
		return
	}

	fmt.Println("Theme successfully updated.")
}

func checkEnv(envFilePath string, errorChannel chan error) {
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		themeID := askForThemeIDAndCreateEnv(envFilePath)
		updateTheme(themeID)
	} else {
		content, err := os.ReadFile(envFilePath)
		if err != nil {
			errorChannel <- fmt.Errorf("failed to read .env file: %w", err)
			return
		}
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "THEME_ID=") {
				themeID := strings.TrimPrefix(line, "THEME_ID=")
				updateTheme(themeID)
				errorChannel <- nil
				return
			}
		}
		errorChannel <- fmt.Errorf("THEME_ID not found in .env")
	}
}

func askForThemeIDAndCreateEnv(envFilePath string) string {
	fmt.Println("Please enter the theme ID:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	themeID := scanner.Text()
	fileContent := fmt.Sprintf("THEME_ID=%s\n", themeID)
	os.WriteFile(envFilePath, []byte(fileContent), 0644)
	fmt.Printf(".env created with ID: %s\n", themeID)
	return themeID
}

func updateTheme(themeID string) {
	fmt.Printf("Updating theme with ID: %s\n", themeID)
	bar := progressbar.Default(100, "Updating theme")
	for i := 0; i < 100; i++ {
		bar.Add(1)
	}
	cmd := exec.Command("zcli", "themes:update", "--themeId="+themeID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to update theme: %v\n", err)
	}
}

func updateVersion(manifestFilePath string) error {
	file, err := os.ReadFile(manifestFilePath)
	if err != nil {
		return fmt.Errorf("failed to read manifest.json: %w", err)
	}

	var manifest map[string]interface{}
	if err := json.Unmarshal(file, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest.json: %w", err)
	}

	version, ok := manifest["version"].(string)
	if !ok {
		return fmt.Errorf("unable to find or parse version in manifest.json")
	}

	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return fmt.Errorf("version format is not correct. Expected format: major.minor.patch")
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("failed to parse patch number: %w", err)
	}
	patch++
	newVersion := fmt.Sprintf("%s.%s.%d", parts[0], parts[1], patch)
	manifest["version"] = newVersion

	updatedFile, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated manifest.json: %w", err)
	}

	if err := os.WriteFile(manifestFilePath, updatedFile, 0644); err != nil {
		return fmt.Errorf("failed to write updated manifest.json: %w", err)
	}

	fmt.Printf("Updating version: from %s to %s\n", version, newVersion)
	return nil
}
