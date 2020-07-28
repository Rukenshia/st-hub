package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sthub/lib/battle"
)

// TestShip represents basic info on a Warship for testing
type TestShip struct {
	ID   uint64
	Name string
}

// TestIteration is a struct representing a raw iteration cycle
type TestIteration struct {
	ClientVersion string
	IterationName string
	Ships         []TestShip
}

// TestIterationFile represents a on-disk stored
type TestIterationFile struct {
	TestIteration

	filename string

	Battles []*battle.Battle
}

// LoadTestIterationFile loads a test iteration file from disk
func LoadTestIterationFile(configPath, iterationName string) (*TestIterationFile, error) {
	data, err := ioutil.ReadFile(filepath.Join(configPath, fmt.Sprintf("sthub-%s.json", iterationName)))
	if err != nil {
		return nil, err
	}

	var t TestIterationFile
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	t.filename = filepath.Join(configPath, fmt.Sprintf("sthub-%s.json", iterationName))

	return &t, nil
}

// LoadOrCreateIterationFile loads an existing or creates an empty iteration file
//
// Due to the nature of Wargaming's ST program, the testships during an iteration can change at
// any time. To prevent inconsistencies (or forcing a new file while an iteration is still running),
// the Ships on an existing file will always be overwritten.
func LoadOrCreateIterationFile(configPath string, currentIteration *TestIteration) (*TestIterationFile, error) {
	filename := filepath.Join(configPath, fmt.Sprintf("sthub-%s-%s.json", currentIteration.ClientVersion, currentIteration.IterationName))

	var ti *TestIterationFile
	if _, err := os.Stat(filename); err == nil {
		ti, err = LoadTestIterationFile(configPath, fmt.Sprintf("%s-%s", currentIteration.ClientVersion, currentIteration.IterationName))
		if err != nil {
			return nil, err
		}

		// Always override ships as they can change at any time
		ti.Ships = currentIteration.Ships

		log.Printf("%v", ti)
	} else {
		ti = &TestIterationFile{
			TestIteration: *currentIteration,
			filename:      filename,
		}
	}

	data, err := json.Marshal(ti)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(filename, data, 0666); err != nil {
		return nil, err
	}

	return ti, nil
}

// Save stores the file on its original path
func (f *TestIterationFile) Save() error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.filename, data, 0666)
}

// HasShip checks whether a ship is part of this test iteration
func (t *TestIteration) HasShip(id uint64) bool {
	for _, s := range t.Ships {
		if s.ID == id {
			return true
		}
	}
	return false
}

// GetShip returns info on a ship if it exists
func (t *TestIteration) GetShip(id uint64) *TestShip {
	for _, s := range t.Ships {
		if s.ID == id {
			return &s
		}
	}
	return nil
}
