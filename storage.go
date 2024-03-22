package main

import (
	"context"
	"os"
	"strings"

	"github.com/luthermonson/go-proxmox"
)

func getStorage(targetNode *proxmox.Node) *proxmox.Storage {
	// Get node storages
	storages, storageErr := targetNode.Storages(context.Background())

	if storageErr != nil {
		logger.Fatal("Error fetching storages for node: ", targetNode.Name)
	}

	// ignore local-lvm storage, as it's the default storage

	// get storage name that's not local-lvm
	targetStorage := &proxmox.Storage{}

	for _, storage := range storages {

		if storage.Enabled == 0 {
			continue
		}
		if storage.Shared == 1 {
			continue
		}
		if checkIgnoredStorage(storage.Name) {
			continue
		}

		targetStorage = storage
		break
	}

	return targetStorage
}

func checkIgnoredStorage(storageName string) bool {
	ignoredStorages := getIgnoredStorages()

	for _, storage := range ignoredStorages {
		if storage == storageName {
			return true
		}
	}

	return false
}

func getIgnoredStorages() []string {
	ignoredStorages, ok := os.LookupEnv("PROXMOX_IGNORE_STORAGE")

	if !ok {
		return []string{}
	}

	// split on comma
	return strings.Split(ignoredStorages, ",")
}
