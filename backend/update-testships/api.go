package main

import (
	"fmt"
	"log"
	"sync"

	resty "github.com/go-resty/resty/v2"
)

// WowsAPI is a module providing access to the required
// API methods of the Wargaming API, masking away the application
// id from every request to make it easier to use.
type WowsAPI struct {
	// The application id (not a secret / private token) used to make calls to
	// the WoWS API
	applicationID string

	// The realm specifies which server of the wargaming API should be targeted
	// Can be one of: eu, ru, com, asia
	realm string

	// Resty client for HTTP requests
	client *resty.Client
}

// APIError contains useful information when the Wargaming API
// returned an error code back to us. This struct implements the
// error interface.
type APIError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
	Value   string `json:"value"`
}

// Error returns a formatted error message
func (a APIError) Error() string {
	return fmt.Sprintf("Wargaming API Error: %d - %s. field=%s value=%s", a.Code, a.Message, a.Field, a.Value)
}

// APIResponse is a wrapping object used by the Wargaming API that
// indicates the state of your request and contains various metadata
type APIResponse struct {
	// Can bei either "ok" or "error".
	Status string `json:"status"`

	// The Error struct only contains data if the Status is "error".
	Error *APIError `json:"error"`

	// the Meta struct only contains data if the Status is "ok".
	Meta *struct {
		Count     uint `json:"count"`
		Total     uint `json:"total"`
		Limit     uint `json:"limit"`
		Page      uint `json:"page"`
		PageTotal uint `json:"page_total"`
	} `json:"meta"`

	// Response data
	Data interface{}
}

// ListShipsResponse is the wargaming API response for
// /wows/encyclopedia/ships
type ListShipsResponse struct {
	APIResponse

	Data map[string]*Ship `json:"data"`
}

// A Ship represents a WoWS Warship from the API. This struct
// only contains information relevant to this project.
type Ship struct {
	ShipID    uint   `json:"ship_id"`
	ShipIDStr string `json:"ship_id_str"`

	Name        string `json:"name"`
	Nation      string `json:"nation"`
	Tier        uint   `json:"tier"`
	Type        string `json:"type"`
	Description string `json:"description"`

	Images struct {
		Contour string `json:"contour"`
		Small   string `json:"small"`
		Medium  string `json:"medium"`
		Large   string `json:"large"`
	} `json:"images"`

	PriceCredit uint `json:"price_credit"`
	PriceGold   uint `json:"price_gold"`

	HasDemoProfile bool `json:"has_demo_profile"`
	IsPremium      bool `json:"is_premium"`
	IsSpecial      bool `json:"is_special"`

	// NextShips is an associative array of ShipID <> Experience required
	NextShips map[string]uint `json:"next_ships"`

	ModSlots uint   `json:"mod_slots"`
	Upgrades []uint `json:"upgrades"`

	// TODO
	DefaultProfile interface{} `json:"default_profile"`
	Modules        interface{} `json:"modules"`
	ModulesTree    interface{} `json:"modules_tree"`
}

func (w *WowsAPI) buildURL(format string, a ...interface{}) string {
	return fmt.Sprintf("https://api.worldofwarships.%s/%s", w.realm, fmt.Sprintf(format, a...))
}

// FindShips returns a slice of ships that match any of the given names.
// Globs can be used inside of the name parameter.
func (w *WowsAPI) FindShips(globNames []string) ([]*Ship, error) {
	return nil, nil
}

// GetWarships returns a map of all ship ids and ship objects present in the
// game.
//
// This function will query all warships in the game, so it can be
// quite slow as multiple requests are required to retrieve all
// the information. This is also why this function accepts and returns
// slices.
func (w *WowsAPI) GetWarships() (map[string]*Ship, error) {
	// Make the first request to determine how many pages we need, afterwards
	// we can parallelise our efforts. Since the result is a map anyway, sorting
	// does not matter.
	resp, err := w.getPage(w.buildURL("wows/encyclopedia/ships/"), 1, &ListShipsResponse{})
	if err != nil {
		return nil, err
	}

	// TODO: is there a way to remove the type assertion here, into getPage?
	// maybe make the last param into a variable that we can then just use?
	firstPage := resp.(*ListShipsResponse)

	if firstPage.Error != nil {
		return nil, firstPage.Error
	}

	responses := []*ListShipsResponse{firstPage}

	var wg sync.WaitGroup
	for i := uint(2); i < firstPage.Meta.PageTotal+1; i++ {
		wg.Add(1)
		go func(pageNo uint) {
			defer wg.Done()

			resp, err := w.getPage(w.buildURL("wows/encyclopedia/ships/"), pageNo, &ListShipsResponse{})
			if err != nil {
				log.Printf("Could not get page: %v", err)
			}

			responses = append(responses, resp.(*ListShipsResponse))
		}(i)
	}

	wg.Wait()

	ships := make(map[string]*Ship, 0)

	for _, page := range responses {
		if page == nil {
			return nil, fmt.Errorf("invalid page retrieved")
		}

		if page.Error != nil {
			return nil, page.Error
		}

		for id, ship := range page.Data {
			ships[id] = ship
		}
	}

	return ships, nil
}

func (w *WowsAPI) getPage(url string, pageNo uint, response interface{}) (interface{}, error) {
	resp, err := w.client.R().
		SetQueryParam("page_no", fmt.Sprintf("%d", pageNo)).
		SetQueryParam("application_id", w.applicationID).
		SetResult(response).
		Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Result(), err
}
