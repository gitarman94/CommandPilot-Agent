package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Config struct {
	BaseDir   string
	ServerURL string
	AgentID   string
	Version   string
}

func LoadConfig() Config {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	baseDir := filepath.Dir(exe)

	cfg := Config{
		BaseDir: baseDir,
	}

	server := strings.TrimSpace(os.Getenv("COMMANDPILOT_SERVER"))
	if server == "" {
		server = readTextFile(filepath.Join(baseDir, "server_url.txt"))
	}
	cfg.ServerURL = normalizeServerURL(server)

	if cfg.ServerURL == "" {
		cfg.ServerURL = "http://127.0.0.1:8080"
	}

	cfg.AgentID = strings.TrimSpace(os.Getenv("COMMANDPILOT_AGENT_ID"))
	if cfg.AgentID == "" {
		cfg.AgentID = strings.TrimSpace(readTextFile(filepath.Join(baseDir, "agent_id.txt")))
	}
	if cfg.AgentID == "" {
		cfg.AgentID = uuid.NewString()
		_ = os.WriteFile(filepath.Join(baseDir, "agent_id.txt"), []byte(cfg.AgentID+"\n"), 0600)
	}

	cfg.Version = strings.TrimSpace(os.Getenv("COMMANDPILOT_VERSION"))
	if cfg.Version == "" {
		cfg.Version = strings.TrimSpace(readTextFile(filepath.Join(baseDir, "version.txt")))
	}
	if cfg.Version == "" {
		cfg.Version = "dev"
	}

	return cfg
}

func readTextFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func normalizeServerURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return strings.TrimRight(raw, "/")
	}

	host := u.Host
	if host == "" {
		host = u.Path
		u.Path = ""
	}

	if host == "" {
		return ""
	}

	if !strings.Contains(host, ":") {
		host = net.JoinHostPort(host, "8080")
	}

	u.Host = host
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawQuery = ""
	u.Fragment = ""

	out := strings.TrimRight(u.String(), "/")
	if out == "" {
		return fmt.Sprintf("http://%s", host)
	}

	return out
}