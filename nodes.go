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
	candidateNodes := []vmCount{}

	vmcount := []int{}

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

		vmcount = append(vmcount, len(node.VirtualMachines))
	}

	// get the minimum numbe from the list
	sort.Ints(vmcount)

	for _, node := range allNodes {
		if node.vmCount == vmcount[0] {
			candidateNodes = append(candidateNodes, node)
		}
	}

	// if there are more than one node select a random one
	if len(candidateNodes) > 1 {
		return candidateNodes[rand.Intn(len(candidateNodes))].node
	}

	// confirm there are possible nodes
	if len(candidateNodes) == 0 {
		logger.Fatal("No nodes available")
	}

	// return the only node
	return candidateNodes[0].node
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
