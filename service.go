package main

const SystemdUnit = `[Unit]
Description=CommandPilot Agent
After=network.target

[Service]
Type=simple
ExecStart=/opt/commandpilot-agent/pilot-agent
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
`