package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sthub/lib"
	"sthub/lib/scraper"

	"github.com/levigross/grequests"
	"github.com/rs/xid"
	"github.com/sqweek/dialog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"

	"github.com/Masterminds/semver"

	rice "github.com/GeertJohan/go.rice"
)

// VERSION represents the current version of StHub (this component)
var VERSION = semver.MustParse("0.7.0")

func main() {
	f, err := setupLogger()
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Check for new version
	if err := selfUpdate(); err != nil {
		dialog.Message("%s: %v. %s", "Could not check for updates", err, "Please contact Rukenshia.").Title("StHub: ERR_SELF_UPDATE").Error()
		log.Println(err)
	}

	migratedFiles, err := Migrate063ConfigFiles("./")
	if err != nil {
		log.Fatalf("Could not migrate files: %v", err)
	}
	log.Printf("Migrated %d files", len(migratedFiles))

	// unpack logo
	iconPath, err := unpackLogo()
	defer os.Remove(iconPath)

	// Initialize astilectron
	var a, _ = astilectron.New(
		log.New(os.Stderr, "", 0),
		astilectron.Options{
			AppName:            "StHub",
			AppIconDefaultPath: iconPath,
		})
	defer a.Close()

	// Start astilectron
	a.Start()

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

	done := make(chan bool)

	if !hasLocalConfig() {
		var w, _ = a.NewWindow("https://sthub.in.fkn.space/#setup", &astilectron.WindowOptions{
			Center: astikit.BoolPtr(true),
			Height: astikit.IntPtr(800),
			Width:  astikit.IntPtr(600),
		})

		w.OnMessage(func(m *astilectron.EventMessage) interface{} {
			// Unmarshal
			var s string
			m.Unmarshal(&s)

			// Process message
			if s == "connect" {
				return "ok"
			} else if s == "selectGameDir" {
				cfg, err := initApp(currentIteration)
				if err != nil {
					dialog.Message("%s", "Could not load configuration.").Title("StHub").Error()
					os.Exit(1)
				}

				start(done, currentIteration)
				return cfg.WowsPath
			}
			return nil
		})

		w.Create()
	} else {
		var w, _ = a.NewWindow(fmt.Sprintf("https://sthub.in.fkn.space?n=%d", rand.Int31()), &astilectron.WindowOptions{
			Center: astikit.BoolPtr(true),
			Height: astikit.IntPtr(800),
			Width:  astikit.IntPtr(600),
		})
		w.Create()
		go func() {
			w.Session.ClearCache()
		}()

		start(done, currentIteration)
	}

	go func() {
		a.Wait()
		done <- true
	}()

	<-done
}

// setupLogger opens the `sthub.log` file and sets go's standard logger output to it.
func setupLogger() (*os.File, error) {
	f, err := os.OpenFile("sthub.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return nil, err
	}

	log.SetOutput(f)

	return f, nil
}

// unpackLogo retrieves the embedded data from go.rice and writes it to a temporary file, so that astilectron can use it
func unpackLogo() (string, error) {
	tmpfile, err := ioutil.TempFile("", "sthub-logo")
	if err != nil {
		return "", err
	}

	assets := rice.MustFindBox("assets")
	icon, err := assets.Bytes("logo.png")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write(icon); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func start(done chan bool, currentIteration *lib.TestIteration) {
	// Initialise the application
	cfg, err := initApp(currentIteration)
	if err != nil {
		dialog.Message("%s", "Could not load configuration.").Title("StHub").Error()
		os.Exit(1)
	}

	configPath, err := GetConfigPath()
	if err != nil {
		dialog.Message("Could not find a configuration path, please contact Rukenshia").Title("ERR_LOAD_CFG_PATH").Error()
		os.Exit(1)
	}

	testController, err := lib.NewTestController(configPath, currentIteration)
	if err != nil {
		dialog.Message("%s: %v. %s", "Could not create API Controller", err, "Please contact Rukenshia.").Title("StHub: LC_CURRENT_ITER").Error()
		log.Fatalln(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5000", "http://100.115.92.205:5000", "https://sthub.in.fkn.space"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	testController.RegisterRoutes(e.Group("/iterations"))
	e.GET("/version", func(c echo.Context) error {
		c.String(200, VERSION.String())
		return nil
	})

	// Start server
	go func() {
		e.Logger.Fatal(e.Start("localhost:1323"))
		done <- true
	}()

	scraper := scraper.New(cfg.WowsPath, testController)

	if err := scraper.Start(currentIteration.ClientVersion); err != nil {
		dialog.Message("%s: %v", "Could not start waiting for info", err).Title("StHub: ERR_SCRAPER_START").Error()
		log.Fatalln(err)
	}
}

func initApp(currentIteration *lib.TestIteration) (*Config, error) {
	// ignore on mac (in-dev)
	if runtime.GOOS == "darwin" {
		return &Config{
			WowsPath: "/tmp",
		}, nil
	}

	config, err := LoadConfigFromDefaultPath()

	if err != nil {
		if os.IsNotExist(err) {
			id := xid.New()
			// Ask for WoWS Directory
			config = &Config{
				WowsPath: "/invalid/path",
				UserID:   &id,
			}
		} else {
			return nil, fmt.Errorf("could not read config file: %v", err)
		}
	}

	// Check if the WoWS path is set and still exists
	if _, err := os.Stat(config.WowsPath); os.IsNotExist(err) {
		dir, err := dialog.Directory().Title("Choose World of Warships Installation Path").Browse()
		if err != nil {
			log.Fatalf("Could not pick directory: %v", err)
		}

		config.WowsPath = dir
	}

	if err := config.Save(); err != nil {
		log.Fatalf("Could not save config: %v", err)
	}

	if message, err := installGameMod(config.WowsPath, currentIteration.ClientVersion); err != nil {
		dialog.Message(message).
			Title("StHub Setup Failure").
			Error()
		os.Exit(1)
	}

	return config, nil
}
