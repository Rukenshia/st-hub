package lib

import (
	"fmt"

	"github.com/levigross/grequests"
)

// GameVersion is the version of World of Warships, for example "0.6.9.0"
type GameVersion string

// BasicShip is the basic information of a ship retrieved from StHub
type BasicShip struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type gameVersionResponse struct {
	Version GameVersion `json:"version"`
}

// GetGameVersion will retrieve the current wows game version from StHub
func GetGameVersion() (GameVersion, error) {
	res, err := grequests.Get("https://testhub.in.fkn.space/api/wows/testing/game_version", nil)
	if err != nil {
		return "", fmt.Errorf("Could not get game version: %v", err)
	}

	data := &gameVersionResponse{}
	if err := res.JSON(data); err != nil {
		return "", fmt.Errorf("Could not parse game version response: %v", err)
	}

	return data.Version, nil
}

// GetTestships will retrieve all test ships currently listed on the StHub website
func GetTestships() ([]BasicShip, error) {
	res, err := grequests.Get("https://testhub.in.fkn.space/api/wows/testing/ships", nil)
	if err != nil {
		return nil, fmt.Errorf("Could not get game version: %v", err)
	}

	data := []BasicShip{}
	if err := res.JSON(&data); err != nil {
		return nil, fmt.Errorf("Could not parse game version response: %v", err)
	}

	return data, nil

}
