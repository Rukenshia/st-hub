package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sthub/lib"
	"sthub/lib/scraper"
	"time"

	"github.com/levigross/grequests"
	"github.com/sqweek/dialog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	//	"gopkg.in/toast.v1"
)

func main() {
	// Initialise the application
	cfg, err := initApp()
	if err != nil {
		dialog.Message("%s", "Could not load configuration.").Title("StHub").Error()
		os.Exit(1)
	}

	// Find current test iteration
	res, err := grequests.Get("https://hv59yay1u3.execute-api.eu-central-1.amazonaws.com/live/iteration/current", nil)
	if err != nil {
		dialog.Message("%s: %v. %s", "Could not retrieve current test iteration", err, "Please contact Rukenshia.").Title("StHub: RETR_CURRENT_ITER").Error()
		log.Fatalln(err)
	}

	currentIteration := new(lib.TestIteration)
	if err := res.JSON(currentIteration); err != nil {
		dialog.Message("%s: %v. %s", "Could not parse current test iteration", err, "Please contact Rukenshia.").Title("StHub: PARSE_CURRENT_ITER").Error()
		log.Fatalln(err)
	}

	testController, err := lib.NewTestController(currentIteration)
	if err != nil {
		dialog.Message("%s: %v. %s", "Could not create API Controller", err, "Please contact Rukenshia.").Title("StHub: LC_CURRENT_ITER").Error()
		log.Fatalln(err)
	}

	scraper := scraper.New(cfg.WowsPath, testController)

	if err := scraper.Start(currentIteration.ClientVersion); err != nil {
		dialog.Message("%s: %v", "Could not start waiting for info", err).Title("StHub: ERR_SCRAPER_START").Error()
		log.Fatalln(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5000", "http://100.115.92.205:5000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Frontend
	e.Static("/frontend", "frontend/public")

	testController.RegisterRoutes(e.Group("/iterations"))

	go func() {
		// Wait and open the browser
		timer := time.NewTimer(100 * time.Millisecond)

		<-timer.C

		// if err := exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:1323/frontend/index.html").Start(); err != nil {
		// 	log.Fatal("Could not open browser")
		// }
	}()

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func initApp() (*Config, error) {
	// Load local config
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, fmt.Errorf("could not find executable path: %v", err)
	}

	var config Config

	data, err := ioutil.ReadFile(filepath.Join(execDir, "sthub-config.json"))
	if err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			dialog.Message("%s", "Your sthub-config.json is corrupted. Please contact Rukenshia").Title("Fatal error").Error()
			os.Exit(1)
		}
	} else {
		if os.IsNotExist(err) {
			// Ask for WoWS Directory
			config = Config{
				WowsPath: "/invalid/path",
			}
		} else {
			return nil, fmt.Errorf("could not read config file: %v", err)
		}
	}

	// Check if the WoWS path is set and still exists
	if _, err := os.Stat(config.WowsPath); os.IsNotExist(err) {
		dialog.Message("%s", "You need to pick the installation directory of your World of Warships client.").
			Title("StHub Setup").
			Info()
		dir, err := dialog.Directory().Title("Choose World of Warships Installation Path").Browse()
		if err != nil {
			log.Fatalf("Could not pick directory: %v", err)
		}

		config.WowsPath = dir
	}

	data, err = json.Marshal(&config)
	if err != nil {
		log.Fatalf("Could not marshal before save: %v", err)
	}
	if err := ioutil.WriteFile(filepath.Join(execDir, "sthub-config.json"), data, 0666); err != nil {
		log.Fatalf("Could not write config: %v", err)
	}

	// Check if the modification exists

	return &config, nil
}
