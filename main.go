package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sthub/lib"
	"sthub/lib/scraper"

	"github.com/levigross/grequests"
	"github.com/sqweek/dialog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"

	"github.com/Masterminds/semver"

	rice "github.com/GeertJohan/go.rice"
)

// VERSION represents the current version of StHub (this component)
var VERSION = semver.MustParse("0.5.0")

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

	// unpack logo
	iconPath, err := unpackLogo()
	defer os.Remove(iconPath)

	// Initialize astilectron
	var a, _ = astilectron.New(
		log.New(os.Stderr, "", 0),
		astilectron.Options{
			AppName:            "StHub",
			AppIconDefaultPath: iconPath,
			// BaseDirectoryPath: "<where you want the provisioner to install the dependencies>",
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
		var w, _ = a.NewWindow("https://sthub.in.fkn.space", &astilectron.WindowOptions{
			Center: astikit.BoolPtr(true),
			Height: astikit.IntPtr(800),
			Width:  astikit.IntPtr(600),
		})
		w.Create()

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

// Returns whether a local config file exists
func hasLocalConfig() bool {
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return false
	}

	if _, err := os.Stat(path.Join(execDir, "sthub-config.json")); os.IsNotExist(err) {
		return false
	}

	return true
}

func start(done chan bool, currentIteration *lib.TestIteration) {
	// Initialise the application
	cfg, err := initApp(currentIteration)
	if err != nil {
		dialog.Message("%s", "Could not load configuration.").Title("StHub").Error()
		os.Exit(1)
	}

	testController, err := lib.NewTestController(currentIteration)
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

	if err := checkRequiredFiles(filepath.Join(config.WowsPath, "res_mods", currentIteration.ClientVersion)); err != nil {
		dialog.Message("%s", "The game modification will now be added to your client. If you have World of Warships running, you will have to restart it.").
			Title("StHub Setup").
			Info()
	}

	// Always install the mod
	box := rice.MustFindBox("mod")

	if err := os.MkdirAll(filepath.Join(config.WowsPath, "res_mods", currentIteration.ClientVersion, "PnFMods", "StHub", "api"), 0666); err != nil {
		dialog.Message("%s: %v", "Could not create required directories for the mod", err).
			Title("StHub Setup").
			Error()
		os.Exit(1)
	}

	pnfModsLoader, err := box.Bytes("PnFModsLoader.py")
	if err != nil {
		dialog.Message("%s: %v", "Could not load PnFModsLoader.py", err).
			Title("StHub Setup").
			Error()
		os.Exit(1)
	}
	sthubMain, err := box.Bytes("PnFMods/StHub/Main.py")
	if err != nil {
		dialog.Message("%s: %v", "Could not load StHub/Main.py", err).
			Title("StHub Setup").
			Error()
		os.Exit(1)
	}

	if err := ioutil.WriteFile(filepath.Join(config.WowsPath, "res_mods", currentIteration.ClientVersion, "PnFModsLoader.py"), pnfModsLoader, 0666); err != nil {
		dialog.Message("%s: %v", "Could not write PnFModsLoader.py", err).
			Title("StHub Setup").
			Error()
		os.Exit(1)
	}

	if err := ioutil.WriteFile(filepath.Join(config.WowsPath, "res_mods", currentIteration.ClientVersion, "PnFMods", "StHub", "Main.py"), sthubMain, 0666); err != nil {
		dialog.Message("%s: %v", "Could not write StHub/Main.py", err).
			Title("StHub Setup").
			Error()
		os.Exit(1)
	}

	return &config, nil
}

func checkRequiredFiles(gamePath string) error {
	if _, err := os.Stat(filepath.Join(gamePath, "PnFModsLoader.py")); os.IsNotExist(err) {
		return errors.New("PnFModsLoader.py does not exist")
	}
	if _, err := os.Stat(filepath.Join(gamePath, "PnFMods")); os.IsNotExist(err) {
		return errors.New("PnFMods does not exist")
	}
	if _, err := os.Stat(filepath.Join(gamePath, "PnFMods", "StHub")); os.IsNotExist(err) {
		return errors.New("StHub does not exist")
	}
	if _, err := os.Stat(filepath.Join(gamePath, "PnFMods", "StHub", "Main.py")); os.IsNotExist(err) {
		return errors.New("Main.py does not exist")
	}
	if _, err := os.Stat(filepath.Join(gamePath, "PnFMods", "StHub", "api")); os.IsNotExist(err) {
		return errors.New("api does not exist")
	}

	return nil
}
