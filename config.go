package main

import (
	"os"
)

type Config struct {
	ServerURL string
	AgentID   string
	Version   string
}

func LoadConfig() Config {
	server := os.Getenv("COMMANDPILOT_SERVER")
	if server == "" {
		server = "http://127.0.0.1:8080"
	}

	id := loadOrCreateAgentID()

	return Config{
		ServerURL: server,
		AgentID:   id,
		Version:   "1.0.0",
	}
}