package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/luthermonson/go-proxmox"
)

// help function
func cloneHelp() {
	fmt.Println("Usage: proxmox clone [options]")
	fmt.Println("Clone a template VM and create on the least used node with the next available VM ID")
	fmt.Println("Options:")
	fmt.Println("  -n, --name <name> - (required) - Name of the VM")
	fmt.Println("  -s, --source <sourceVM> - (required) - Source VM to clone")
	fmt.Println("  -e, --environment <environment> - (required) - Environment to use")
	fmt.Println("  -p, --poweron - (optional) - Power on the VM after creation")

	fmt.Println("Examples:")
	fmt.Println("  proxmox create -n test-vm -s test-vm-template")
}

func cloneMapper() {
	// get options from parameters
	params := os.Args[2:]

	if checkParams(params, "help", false) {
		cloneHelp()
		os.Exit(0)
	}

	if !checkParams(params, "name", true) {
		logger.Print("Name is required")
		cloneHelp()
		os.Exit(0)
	}

	if !checkParams(params, "source", true) {
		logger.Print("Source is required")
		cloneHelp()
		os.Exit(0)
	}

	if !checkParams(params, "environment", true) {
		logger.Print("Environment is required")
		cloneHelp()
		os.Exit(0)
	}

	cloneVM(params)

}

func cloneVM(params []string) {

	if len(params) < 3 {
		cloneHelp()
		os.Exit(1)
	}

	// get source and destination
	source := getParams(params, "source")
	vmName := getParams(params, "name")
	environment := getParams(params, "environment")
	powerOn := false

	if checkParams(params, "poweron", false) {
		powerOn = true
	}

	if checkParams(params, "target", false) {
		logger.Fatal("Target not implemented")

	}

	sourceint, err := strconv.ParseUint(source, 10, 64)

	if err != nil {
		logger.Fatal("Source is not a VM ID")
	}

	// get source template
	templateVM := getTemplate(sourceint)

	// get the least used node
	targetNode := getLeastUsedNode()

	// clone the VM
	vmID := cloneLocally(templateVM, environment, vmName)

	// update nodes
	updateNodes()

	// get current node for VM
	vmNode := getNodeByVMID(vmID)

	// get the new VM
	newVM := getVMbyID(vmID)

	if vmNode.Node.Name != targetNode.Node.Name {
		// migrate to the least used node
		migrateVM(newVM, targetNode)
	}

	migratedVM := getVMbyID(vmID)

	if powerOn {
		powerOnVM(migratedVM)
	}

	logger.Print(
		"All operations completed successfully for creation of VM: ", vmName,
		" with ID: ", vmID, " in environment: ", environment, " on node: ", targetNode.Node.Name,
		" from template: ", source, " and is currently powered: ", newVM.Status,
	)
}

func cloneLocally(templateVM VirtualMachine, environment string, name string) uint64 {

	intRange, parseErr := strconv.ParseInt(getEnvironment(environment), 10, 64)

	if parseErr != nil {
		logger.Fatal("Error parsing environment range")
	}

	//get the next available VM ID

	for _, vm := range activeVMids {
		if vm == uint64(intRange) {
			intRange++
		}
	}

	cloneOptions := proxmox.VirtualMachineCloneOptions{
		Full:  1,
		NewID: int(intRange),
		Pool:  environment,
		Name:  name,
	}

	_, proxmoxTask, err := templateVM.Clone(context.Background(), &cloneOptions)

	if err != nil {
		logger.Fatal("Error cloning VM: ", err.Error())
	}

	for {
		status, complete, err := proxmoxTask.WaitForCompleteStatus(context.Background(), 60, 10)

		if err != nil {
			logger.Fatal("Error waiting for clone operation to complete: ", err.Error())
		}

		if complete {
			if status {
				break
			} else {
				logger.Fatal("Clone failed")
			}
		}
	}

	logger.Print("Successfully cloned VM: ", name, " with ID: ", intRange, " from template: ", templateVM.Name)
	return uint64(intRange)
}

func migrateVM(vm VirtualMachine, targetNode ProxmoxNode) {

	migrateOptions := proxmox.VirtualMachineMigrateOptions{
		Target:        targetNode.Node.Name,
		TargetStorage: getStorage(&targetNode.Node).Name,
	}

	proxmoxTask, err := vm.Migrate(context.Background(), &migrateOptions)

	if err != nil {
		logger.Fatal("Error migrating VM: ", err.Error())
	}

	// wait for the task to complete
	for {
		status, complete, err := proxmoxTask.WaitForCompleteStatus(context.Background(), 60, 10)

		if err != nil {
			logger.Fatal("Error waiting for migrate operation to complete: ", err.Error())
		}

		if complete {
			if status {
				break
			} else {
				logger.Fatal("Migrate failed")
			}
		}
	}

	logger.Print("Migration of ", vm.Name, " to ", targetNode.Node.Name, " complete")

}

func powerOnVM(vm VirtualMachine) {

	_, err := vm.Start(context.Background())

	if err != nil {
		logger.Fatal("Error starting VM: ", err.Error())
	}

	for {
		targetVM := getVMbyID(uint64(vm.VMID))
		if targetVM.IsRunning() {
			break
		} else {
			logger.Fatal("Start failed")
		}
	}

	logger.Print(vm.Name, " has been successfully powered on")
}
