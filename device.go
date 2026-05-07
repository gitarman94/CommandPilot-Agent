package main

import (
	"os"

	"github.com/google/uuid"
)

const agentIDFile = "/etc/commandpilot-agent-id"

func loadOrCreateAgentID() string {
	data, err := os.ReadFile(agentIDFile)
	if err == nil && len(data) > 0 {
		return string(data)
	}

	id := uuid.NewString()

	_ = os.WriteFile(agentIDFile, []byte(id), 0600)

	return id
}