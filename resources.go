package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// Help function

func resourceHelp() {
	fmt.Println("Usage: proxmox resources [options]")
	fmt.Println("List resource distribution between nodes")
	fmt.Println("Options:")
	fmt.Println("  -h, --help")

	fmt.Println("Examples:")
	fmt.Println("  proxmox resources")
}

func resources() {
	// get options from parameters
	params := os.Args[2:]

	if checkParams(params, "help", false) {
		resourceHelp()
		os.Exit(0)
	}

	listAllResources()
}

// List all node resources
func listAllResources() {
	// create tabwriter
	tabWriter := new(tabwriter.Writer)
	tabWriter.Init(os.Stdout, 5, 8, 1, '\t', 0)

	// print header
	fmt.Fprintln(tabWriter, "Node\tCPU (%)\tMemory Used (%)\tVMs\t")
	fmt.Fprintln(tabWriter, "----\t---\t------\t-------\t")

	// loop nodes
	for _, node := range nodes {
		// get node resources
		cpu := node.CPU / 100
		memory := float64(node.Memory.Total) / float64(node.Memory.Used)
		vms := len(node.VirtualMachines)

		// print node resources
		fmt.Fprintf(tabWriter, "%s\t%.2f\t%.2f\t%d\t\n", node.Name, cpu, memory, vms)
	}

	// flush tabwriter
	tabWriter.Flush()

}
