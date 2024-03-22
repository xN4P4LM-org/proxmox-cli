package main

import (
	"net/http"

	"github.com/luthermonson/go-proxmox"
)

// Get Proxmox Client

func connect() *proxmox.Client {
	httpClient := http.Client{}

	client := proxmox.NewClient(
		getEnvVariable("PROXMOX_URL"),
		proxmox.WithHTTPClient(&httpClient),
		proxmox.WithAPIToken(getEnvVariable("PROXMOX_TOKEN_ID"), getEnvVariable("PROXMOX_TOKEN_SECRET")),
		proxmox.WithLogger(proxmoxLogger),
	)

	return client
}
