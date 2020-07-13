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
func (s *Scraper) Start(clientPath string) error {
	dir := filepath.Join(s.wowsPath, "bin", clientPath, "res_mods", "PnFMods", "StHub", "api")
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

					if strings.Contains(filepath.Base(event.Name), "results.") {
						if err := s.handleResultsFile(event.Name); err != nil {
							log.Printf("scraper: could not handle results file: %v", err)
							dialog.Message("Could not process battle results. Please contact Rukenshia").Title("StHub: ERR_SCRAPER_HANDLE_RESULTS_FILE").Error()
							continue
						}
					}

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

							if err := os.Remove(event.Name); err != nil {
								log.Printf("scraper: could not remove file: %v", err)
								dialog.Message("Could not remove a rejected battle file. Please contact Rukenshia").Title("StHub: ERR_BATTLE_FLOW_REPORT_REJECTED_RM").Error()
							}
							continue
						}
						ab := s.c.GetActiveBattle()
						if ab == nil {
							log.Println("scraper: no active battle found, ignoring battle result")
							continue
						}

						if ab.Timestamp != info.Timestamp {
							log.Printf("scraper: battle end is for timestamp %s, but active battle is %s", info.Timestamp, ab.Timestamp)
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

		if strings.Contains(filepath.Base(path), "results.") {
			if err := s.handleResultsFile(path); err != nil {
				log.Printf("scraper: could not handle results file: %v", err)
				dialog.Message("Could not process battle results. Please contact Rukenshia").Title("StHub: ERR_SCRAPER_HANDLE_RESULTS_FILE").Error()
				return err
			}
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
				if err := os.Remove(path); err != nil {
					log.Printf("scraper: could not file: %v", err)
					dialog.Message("Could not remove a non-testship file. Please contact Rukenshia").Title("StHub: ERR_SCRAPER_LOAD_START_RM_ALREADY_ACTIVE").Error()
				}
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
			if err := os.Remove(path); err != nil {
				log.Printf("scraper: could not remove start file: %v", err)
				dialog.Message("Could not remove a start battle file. Please contact Rukenshia").Title("StHub: ERR_SCRAPER_LOAD_START_RM").Error()
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
				if ab.Timestamp < info.Timestamp {
					log.Printf("the current active battle (id %s, ts %s) is older than the loaded file %s. force abandoning current active battle", ab.ID, ab.Timestamp, info.Timestamp)

					if err := s.abandonActiveBattle(ab); err != nil {
						log.Printf("scraper: could not abandon active battle: %v", err)
						dialog.Message("There was an error abandoning an old battle. Please contact Rukenshia").Title("StHub: ERR_SCRAPER_LOAD_REPORT_END_ABANDON_ACTIVE").Error()
						return err
					}

					if _, err := s.reportBattleStart(info); err != nil {
						if err.Error() == "ERR_NOT_IN_TESTING" {
							log.Printf("scraper: ship not in testing, aborting")
							if err := os.Remove(path); err != nil {
								log.Printf("scraper: could not remove file: %v", err)
								dialog.Message("Could not remove a battle file. Please contact Rukenshia").Title("StHub: ERR_SCRAPER_LOAD_REPORT_END_RM").Error()
							}
							return nil
						}
						log.Printf("scraper: could not report battle start after force-abandon: %v", err)
						dialog.Message("There was an error reporting a battle, please contact Rukenshia").Title("StHub: ERR_SCRAPER_LOAD_REPORT_END_ABANDONED_START").Error()
						return err
					}
				} else {
					log.Printf("scraper: battle end is for timestamp %s, but active battle is %s", info.Timestamp, ab.Timestamp)
					dialog.Message("It seems like there are multiple active battles, please contact Rukenshia to fix this bug").Title("StHub: ERR_SCRAPER_LOAD_TIMESTAMP_MISMATCH_NEWER").Error()
				}
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

func (s *Scraper) handleResultsFile(path string) error {
	log.Printf("scraper: found %s", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	info := new(battle.Results)
	if err := json.Unmarshal(data, info); err != nil {
		return err
	}

	if _, ok := s.rejectedBattles[info.Timestamp]; ok {
		log.Printf("scraper: ignoring previously rejected battle %s", info.Timestamp)
		delete(s.rejectedBattles, info.Timestamp)

		if err := os.Remove(path); err != nil {
			log.Printf("scraper: could not remove file: %v", err)
			dialog.Message("Could not remove a rejected battle file. Please contact Rukenshia").Title("StHub: ERR_BATTLE_FLOW_REPORT_REJECTED_RM").Error()
		}
		return nil
	}

	// Find the battle with the given timestamp
	for _, b := range s.c.GetCurrentIterationRaw().Battles {
		if b.Timestamp == info.Timestamp {
			log.Printf("scraper: found battle with correct timestamp, adding results")

			b.Results = info

			if err := os.Remove(path); err != nil {
				return err
			}

			return s.c.SaveCurrentIteration()
		}
	}

	log.Printf("scraper: could not map results to battle. ignoring file")

	// Ignore the file
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
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

func (s *Scraper) abandonActiveBattle(activeBattle *battle.Battle) error {
	// Construct a new Battle object
	b := &battle.Battle{
		ID:         activeBattle.ID,
		Statistics: activeBattle.Statistics,
		FinishedAt: activeBattle.FinishedAt,
		StartedAt:  activeBattle.StartedAt,
		ShipID:     activeBattle.ShipID,
		ShipName:   activeBattle.ShipName,
		Status:     "abandoned",
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
