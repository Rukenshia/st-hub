package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
	"github.com/sqweek/dialog"
)

// checkRequiredFiles checks whether the game modification files exist yet. The first
// return value is a flag whether the mod needs to be installed, the error will contain
// the specific file that failed this function.
func checkRequiredFiles(gamePath string) (bool, error) {
	if _, err := os.Stat(filepath.Join(gamePath, "PnFModsLoader.py")); os.IsNotExist(err) {
		return false, errors.New("PnFModsLoader.py does not exist")
	}
	if _, err := os.Stat(filepath.Join(gamePath, "PnFMods")); os.IsNotExist(err) {
		return false, errors.New("PnFMods does not exist")
	}
	if _, err := os.Stat(filepath.Join(gamePath, "PnFMods", "StHub")); os.IsNotExist(err) {
		return false, errors.New("StHub does not exist")
	}
	if _, err := os.Stat(filepath.Join(gamePath, "PnFMods", "StHub", "Main.py")); os.IsNotExist(err) {
		return false, errors.New("Main.py does not exist")
	}
	if _, err := os.Stat(filepath.Join(gamePath, "PnFMods", "StHub", "api")); os.IsNotExist(err) {
		return false, errors.New("api does not exist")
	}

	return true, nil
}

// installGameMod installs the game modification into the game directory of the user
func installGameMod(wowsPath string, clientVersion string) (string, error) {
	if exists, err := checkRequiredFiles(filepath.Join(wowsPath, "res_mods", clientVersion)); !exists || err != nil {
		dialog.Message("%s", "The game modification will now be added to your client. If you have World of Warships running, you will have to restart it.").
			Title("StHub Setup").
			Info()
	}

	// Always install the mod
	box := rice.MustFindBox("mod")

	if err := os.MkdirAll(filepath.Join(wowsPath, "res_mods", clientVersion, "PnFMods", "StHub", "api"), 0666); err != nil {
		return fmt.Sprintf("Could not create required directories for the mod: %v", err), err
	}

	pnfModsLoader, err := box.Bytes("PnFModsLoader.py")
	if err != nil {
		return fmt.Sprintf("Could not load PnFModsLoader.py: %v", err), err
	}
	sthubMain, err := box.Bytes("PnFMods/StHub/Main.py")
	if err != nil {
		return fmt.Sprintf("Could not load Main.py: %v", err), err
	}

	if err := ioutil.WriteFile(filepath.Join(wowsPath, "res_mods", clientVersion, "PnFModsLoader.py"), pnfModsLoader, 0666); err != nil {
		return fmt.Sprintf("Could not write PnFModsLoader.py: %v", err), err
	}

	if err := ioutil.WriteFile(filepath.Join(wowsPath, "res_mods", clientVersion, "PnFMods", "StHub", "Main.py"), sthubMain, 0666); err != nil {
		return fmt.Sprintf("Could not write Main.py: %v", err), err
	}

	return "", nil
}
