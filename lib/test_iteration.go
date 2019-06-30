package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sthub/lib/battle"
)

// TestIterationFile represents a on-disk stored
type TestIterationFile struct {
	ClientVersion string
	IterationName string

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

	return &t, nil
}
