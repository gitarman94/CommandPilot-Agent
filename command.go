package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
)

type AgentCommand struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
}

type CommandResult struct {
	CommandID int    `json:"command_id"`
	Output    string `json:"output"`
	Success   bool   `json:"success"`
}

func StartCommandLoop(cfg Config) {
	for {
		processCommands(cfg)
		sleepSeconds(10)
	}
}

func processCommands(cfg Config) {
	resp, err := http.Get(
		cfg.ServerURL + "/api/agent/commands?agent_id=" + cfg.AgentID,
	)

	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	var commands []AgentCommand

	if err := json.NewDecoder(resp.Body).Decode(&commands); err != nil {
		return
	}

	for _, c := range commands {
		runCommand(cfg, c)
	}
}

func runCommand(cfg Config, c AgentCommand) {
	cmd := exec.Command("sh", "-c", c.Command)

	out, err := cmd.CombinedOutput()

	result := CommandResult{
		CommandID: c.ID,
		Output:    string(out),
		Success:   err == nil,
	}

	body, _ := json.Marshal(result)

	resp, err := http.Post(
		cfg.ServerURL+"/api/agent/result",
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		log.Println(err)
		return
	}

	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}