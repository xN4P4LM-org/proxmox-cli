package main

import (
	"encoding/json"
	"os"
)

func getEnvironment(environment string) string {

	type Environment struct {
		Environment string `json:"enviroment"`
		VMID        string `json:"vmid"`
	}

	var environments map[string]Environment

	// get the environment from environment.json
	environmentsFile, readErr := os.ReadFile(findFile())

	if readErr != nil {
		logger.Fatal("Error opening environments.json")
	}

	unmarshalErr := json.Unmarshal(environmentsFile, &environments)

	if unmarshalErr != nil {
		logger.Fatal("Error unmarshalling environments.json")
	}

	// get the environment
	for _, env := range environments {
		if env.Environment == environment {
			return env.VMID
		}
	}

	logger.Fatal("Environment not found")
	return ""
}

func findFile() string {
	possible_locations := []string{
		"$HOME/.config/proxmox-cli/environments.json",
		"/etc/proxmox-cli/environments.json",
	}

	for _, location := range possible_locations {
		if checkForFile(location) {
			return location
		}
	}

	logger.Fatal("environments.json not found")
	return ""
}

func checkForFile(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
