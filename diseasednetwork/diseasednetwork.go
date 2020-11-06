package diseasednetwork

import (
	"math/rand"
	"time"
)

// Void is an empty struct used for more space efficient maps/sets
type Void struct{}

// DiseasedNetwork represents a dynamic network with a disease that tries to adapt to slow
// the spread of the disease
type DiseasedNetwork struct {
	diseases []Disease
	adjMat   Network
}

// NumNodes returns the number of nodes in a network
func (n *DiseasedNetwork) NumNodes() int {
	return n.adjMat.NumNodes()
}

// NewDiseasedNetwork creates a new instance of DiseasedNetwork
func NewDiseasedNetwork(diseases []Disease, underlyingNet Network) DiseasedNetwork {

	net := DiseasedNetwork{
		diseases: diseases,
		adjMat:   underlyingNet.MakeCopy(),
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
func (n *DiseasedNetwork) Step() time.Duration {
	stepStart := time.Now()
	n.spreadInfection()
	n.updateStates()
	for _, dis := range n.diseases {
		for node := 0; node < n.NumNodes(); node++ {
			dis.IncTimeInState(node)
		}
	}
	return time.Now().Sub(stepStart)
}

// spreadInfection of all the diseases by finding the infectious nodes, the nodes they could infect,
// and then randomly determining if each one will get infected
func (n *DiseasedNetwork) spreadInfection() {
	for i, disease := range n.diseases {
		infectiousNodes := disease.FindNodesInState(StateI)
		atRiskGroups := make([]map[int]uint8, len(infectiousNodes))
		groupIndex := 0
		for node := range infectiousNodes {
			atRiskGroups[groupIndex] = n.findNeighbors(node, StateS, i)
			groupIndex++
		}

		for _, group := range atRiskGroups {
			for node := range group {
				if rand.Float32() < disease.InfectionProbability() {
					disease.SetState(node, StateE)
				}
			}
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
		exposedNodes := disease.FindNodesInState(StateE)
		infectedNodes := disease.FindNodesInState(StateI)
		for node := range exposedNodes {
			if disease.TimeInState(node) == disease.TimeToI() {
				disease.SetState(node, StateI)
			}
		}
		for node := range infectedNodes {
			if disease.TimeInState(node) == disease.TimeToR() {
				disease.SetState(node, StateR)
			}
		}
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
