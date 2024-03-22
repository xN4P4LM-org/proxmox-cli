package main

import (
	"context"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/luthermonson/go-proxmox"
)

type ProxmoxNode struct {
	proxmox.Node
	VirtualMachines []VirtualMachine
}

// get nodes to ignore
func getNodesToIgnore() []string {
	nodesToIgnore, ok := os.LookupEnv("PROXMOX_IGNORE_NODE")

	if !ok {
		return []string{}
	}

	// split on comma
	return strings.Split(nodesToIgnore, ",")
}

func getNodes() {
	api_nodes, err := client.Nodes(context.Background())

	//empty existing nodes
	nodes = []ProxmoxNode{}

	if err != nil {
		logger.Fatal("Error fetching nodes: ", err.Error())
	}

	for _, nodeStatus := range api_nodes {
		node, err := client.Node(context.Background(), nodeStatus.Node)

		if err != nil {
			logger.Fatal("Error fetching node: ", err.Error())
		}

		vms := getVirtualMachines(*node)

		nodes = append(nodes, ProxmoxNode{
			//*nodeStatus,
			*node,
			vms,
		})

	}

	sort.Slice(activeVMids, func(i, j int) bool {
		return activeVMids[i] < activeVMids[j]
	})

}

type vmCount struct {
	vmCount int
	node    ProxmoxNode
}

func getLeastUsedNode() ProxmoxNode {
	allNodes := []vmCount{}
	canidateNodes := []vmCount{}

	leastNumberOfVMs := 0

Loop:
	for _, node := range nodes {
		for _, excluded_node := range getNodesToIgnore() {
			if excluded_node == node.Name {
				continue Loop
			}
		}

		allNodes = append(allNodes, vmCount{
			vmCount: len(node.VirtualMachines),
			node:    node,
		})

		if leastNumberOfVMs == 0 || len(node.VirtualMachines) < leastNumberOfVMs {
			leastNumberOfVMs = len(node.VirtualMachines)
		}
	}

	// filter out nodes with the least number of VMs
	for index, vm := range allNodes {
		// if the node has more VMs than the least number of VMs, remove it
		if vm.vmCount > leastNumberOfVMs {
			continue
		}
		canidateNodes = append(canidateNodes, allNodes[index])
	}

	// if there are more than one node select a random one
	if len(canidateNodes) > 1 {
		return canidateNodes[rand.Intn(len(canidateNodes))].node
	}

	// confirm there are possible nodes
	if len(canidateNodes) == 0 {
		logger.Fatal("No nodes available")
	}

	// return the only node
	return canidateNodes[0].node
}

func getNodeByVMID(VMID uint64) ProxmoxNode {
	for _, node := range nodes {
		for _, vm := range node.VirtualMachines {
			if uint64(vm.VMID) == VMID {
				return node
			}
		}
	}

	logger.Fatal("VM not found")
	return ProxmoxNode{}
}
