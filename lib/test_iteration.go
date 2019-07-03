package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
func LoadTestIterationFile(iterationName string) (*TestIterationFile, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("sthub-%s.json", iterationName))
	if err != nil {
		return nil, err
	}

	var t TestIterationFile
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	t.filename = fmt.Sprintf("sthub-%s.json", iterationName)

	return &t, nil
}

// LoadOrCreateIterationFile loads or creates an empty iteration file
func LoadOrCreateIterationFile(currentIteration *TestIteration) (*TestIterationFile, error) {
	filename := fmt.Sprintf("sthub-%s-%s.json", currentIteration.ClientVersion, currentIteration.IterationName)

	if _, err := os.Stat(filename); err == nil {
		return LoadTestIterationFile(fmt.Sprintf("%s-%s", currentIteration.ClientVersion, currentIteration.IterationName))
	}

	ti := &TestIterationFile{
		TestIteration: *currentIteration,
		filename:      filename,
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
