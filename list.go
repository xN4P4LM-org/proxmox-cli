package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// Help function
func listHelp() {
	fmt.Println("Usage: proxmox list [options]")
	fmt.Println("List Virtual Machines")
	fmt.Println("Options:")
	fmt.Println("  -h, --help")
	fmt.Println("  -n, --node <node>")
	fmt.Println("  -t, --template")

	fmt.Println("Examples:")
	fmt.Println("  proxmox list")
}

// List operation
func list() {
	// get options from parameters
	params := os.Args[2:]

	if checkParams(params, "help", false) {
		listHelp()
		os.Exit(0)
	}

	// if leng of parameters is less than 1, print all VMs
	if len(params) < 1 {
		listAllNodes()
		return
	}

	// if length of parameters is 1+ get the --node parameter
	if checkParams(params, "node", true) {
		listNode(params[1])
		return
	}

	// if length of parameters is 1+ get the --template parameter
	if checkParams(params, "template", false) {
		listTemplates()
		return
	}

	logger.Fatal("Invalid parameters")
}

// List single node virtual machines
func listNode(node string) {

}

// List all virtual machines
func listAllNodes() {

	// create tabwriter
	tabWriter := new(tabwriter.Writer)
	tabWriter.Init(os.Stdout, 5, 8, 1, '\t', 0)

	// print header
	fmt.Fprintln(tabWriter, "Node\tVMID\tName\tStatus\tTags\t")
	fmt.Fprintln(tabWriter, "----\t----\t----\t------\t----\t")

	// loop nodes
	for _, node := range nodes {

		// loop virtual machines
		for _, vm := range node.VirtualMachines {
			// print virtual machine
			fmt.Fprintf(tabWriter, "%s\t%d\t%s\t%s\t%s\t\n", node.Name, vm.VMID, vm.Name, vm.Status, vm.Tags)
		}

	}

	// flush tabwriter
	tabWriter.Flush()

}

// List all templates
func listTemplates() {
	templates := getTemplates()

	// create tabwriter
	tabWriter := new(tabwriter.Writer)
	tabWriter.Init(os.Stdout, 5, 8, 1, '\t', 0)

	// print header
	fmt.Fprintln(tabWriter, "Node\tVMID\tName\t")

	// loop templates
	for _, template := range templates {
		// print template
		fmt.Fprintf(tabWriter, "%s\t%d\t%s\t\n", template.Node, template.VMID, template.Name)
	}

	// flush tabwriter
	tabWriter.Flush()

}
