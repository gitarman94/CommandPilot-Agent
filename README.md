# ⚠️ EARLY DEVELOPMENT — EXPECT BREAKAGE

# CommandPilot Agent

CommandPilot Agent is the endpoint component of the CommandPilot platform.

It is a lightweight cross-platform system agent written entirely in Go and designed to integrate directly with the CommandPilot control plane.

The agent continuously reports device telemetry, system inventory, heartbeat status, and executes actions dispatched from the CommandPilot server.

---

# 🧠 Overview

The CommandPilot ecosystem consists of:

```text
CommandPilot
├── pilot-core      (central server / control plane)
├── pilot-agent     (endpoint agent)
├── web UI          (administrative interface)
└── REST API        (integration layer)
```

This repository contains the Go-based endpoint agent.

Repository:

* [CommandPilot-Agent Repository](https://github.com/gitarman94/CommandPIlot-Agent?utm_source=chatgpt.com)

---

# ✨ Features

* Pure Go implementation
* Cross-platform architecture
* Windows and Linux support
* Persistent heartbeat reporting
* Device inventory synchronization
* Remote action polling
* Command execution support
* Server-driven update architecture
* Lightweight runtime footprint
* systemd service support
* Windows service support
* Automatic reconnect behavior
* SHA256 binary validation
* UUID-based endpoint identity
* Configurable server targeting

---

# 📁 Repository Structure

```text
CommandPilot-Agent/
├── main.go
├── action.go
├── command.go
├── device.go
├── heartbeat.go
├── inventory.go
├── service.go
├── system_info.go
├── update.go
├── go.mod
├── README.md
├── setup_or_update_agent.sh
└── setup_or_update_agent.ps1
```

---

# 🔧 Core Responsibilities

The agent is responsible for:

* registering devices with the server
* maintaining heartbeat connectivity
* reporting system metadata
* reporting operating system information
* polling for queued actions
* executing server-issued commands
* returning action output/results
* receiving update metadata from the server
* performing controlled binary upgrades

---

# 🌐 Server Communication

The agent communicates directly with the CommandPilot server over HTTP/HTTPS.

Typical server endpoint:

```text
http://<server-ip>:8080
```

The installer prompts for the server automatically.

---

# 🖥️ Supported Platforms

## Linux

Supported distributions:

* Debian
* Ubuntu
* Linux Mint
* Pop!_OS
* Raspberry Pi OS

## Windows

Supported versions:

* Windows 10
* Windows 11
* Windows Server 2019+
* Windows Server 2022+

---

# ⚡ Linux Installation

Installer:

* [setup_or_update_agent.sh](https://github.com/gitarman94/CommandPIlot-Agent/blob/main/setup_or_update_agent.sh?utm_source=chatgpt.com)

## Standard Install

```bash
curl -fsSL https://raw.githubusercontent.com/gitarman94/CommandPIlot-Agent/main/setup_or_update_agent.sh | bash
```

## Silent Install

```bash
curl -fsSL https://raw.githubusercontent.com/gitarman94/CommandPIlot-Agent/main/setup_or_update_agent.sh | bash -s -- --silent --server 192.168.1.10
```

## Upgrade Existing Installation

```bash
curl -fsSL https://raw.githubusercontent.com/gitarman94/CommandPIlot-Agent/main/setup_or_update_agent.sh | bash -s -- --upgrade
```

---

# ⚡ Windows Installation

Installer:

* [setup_or_update_agent.ps1](https://github.com/gitarman94/CommandPIlot-Agent/blob/main/setup_or_update_agent.ps1?utm_source=chatgpt.com)

## Interactive Install

```powershell
powershell.exe -ExecutionPolicy Bypass -File .\setup_or_update_agent.ps1
```

## Silent Install

```powershell
powershell.exe -ExecutionPolicy Bypass -File .\setup_or_update_agent.ps1 -Silent -ServerUrl 192.168.1.10
```

## Uninstall

```powershell
powershell.exe -ExecutionPolicy Bypass -File .\setup_or_update_agent.ps1 -Uninstall
```

---

# 🔄 Update Architecture

The CommandPilot agent uses a server-driven update model.

The server will:

1. store approved agent builds
2. expose update metadata
3. notify agents of newer versions
4. provide update download packages

The agent will:

1. compare installed version vs server version
2. download approved update packages
3. validate integrity
4. replace the running binary
5. restart itself safely

---

# 🧾 Configuration Files

## Linux

```text
/opt/commandpilot-agent/
```

## Windows

```text
C:\CommandPilot_Agent\
```

Common files:

```text
server_url.txt
agent_id.txt
version.txt
```

---

# 🔐 Security Model

Current security controls include:

* UUID-based agent identity
* SHA256 binary verification
* bcrypt password authentication on server
* session-authenticated administration
* controlled update channels
* authenticated action APIs

Planned improvements:

* mutual TLS
* signed update packages
* API tokens
* certificate pinning
* encrypted agent registration

---

# 📡 Agent Telemetry

The agent continuously reports:

* hostname
* operating system
* architecture
* IP addresses
* uptime
* last check-in
* agent version
* inventory metadata

---

# 🛠️ Remote Actions

The agent supports server-issued actions including:

* shell command execution
* PowerShell execution
* inventory refresh
* telemetry refresh
* service restart
* update operations

---

# 🔁 Services

## Linux

Installed as:

```text
commandpilot-agent.service
```

Managed via:

```bash
systemctl status commandpilot-agent
systemctl restart commandpilot-agent
journalctl -u commandpilot-agent
```

## Windows

Installed as:

```text
CommandPilotAgent
```

Managed through:

* Services.msc
* PowerShell
* sc.exe

---

# ⚠️ Current Limitations

* Early-stage platform
* No TLS by default
* No signed update packages yet
* Limited RBAC enforcement
* Minimal endpoint hardening
* No offline action queueing
* No delta updates
* No multi-tenant separation

---

# 🔮 Planned Features

* encrypted transport
* WebSocket communications
* streaming telemetry
* remote shell sessions
* file transfer
* policy enforcement
* patch orchestration
* endpoint isolation
* certificate authentication
* fleet grouping
* live event streaming

---

# 📜 License

TBD

---

# 🙋 Related Repositories

## Server

* [CommandPilot Server Repository](https://github.com/gitarman94/CommandPilot?utm_source=chatgpt.com)

## Agent

* [CommandPilot-Agent Repository](https://github.com/gitarman94/CommandPIlot-Agent?utm_source=chatgpt.com)
