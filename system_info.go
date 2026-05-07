package main

import (
	"os"
	"runtime"
)

type SystemInfo struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	AgentID  string `json:"agent_id"`
}

func GetSystemInfo(cfg Config) SystemInfo {
	host, _ := os.Hostname()

	return SystemInfo{
		Hostname: host,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Version:  cfg.Version,
		AgentID:  cfg.AgentID,
	}
}