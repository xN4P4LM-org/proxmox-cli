package main

import (
	"context"
	"strings"

	"github.com/luthermonson/go-proxmox"
)

var activeVMids = []uint64{}

// VirtualMachine struct
type VirtualMachine struct {
	proxmox.VirtualMachine
	vm_type string
	tags    []string
}

func getVirtualMachines(node proxmox.Node) []VirtualMachine {
	vms, err := node.VirtualMachines(context.Background())

	if err != nil {
		logger.Fatal("Error fetching VMs for node: ", node.Name)
	}

	vmArray := []VirtualMachine{}

	for _, vm := range vms {
		tags := strings.Split(vm.Tags, ";")

		vmType := "VM"
		if vm.Template {
			vmType = "Template"
		}

		vmArray = append(vmArray, VirtualMachine{
			*vm,
			vmType,
			tags,
		})

		activeVMids = append(activeVMids, uint64(vm.VMID))
	}

	return vmArray
}

func getTemplates() []VirtualMachine {
	// update nodes
	updateNodes()

	templates := []VirtualMachine{}

	for _, node := range nodes {
		for _, vm := range node.VirtualMachines {
			if vm.vm_type == "Template" {
				templates = append(templates, vm)
			}
		}
	}

	return templates
}

func getTemplate(source uint64) VirtualMachine {
	// update nodes
	updateNodes()

	for _, node := range nodes {
		for _, vm := range node.VirtualMachines {
			if uint64(vm.VMID) == source {
				return vm
			}
		}
	}

	logger.Fatal("Template not found")
	return VirtualMachine{}
}

func getVMbyID(source uint64) VirtualMachine {
	// update nodes
	updateNodes()

	for _, node := range nodes {
		for _, vm := range node.VirtualMachines {
			if uint64(vm.VMID) == source {
				return vm
			}
		}
	}
	logger.Fatal("VM with id ", source, " not found")
	return VirtualMachine{}
}
