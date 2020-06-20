package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/rs/xid"
)

// Config holds basic information required to run StHub
type Config struct {
	path string

	// WowsPath is the file path to the game directory
	WowsPath string

	// UserID is a unique id for the user, allowing us to be
	// able to link information later on
	UserID *xid.ID
}

// LoadConfigFromDefaultPath loads the config from the default path
// and returns an error if it does not exist or cannot be loaded.
func LoadConfigFromDefaultPath() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	var config Config

	data, err := ioutil.ReadFile(filepath.Join(configPath, "sthub-config.json"))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Could not unmarshal config: %v", err)
	}

	return &config, nil
}

// InitUser sets the UserID to a new xid if it has not been generated yet
func (c *Config) InitUser() {
	if c.UserID == nil {
		id := xid.New()
		c.UserID = &id
	}
}

// Save writes the config
func (c *Config) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Could not marshal before save: %v", err)
	}
	if err := ioutil.WriteFile(filepath.Join(c.path, "sthub-config.json"), data, 0666); err != nil {
		return fmt.Errorf("Could not write file: %v", err)
	}

	return nil
}

// GetConfigPath returns the path where sthub data should be stored.
// If the path for the "sthub" directory does not exist yet, it will be
// created by this function.
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
