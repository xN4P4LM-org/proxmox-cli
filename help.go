package main

import "fmt"

// Print help functions

func help() {
	// print help
	fmt.Println("Usage: proxmox [options] [command]")
	fmt.Println("Global Options:")
	fmt.Println("  -h, --help")
	fmt.Println("  -d, --debug - Enable debug mode")
	fmt.Println("Commands:")
	fmt.Println("  list")
	fmt.Println("  create")
	fmt.Println("  delete")
}
