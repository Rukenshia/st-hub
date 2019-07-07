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
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/levigross/grequests"
	"github.com/sqweek/dialog"
)

// Scraper is a module that keeps track of changes in the World of Warships
// installation directory. More specifically, it looks for changes in files
// of the clientside mod to detect new battle information.
type Scraper struct {
	wowsPath        string
	c               *lib.TestController
	rejectedBattles map[string]error
}

// New creates a new Scraper
func New(wowsPath string, c *lib.TestController) *Scraper {
	return &Scraper{c: c, wowsPath: wowsPath, rejectedBattles: make(map[string]error, 0)}
}

// Start starts the scraping process
func (s *Scraper) Start(clientVersion string) error {
	dir := filepath.Join(s.wowsPath, "res_mods", clientVersion, "PnFMods", "StHub", "api")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("The modification is not installed correctly to %s", dir)
	}

	// Parse existing files
	if err := s.loadCurrentFiles(dir); err != nil {
		return fmt.Errorf("Could not load existing files with scraper: %v", err)
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

					if info.IsBattleStart() {
						ab := s.c.GetActiveBattle()
						if ab != nil && ab.Timestamp == info.Timestamp {
							log.Println("scraper: active battle is already the fs event")
							continue
						}

						if _, err := s.reportBattleStart(info); err != nil {
							if err.Error() == "ERR_NOT_IN_TESTING" {
								s.rejectedBattles[info.Timestamp] = err
								continue
							}
							log.Printf("scraper: Error reporting battle: %v", err)
							continue
						}
					} else if info.IsBattleEnd() || info.IsBattleQuit() {
						if _, ok := s.rejectedBattles[info.Timestamp]; ok {
							log.Printf("scraper: ignoring previously rejected battle %s", info.Timestamp)
							delete(s.rejectedBattles, info.Timestamp)
							continue
						}
						ab := s.c.GetActiveBattle()
						if ab == nil {
							log.Println("scraper: no active battle found, ignoring battle result")
							continue
						}

						if ab.Timestamp != info.Timestamp {
							log.Printf("scraper: battle end is for timestamp %s, but active battle is %s. ignoring battle", info.Timestamp, ab.Timestamp)
							dialog.Message("It seems like there are multiple active battles, please contact Rukenshia to fix this bug").Title("StHub: ERR_BATTLE_FLOW_TIMESTAMP_MISMATCH").Error()
							continue
						}

						if err := s.reportBattleEnd(ab, info); err != nil {
							log.Printf("scraper: could not report battle end: %v", err)
							dialog.Message("Could not report battle end to the API. Please contact Rukenshia").Title("StHub: ERR_BATTLE_FLOW_REPORT_END").Error()
							continue
						}

						if err := os.Remove(event.Name); err != nil {
							log.Printf("scraper: could not remove file: %v", err)
							dialog.Message("Could not remove a battle file. Please contact Rukenshia").Title("StHub: ERR_BATTLE_FLOW_REPORT_END_RM").Error()
						}
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

func (s *Scraper) loadCurrentFiles(dir string) error {
	log.Printf("scraper: loading current files")
	return filepath.Walk(dir, func(path string, stat os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// remove 0.1.x files
		filename := filepath.Base(path)
		if filename == "battle.start" || filename == "battle.end" || filename == "battle.response" {
			if err := os.Remove(path); err != nil {
				log.Printf("scraper: could not remove 0.1.x file: %v", err)
				dialog.Message("Could not remove a 0.1.x file. Please contact Rukenshia").Title("StHub: ERR_SCRAPER_LOAD_REPORT_01x_RM").Error()
			}
			return nil
		}

		if !strings.Contains(filepath.Base(path), "battle.") {
			return nil
		}

		log.Printf("scraper: found %s", path)

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		info := new(ModBattleInfo)
		if err := json.Unmarshal(data, info); err != nil {
			return err
		}

		if info.IsBattleStart() {
			ab := s.c.GetActiveBattle()
			if ab != nil && ab.Timestamp == info.Timestamp {
				return nil
			}

			log.Printf("scraper: reporting battle start for %s", info.Timestamp)

			if _, err := s.reportBattleStart(info); err != nil {
				if err.Error() == "ERR_NOT_IN_TESTING" {
					s.rejectedBattles[info.Timestamp] = err
					return nil
				}
				log.Printf("scraper: Error reporting battle: %v", err)
				return err
			}
		} else if info.IsBattleEnd() || info.IsBattleQuit() {
			if _, ok := s.rejectedBattles[info.Timestamp]; ok {
				log.Printf("scraper: ignoring previously rejected battle %s", info.Timestamp)
				delete(s.rejectedBattles, info.Timestamp)
				return nil
			}
			ab := s.c.GetActiveBattle()
			if ab == nil {
				log.Printf("scraper: found possible battle that ended while sthub was not running. registering as active battle")

				if _, err := s.reportBattleStart(info); err != nil {
					if err.Error() == "ERR_NOT_IN_TESTING" {
						log.Printf("scraper: battle not test ship, deleting file")

						if err := os.Remove(path); err != nil {
							log.Printf("scraper: could not remove non-testship file: %v", err)
							dialog.Message("Could not remove a non-testship file. Please contact Rukenshia").Title("StHub: ERR_SCRAPER_LOAD_REPORT_ERR_RM").Error()
						}
						return nil
					}
					log.Printf("scraper: Error reporting battle: %v", err)
					return err
				}

				ab = s.c.GetActiveBattle()
			}

			if ab.Timestamp != info.Timestamp {
				log.Printf("scraper: battle end is for timestamp %s, but active battle is %s. ignoring battle", info.Timestamp, ab.Timestamp)
				dialog.Message("It seems like there are multiple active battles, please contact Rukenshia to fix this bug").Title("StHub: ERR_SCRAPER_LOAD_TIMESTAMP_MISMATCH").Error()
				return err
			}

			log.Printf("scraper: reporting battle end for %s", info.Timestamp)
			if err := s.reportBattleEnd(ab, info); err != nil {
				log.Printf("scraper: could not report battle end: %v", err)
				dialog.Message("Could not report battle end to the API. Please contact Rukenshia").Title("StHub: ERR_SCRAPER_LOAD_REPORT_END").Error()
				return err
			}

			if err := os.Remove(path); err != nil {
				log.Printf("scraper: could not remove file: %v", err)
				dialog.Message("Could not remove a battle file. Please contact Rukenshia").Title("StHub: ERR_SCRAPER_LOAD_REPORT_END_RM").Error()
			}
		}
		return nil
	})
}

func (s *Scraper) reportBattleStart(info *ModBattleInfo) (*battle.Battle, error) {
	data, err := json.Marshal(info.ToBattleStartRequest())
	if err != nil {
		return nil, err
	}

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

func (s *Scraper) reportBattleEnd(activeBattle *battle.Battle, info *ModBattleInfo) error {
	// Construct a new Battle object
	b := &battle.Battle{
		ID:         activeBattle.ID,
		Statistics: info.GetStatistics(),
		FinishedAt: activeBattle.FinishedAt,
		StartedAt:  activeBattle.StartedAt,
		ShipID:     activeBattle.ShipID,
		ShipName:   activeBattle.ShipName,
		Status:     info.Status,
		Timestamp:  activeBattle.Timestamp,
	}

	data, err := json.Marshal(b)
	if err != nil {
		return nil
	}

	res, err := grequests.Post("http://localhost:1323/iterations/current/battles/active", &grequests.RequestOptions{
		RequestBody: bytes.NewBuffer(data),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(res.String())
	}
	return nil
}
