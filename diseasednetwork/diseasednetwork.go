package diseasednetwork

import (
	"math/rand"
	"time"
)

// Void is an empty struct used for more space efficient maps/sets
type Void struct{}

// DiseasedNetwork represents a dynamic network with a disease that tries to adapt to slow
// the spread of the disease. The bad disease should be put in slot 0.
type DiseasedNetwork struct {
	diseases   []Disease
	adjMat     Network
	stepNum    uint
	PlotMakers []PlotMaker
}

// NumNodes returns the number of nodes in a network
func (n *DiseasedNetwork) NumNodes() int {
	return n.adjMat.NumNodes()
}

// NewDiseasedNetwork creates a new instance of DiseasedNetwork
func NewDiseasedNetwork(underlyingNet *Network, diseases []Disease, plotMakers []PlotMaker) DiseasedNetwork {
	net := DiseasedNetwork{
		diseases:   diseases,
		adjMat:     underlyingNet.MakeCopy(),
		stepNum:    0,
		PlotMakers: plotMakers,
	}

	for _, disease := range net.diseases {
		disease.SetNumNodes(net.adjMat.NumNodes())
		for node := 0; node < disease.NumNodes(); node++ {
			disease.SetState(node, StateS)
		}
		infectionStrategy := disease.InitialInfection()
		infectionStrategy.apply(disease)
	}

	return net
}

// Step through one time step
func (n *DiseasedNetwork) Step() (time.Duration, float64) {
	stepStart := time.Now()
	n.spreadInfection()
	n.updateStates()
	for _, dis := range n.diseases {
		for node := 0; node < n.NumNodes(); node++ {
			dis.IncTimeInState(node)
		}
	}
	// let the plotMakers measure the network
	for _, plotMaker := range n.PlotMakers {
		plotMaker.feedInformation(n)
	}
	// let the diseases reset their counting
	r0 := n.diseases[0].R0()
	for _, dis := range n.diseases {
		dis.endStep()
	}
	n.stepNum++
	return time.Now().Sub(stepStart), r0
}

// spreadInfection of all the diseases by finding the infectious nodes, the nodes they could infect,
// and then randomly determining if each one will get infected
func (n *DiseasedNetwork) spreadInfection() {
	for i, disease := range n.diseases {
		infectiousNodes := disease.FindNodesInState(StateI)
		atRiskGroups := make(map[int]map[int]uint8, len(infectiousNodes))
		for node := range infectiousNodes {
			atRiskGroups[node] = n.findNeighbors(node, StateS, i)
		}

		for infectiousNode, group := range atRiskGroups {
			nodesInfected := uint(0)
			for atRiskNode := range group {
				if rand.Float32() < disease.InfectionProbability() {
					disease.SetState(atRiskNode, StateE)
					nodesInfected++
				}
			}
			disease.ReportInfections(infectiousNode, nodesInfected)
		}
	}
}

// findNeighbors finds all the neighbors of node with the indicated state.
// Use a negative value to find all neighbors.
func (n *DiseasedNetwork) findNeighbors(node int, state int, diseaseIndex int) map[int]uint8 {
	neighbors := n.adjMat.NeighborsOf(node)
	if state < 0 {
		return neighbors
	}

	neighborsInState := make(map[int]uint8)
	for neighbor, edgeWeight := range neighbors {
		if n.diseases[diseaseIndex].State(neighbor) == uint8(state) {
			neighborsInState[neighbor] = edgeWeight
		}
	}
	return neighborsInState
}

// updateStates changes the state of exposed and infected nodes if they have been I/E for long enough.
// This happens for each disease in diseases.
func (n *DiseasedNetwork) updateStates() {
	for _, disease := range n.diseases {
		disease.updateStates()
	}
}

// FindNodesInState returns a set of all the nodes in a certain state in the specified disease
func (n *DiseasedNetwork) FindNodesInState(state int, diseaseIndex int) map[int]Void {
	s := uint8(state)
	nodes := make(map[int]Void)
	for node := 0; node < n.NumNodes(); node++ {
		if n.diseases[diseaseIndex].State(node) == s {
			nodes[node] = Void{}
		}
	}
	return nodes
}

// GetNodeStates returns a slice where the ith index contains the state of the ith node
// in the disease number provided. This is useful for visualization.
func (n *DiseasedNetwork) GetNodeStates(diseaseNumber int) []uint8 {
	nodeStates := make([]uint8, n.NumNodes())
	for node := 0; node < n.NumNodes(); node++ {
		nodeStates[node] = n.diseases[diseaseNumber].State(node)
	}
	return nodeStates
}

// R0 gives the R0 of the specified disease.
func (n *DiseasedNetwork) R0(diseaseNum int) float64 {
	return n.diseases[diseaseNum].R0()
}
