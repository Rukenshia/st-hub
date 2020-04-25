package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Migrate063ConfigFiles migrates config files from 0.6.3 (or below) to use
// the AppData directory instead of the directory
// the binary is located in
func Migrate063ConfigFiles(from string) ([]string, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("Could not get config path: %v", err)
	}

	filesToMigrate := []string{}

	// Go through the given directory and find all old config / iteration
	// files. They all start with "sthub-".
	err = filepath.Walk(from, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.Contains(path, "sthub-") {
			return nil
		}

		filesToMigrate = append(filesToMigrate, path)
		newPath := filepath.Join(configPath, path)

		log.Printf("Migrating file: %s to %s", path, newPath)

		if err := os.Rename(path, newPath); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return filesToMigrate, nil
}
