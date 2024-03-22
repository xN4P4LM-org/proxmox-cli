package main

import (
	"log"
	"os"

	"github.com/luthermonson/go-proxmox"
)

var nodes = []ProxmoxNode{}

var client = connect()

var proxmoxLogger = &proxmox.LeveledLogger{Level: proxmox.LevelWarn}

var logger = log.New(os.Stdout, "", log.Default().Flags())

// Main function

func main() {

	args := os.Args

	// set verbose logging
	if checkParams(args, "debug", false) {
		logger.SetFlags(log.LstdFlags | log.Lshortfile)
		proxmoxLogger.Level = proxmox.LevelDebug
	}

	if os.Args == nil || len(args) < 2 {
		help()
		os.Exit(1)
	}

	// switch operation
	switch os.Args[1] {
	case "resources":
		updateNodes()
		resources()
		os.Exit(0)
	case "list":
		updateNodes()
		list()
		os.Exit(0)
	case "clone":
		updateNodes()
		cloneMapper()
		os.Exit(0)
	case "delete":
		updateNodes()
		delete()
		os.Exit(0)
	case "nextNode":
		updateNodes()
		print("The next least used node is ", getLeastUsedNode().Name)
		os.Exit(0)
	default:
		help()
		os.Exit(1)
	}

}

// update node loop
func updateNodes() {
	// get nodes from config
	getNodes()
}
