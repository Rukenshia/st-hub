package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/sqweek/dialog"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v30/github"
	"github.com/inconshreveable/go-update"
)

// Tries to pull the latest version of sthub
func selfUpdate() error {
	log.Printf("Checking for updates")

	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "Rukenshia", "st-hub")
	if err != nil {
		return err
	}

	version, err := semver.NewVersion(release.GetTagName())
	if err != nil {
		return err
	}

	if VERSION.GreaterThan(version) {
		log.Printf("Current version is latest, skip download")
		return nil
	}

	log.Printf("Update %v available, asking for user input", version)

	yes := dialog.Message("A new update, %s, is available. Do you want to download it in the background now? You'll be notified once it is installed.'", version.String()).Title("StHub: New version available").YesNo()

	if !yes {
		log.Printf("User aborted update")
		return nil
	}

	//if runtime.GOOS != "windows" {
	//log.Printf("OS is not windows, cancelling update")
	//dialog.Message("Unsupported operating system. Please download manually.", version.String()).Title("StHub: New version available").Info()
	//return nil
	//}

	log.Printf("Searching update assets", version)

	// Find the release asset and download it
	var url string

	for _, asset := range release.Assets {
		if asset.GetName() == "sthub.exe" {
			url = asset.GetBrowserDownloadURL()
			break
		}
	}

	if url == "" {
		log.Printf("sthub.exe asset not found")
		return fmt.Errorf("Could not find sthub.exe download url")
	}

	log.Printf("Starting goroutine for update")

	go func() {
		log.Printf("Downloading update %v", version)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error downloading update: %v", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("Applying update")

		if err := update.Apply(resp.Body, update.Options{}); err != nil {
			log.Printf("Error applying update: %v", err)
			dialog.Message("An error occured while trying to update: %v. Please contact Rukenshia", err).Title("StHub: ERR_UPDATE_APPLY").Error()
			return
		}

		log.Printf("Update applied successfully")

		dialog.Message("Version %s is now installed. To apply the update, you need to restart StHub.", version.String()).Title("StHub: Update finsihed").Info()
	}()

	return nil
}
