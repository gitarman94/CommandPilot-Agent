package main

import (
	"log"
	"time"
)

func main() {
	cfg := LoadConfig()

	log.Printf("pilot-agent starting")
	log.Printf("server=%s", cfg.ServerURL)

	go StartTelemetryLoop(cfg)
	go StartCommandLoop(cfg)
	go StartUpdateLoop(cfg)

	select {}
}

func sleepSeconds(v int) {
	time.Sleep(time.Duration(v) * time.Second)
}