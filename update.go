package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type UpdateInfo struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

func StartUpdateLoop(cfg Config) {
	for {
		checkForUpdate(cfg)
		sleepSeconds(300)
	}
}

func checkForUpdate(cfg Config) {
	resp, err := http.Get(
		cfg.ServerURL + "/api/agent/update?version=" + cfg.Version,
	)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}

	var update UpdateInfo

	if err := json.NewDecoder(resp.Body).Decode(&update); err != nil {
		return
	}

	if update.Version == "" || update.Version == cfg.Version {
		return
	}

	log.Printf("updating to %s", update.Version)

	downloadAndInstall(update.URL)
}

func downloadAndInstall(url string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	tmp := "/tmp/pilot-agent-update.zip"

	f, err := os.Create(tmp)
	if err != nil {
		return
	}

	defer f.Close()

	io.Copy(f, resp.Body)

	log.Printf("update downloaded: %s", tmp)

	// hook installer here
}