package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

// Config holds basic information required to run StHub
type Config struct {
	// WowsPath is the file path to the game directory
	WowsPath string

	// APIPath is the path used for the "scraper" component
	APIPath string `json:""`
}

// GetConfigPath returns the path where sthub data should be stored
func GetConfigPath() (string, error) {
	dir := os.Getenv("LocalAppData")

	if dir != "" {
		return createConfigDir(dir)
	}

	log.Printf("Falling back to home dir, LocalAppData is not available")
	dir = os.Getenv("HOME")

	if dir != "" {
		return createConfigDir(dir)
	}

	return "", fmt.Errorf("LocalAppData and HOME not available")
}

// createConfigDir creates the sthub config dir inside of the
// configuration base path
func createConfigDir(baseDir string) (string, error) {
	if err := os.Mkdir(path.Join(baseDir, "sthub"), os.ModePerm); err != nil && !os.IsExist(err) {
		return baseDir, err
	}

	return path.Join(baseDir, "sthub"), nil
}

// Returns whether a local config file exists
func hasLocalConfig() bool {
	configPath, err := GetConfigPath()
	if err != nil {
		log.Fatalf("Could not retrieve a config path: %v", err)
		return false
	}

	if _, err := os.Stat(path.Join(configPath, "sthub-config.json")); os.IsNotExist(err) {
		return false
	}

	return true
}
