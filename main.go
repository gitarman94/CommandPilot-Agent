package main

import (
	"log"
)

func main() {
	cfg := LoadConfig()

	log.Printf("pilot-agent starting")
	log.Printf("server=%s", cfg.ServerURL)
	log.Printf("agent_id=%s", cfg.AgentID)
	log.Printf("version=%s", cfg.Version)

	go StartHeartbeatLoop(cfg)
	go StartUpdateLoop(cfg)

	select {}
}