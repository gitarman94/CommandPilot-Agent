package main

import (
	"net"
	"os"
	"runtime"
	"strings"
)

type SystemInfo struct {
	AgentID         string  `json:"agent_id"`
	Version         string  `json:"version"`
	Hostname        string  `json:"hostname"`
	FQDN            string  `json:"fqdn"`
	IP              string  `json:"ip"`
	OS              string  `json:"os"`
	Architecture    string  `json:"architecture"`
	DeviceType      string  `json:"device_type"`
	DeviceModel     string  `json:"device_model"`
	CPUModel        string  `json:"cpu_model"`
	CPUUsage        float64 `json:"cpu_usage"`
	RAMTotal        int64   `json:"ram_total"`
	RAMUsed         int64   `json:"ram_used"`
	RAMUsagePercent float64 `json:"ram_usage_percent"`
	DiskTotal       int64   `json:"disk_total"`
	DiskUsed        int64   `json:"disk_used"`
	DiskFree        int64   `json:"disk_free"`
	DiskFreeHuman   string  `json:"disk_free_human"`
}

func GetSystemInfo(cfg Config) SystemInfo {
	host, _ := os.Hostname()

	return SystemInfo{
		AgentID:      cfg.AgentID,
		Version:      cfg.Version,
		Hostname:     host,
		FQDN:         host,
		IP:           primaryIPv4(),
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		DeviceType:   "endpoint",
		DeviceModel:  "",
		CPUModel:     "",
		CPUUsage:     0,
		RAMTotal:     0,
		RAMUsed:      0,
		RAMUsagePercent: 0,
		DiskTotal:    0,
		DiskUsed:     0,
		DiskFree:     0,
		DiskFreeHuman: "",
	}
}

func primaryIPv4() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		if (iface.Flags & net.FlagLoopback) != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			s := ip.String()
			if strings.TrimSpace(s) != "" {
				return s
			}
		}
	}

	return ""
}