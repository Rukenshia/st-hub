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

					if filepath.Base(event.Name) == "battle.start" {
						data, err := ioutil.ReadFile(event.Name)
						if err != nil {
							log.Println(err)
							dialog.Message("Could not read battle start file. Please contact Rukenshia. The program will exit now.").Title("StHub: ERR_BATTLE_FLOW_READ_START").Error()
							os.Exit(1)
						}

						res, err := s.reportBattleStart(data)
						if err != nil {
							if err.Error() == "ERR_NOT_IN_TESTING" {
								if err := ioutil.WriteFile(filepath.Join(filepath.Dir(event.Name), "battle.response"), []byte("ERR_NOT_IN_TESTING"), 0666); err != nil {
									log.Println(err)
									dialog.Message("Could not send battle start request. Please contact Rukenshia").Title("StHub: ERR_BATTLE_FLOW_START_FAIL_WRITE").Error()
									continue
								}
								continue
							} else {
								log.Println(err)
								dialog.Message("Could not send battle start request. Please contact Rukenshia").Title("StHub: ERR_BATTLE_FLOW_START_SEND").Error()
								continue
							}
						}

						// Write battle result
						data, err = json.Marshal(res)
						if err != nil {
							log.Println(err)
							dialog.Message("Could not marshal data. Please contact Rukenshia").Title("StHub: ERR_BATTLE_FLOW_START_MARSHAL").Error()
							continue
						}

						if err := ioutil.WriteFile(filepath.Join(filepath.Dir(event.Name), "battle.response"), data, 0666); err != nil {
							log.Println(err)
							dialog.Message("Could not send battle start request. Please contact Rukenshia").Title("StHub: ERR_BATTLE_FLOW_START_WRITE").Error()
							continue
						}
					} else if filepath.Base(event.Name) == "battle.end" {
						data, err := ioutil.ReadFile(event.Name)
						if err != nil {
							log.Println(err)
							dialog.Message("Could not read battle end file. Please contact Rukenshia. The program will exit now.").Title("StHub: ERR_BATTLE_FLOW_READ_END").Error()
							os.Exit(1)
						}

						s.reportBattleEnd(data)
					} else {
						log.Println("scraper: untracked file changed")
					}
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
