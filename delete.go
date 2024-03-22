package main

import (
	"fmt"
	"os"
)

// Help function
func deleteHelp() {
	fmt.Println("Usage: proxmox delete [options]")
	fmt.Println("Delete a VM or container")
	fmt.Println("Options:")
	fmt.Println("  -n, --name <name>")

	fmt.Println("Examples:")
	fmt.Println("  proxmox delete -n test-vm")
}

func delete() {
	// get options from parameters
	params := os.Args[2:]

	if checkParams(params, "help", false) {
		deleteHelp()
		os.Exit(0)
	}
}
