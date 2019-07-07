package scraper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sthub/lib"
	"sthub/lib/battle"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/levigross/grequests"
	"github.com/sqweek/dialog"
)

// Scraper is a module that keeps track of changes in the World of Warships
// installation directory. More specifically, it looks for changes in files
// of the clientside mod to detect new battle information.
type Scraper struct {
	wowsPath string
	c        *lib.TestController
}

// New creates a new Scraper
func New(wowsPath string, c *lib.TestController) *Scraper {
	return &Scraper{c: c, wowsPath: wowsPath}
}

// Start starts the scraping process
func (s *Scraper) Start(clientVersion string) error {
	dir := filepath.Join(s.wowsPath, "res_mods", clientVersion, "PnFMods", "StHub", "api")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("The modification is not installed correctly to %s", dir)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("Could not create watcher (%v)", err)
	}

	go func() {
		debounce := false
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("scraper: detected change in", event.Name)

					if !strings.Contains(filepath.Base(event.Name), "battle.") {
						log.Printf("skipping unknown file")
						continue
					}

					if debounce {
						continue
					}
					log.Println("scraper: passing change, not debounced")

					debounce = true
					go func() {
						to := time.NewTimer(500 * time.Millisecond)
						<-to.C

						debounce = false
					}()

					data, err := ioutil.ReadFile(event.Name)
					if err != nil {
						log.Println(err)
						dialog.Message("Could not read battle file. Please contact Rukenshia. The program will exit now.").Title("StHub: ERR_BATTLE_FLOW_READ").Error()
						os.Exit(1)
					}

					info := new(ModBattleInfo)
					if err := json.Unmarshal(data, info); err != nil {
						log.Println(err)
						dialog.Message("Could not parse battle file. Please contact Rukenshia. The program will exit now.").Title("StHub: ERR_BATTLE_FLOW_PARSE").Error()
						os.Exit(1)
					}

					// Decide what to send
					// Send to API
					// Handle NOT_IN_TESTING for battle end (store rejected timestamps, clear on end?)
					// Handle active battle errors?
					// Walk fs on startup
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Scraper) reportBattleStart(data []byte) (*battle.Battle, error) {
	res, err := grequests.Post("http://localhost:1323/iterations/current/battles", &grequests.RequestOptions{
		RequestBody: bytes.NewBuffer(data),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New(res.String())
	}

	battle := new(battle.Battle)
	if err := res.JSON(battle); err != nil {
		return nil, err
	}
	return battle, nil
}

func (s *Scraper) reportBattleEnd(data []byte) error {
	_, err := grequests.Post("http://localhost:1323/iterations/current/battles/active", &grequests.RequestOptions{
		RequestBody: bytes.NewBuffer(data),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return err
	}
	return nil
}
