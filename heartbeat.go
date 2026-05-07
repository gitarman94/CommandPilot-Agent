package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func StartHeartbeatLoop(cfg Config) {
	sendHeartbeat(cfg)

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		sendHeartbeat(cfg)
	}
}

func sendHeartbeat(cfg Config) {
	payload := GetSystemInfo(cfg)

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("heartbeat marshal failed: %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(cfg.ServerURL, "/")+"/api/agent/checkin", bytes.NewReader(body))
	if err != nil {
		log.Printf("heartbeat request build failed: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("heartbeat failed: %v", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("heartbeat rejected: status=%s body=%s", resp.Status, strings.TrimSpace(string(respBody)))
		return
	}

	var ack AgentCheckinResponse
	if err := json.Unmarshal(respBody, &ack); err == nil {
		log.Printf("heartbeat ok: device_id=%d approved=%v", ack.DeviceID, ack.Approved)
		return
	}

	log.Printf("heartbeat ok")
}