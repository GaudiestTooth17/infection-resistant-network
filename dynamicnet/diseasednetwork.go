package dynamicnet

import (
	"math/rand"
	"time"
)

// DiseasedNetwork represents a dynamic network with a disease that tries to adapt to slow
// the spread of the disease
type DiseasedNetwork struct {
	nodeState   []uint8
	numInfected uint32
	timeInState []int16
	rand        *rand.Rand
	disease     Disease
	adjMat      Network
	behavior    AgentBehavior
}

// NumNodes returns the number of nodes in a network
func (n *DiseasedNetwork) NumNodes() int {
	return n.adjMat.numNodes()
}

// NewDiseasedNetwork creates a new instance of DiseasedNetwork
func NewDiseasedNetwork(dis Disease, adjMat Network,
	infectionStrat InitialInfectionStrategy, behavior AgentBehavior) DiseasedNetwork {

	net := DiseasedNetwork{
		nodeState: make([]uint8, adjMat.numNodes()), numInfected: 0,
		timeInState: make([]int16, adjMat.numNodes()),
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
		disease:     dis,
		adjMat:      adjMat,
		behavior:    behavior,
	}

	for node := 0; node < net.adjMat.numNodes(); node++ {
		net.nodeState[node] = StateS
	}

	infectionStrat.apply(&net)
	return net
}

// Step through one time step
func (n *DiseasedNetwork) Step() {
	n.updateConnections()
	n.spreadInfection()
	n.updateStates()
	for i := range n.timeInState {
		n.timeInState[i]++
	}
}

func (n *DiseasedNetwork) spreadInfection() {
	infectiousNodes := n.FindNodesInState(StateI)
	atRiskGroups := make([][]int, len(infectiousNodes))
	for i, node := range infectiousNodes {
		atRiskGroups[i] = n.findNeighbors(node, StateS)
	}

	for _, group := range atRiskGroups {
		for _, node := range group {
			if n.rand.Float32() < n.disease.InfectionProbability() {
				n.changeState(node, StateE)
			}
		}
	}
}

// findNeighbors finds all the neighbors of node with the indicated state. Use -1 to find all neighbors.
func (n *DiseasedNetwork) findNeighbors(node int, state int) []int {
	neighbors := make([]int, 0)
	for neighbor := 0; neighbor < n.adjMat.numNodes(); neighbor++ {
		if n.adjMat.edgeWeight(node, neighbor) > 0 {
			if state == -1 || n.nodeState[neighbor] == uint8(state) {
				neighbors = append(neighbors, neighbor)
			}
		}
	}
	return neighbors
}

// change the state of node to StateS, StateE, StateI, or StateR (from disease package )
func (n *DiseasedNetwork) changeState(node, state int) {
	n.nodeState[node] = uint8(state)
	n.timeInState[node] = 0
}

// updateStates changes the state of exposed and infected nodes if they have been I/E for long enough
func (n *DiseasedNetwork) updateStates() {
	exposedNodes := n.FindNodesInState(StateE)
	infectedNodes := n.FindNodesInState(StateI)
	for _, node := range exposedNodes {
		if n.timeInState[node] == n.disease.TimeToI() {
			n.changeState(node, StateI)
		}
	}
	for _, node := range infectedNodes {
		if n.timeInState[node] == n.disease.TimeToR() {
			n.changeState(node, StateR)
		}
	}
}

// FindNodesInState finds all the nodes in the network with the given state
func (n *DiseasedNetwork) FindNodesInState(state int) []int {
	s := uint8(state)
	nodes := make([]int, 0)
	for node, st := range n.nodeState {
		if s == st {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// if this is slow, consider working with sets (maps instead)
func (n *DiseasedNetwork) updateConnections() {
	toAdd := make([][]int, n.NumNodes())
	toRemove := make([][]int, n.NumNodes())

	for node := 0; node < n.NumNodes(); node++ {
		toAdd[node] = n.findNeighborsToAdd(node)
		toRemove[node] = n.findNeighborsToRemove(node)
	}

	// add all the requested edges
	for node := 0; node < n.NumNodes(); node++ {
		for _, neighbor := range toAdd[node] {
			n.adjMat.addEdge(node, neighbor, 1)
		}
	}

	// removing takes precedence over adding, so remove all the requested edges next
	for node := 0; node < n.NumNodes(); node++ {
		for _, neighbor := range toRemove[node] {
			n.adjMat.removeEdge(node, neighbor)
		}
	}
}

func (n *DiseasedNetwork) findNeighborsToRemove(node int) []int {
	toRemove := make([]int, 0)
	infectedNeighbors := n.findNeighbors(node, StateI)
	for _, neighbor := range infectedNeighbors {
		if n.rand.Float32() < n.behavior.removeInfectedNeighborProb() {
			toRemove = append(toRemove, neighbor)
		}
	}
	return toRemove
}

func (n *DiseasedNetwork) findNeighborsToAdd(node int) []int {
	currentNeighbors := n.findNeighbors(node, -1)
	toAdd := make([]int, 0)
	for _, neighbor := range currentNeighbors {
		nOfn := n.findNeighbors(neighbor, -1)
		for _, nn := range nOfn {
			if n.rand.Float32() < n.behavior.addNeighborOfNeighborProb() {
				toAdd = append(toAdd, nn)
			}
		}
	}
	return toAdd
}