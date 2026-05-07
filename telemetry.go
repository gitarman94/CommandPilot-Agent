package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func StartTelemetryLoop(cfg Config) {
	for {
		sendTelemetry(cfg)
		sleepSeconds(30)
	}
}

func sendTelemetry(cfg Config) {
	info := GetSystemInfo(cfg)

	body, _ := json.Marshal(info)

	req, err := http.NewRequest(
		http.MethodPost,
		cfg.ServerURL+"/api/agent/checkin",
		bytes.NewBuffer(body),
	)

	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()
}