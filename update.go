package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type UpdateInfo struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

func StartUpdateLoop(cfg Config) {
	checkForUpdate(cfg)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		checkForUpdate(cfg)
	}
}

func checkForUpdate(cfg Config) {
	endpoint := strings.TrimRight(cfg.ServerURL, "/") + "/api/agent/update/check?version=" + url.QueryEscape(cfg.Version)

	resp, err := http.Get(endpoint)
	if err != nil {
		log.Printf("update check failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("update check returned: %s", resp.Status)
		return
	}

	var update UpdateInfo
	if err := json.NewDecoder(resp.Body).Decode(&update); err != nil {
		log.Printf("update decode failed: %v", err)
		return
	}

	if update.Version == "" || update.Version == cfg.Version || update.URL == "" {
		return
	}

	log.Printf("update available: %s -> %s", cfg.Version, update.Version)
	downloadAndInstall(cfg, update)
}

func downloadAndInstall(cfg Config, update UpdateInfo) {
	target := update.URL
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = strings.TrimRight(cfg.ServerURL, "/") + "/" + strings.TrimLeft(target, "/")
	}

	resp, err := http.Get(target)
	if err != nil {
		log.Printf("update download failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("update download returned: %s", resp.Status)
		return
	}

	tmp := "/tmp/pilot-agent-update.bin"
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("update read failed: %v", err)
		return
	}

	if err := os.WriteFile(tmp, out, 0600); err != nil {
		log.Printf("update write failed: %v", err)
		return
	}

	log.Printf("update downloaded: %s", tmp)
	log.Printf("install hook not yet implemented")
}