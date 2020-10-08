package dynamicnet

import (
	"math/rand"
	"time"
)

type void struct{}

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
	return n.adjMat.NumNodes()
}

// NewDiseasedNetwork creates a new instance of DiseasedNetwork
func NewDiseasedNetwork(dis Disease, adjMat Network,
	infectionStrat InitialInfectionStrategy, behavior AgentBehavior) DiseasedNetwork {

	net := DiseasedNetwork{
		nodeState: make([]uint8, adjMat.NumNodes()), numInfected: 0,
		timeInState: make([]int16, adjMat.NumNodes()),
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
		disease:     dis,
		adjMat:      adjMat,
		behavior:    behavior,
	}

	for node := 0; node < net.adjMat.NumNodes(); node++ {
		net.nodeState[node] = StateS
	}

	infectionStrat.apply(&net)
	return net
}

// Step through one time step
func (n *DiseasedNetwork) Step() time.Duration {
	stepStart := time.Now()
	n.updateConnections()
	n.spreadInfection()
	n.updateStates()
	for i := range n.timeInState {
		n.timeInState[i]++
	}
	return time.Now().Sub(stepStart)
}

func (n *DiseasedNetwork) spreadInfection() {
	infectiousNodes := n.FindNodesInState(StateI)
	atRiskGroups := make([]map[int]uint8, len(infectiousNodes))
	for i, node := range infectiousNodes {
		atRiskGroups[i] = n.findNeighbors(node, StateS)
	}

	for _, group := range atRiskGroups {
		for node := range group {
			if n.rand.Float32() < n.disease.InfectionProbability() {
				n.changeState(node, StateE)
			}
		}
	}
}

// findNeighbors finds all the neighbors of node with the indicated state. Use -1 to find all neighbors.
func (n *DiseasedNetwork) findNeighbors(node int, state int) map[int]uint8 {
	neighbors := n.adjMat.NeighborsOf(node)
	if state == -1 {
		return neighbors
	}

	neighborsInState := make(map[int]uint8)
	for neighbor, edgeWeight := range neighbors {
		if n.nodeState[neighbor] == uint8(state) {
			neighborsInState[neighbor] = edgeWeight
		}
	}
	return neighborsInState
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
	toAdd := make([]map[int]void, n.NumNodes())
	toRemove := make([]map[int]void, n.NumNodes())

	for node := 0; node < n.NumNodes(); node++ {
		toAdd[node] = n.findNeighborsToAdd(node)
		toRemove[node] = n.findNeighborsToRemove(node)
	}

	// add all the requested edges
	for node := 0; node < n.NumNodes(); node++ {
		for neighbor := range toAdd[node] {
			n.adjMat.AddEdge(node, neighbor, 1)
		}
	}

	// removing takes precedence over adding, so remove all the requested edges next
	for node := 0; node < n.NumNodes(); node++ {
		for neighbor := range toRemove[node] {
			n.adjMat.removeEdge(node, neighbor)
		}
	}
}

func (n *DiseasedNetwork) findNeighborsToRemove(node int) map[int]void {
	toRemove := make(map[int]void)
	infectedNeighbors := n.findNeighbors(node, StateI)
	for neighbor := range infectedNeighbors {
		if n.rand.Float32() < n.behavior.removeInfectedNeighborProb() {
			toRemove[neighbor] = void{}
		}
	}
	return toRemove
}

func (n *DiseasedNetwork) findNeighborsToAdd(node int) map[int]void {
	currentNeighbors := n.findNeighbors(node, -1)
	toAdd := make(map[int]void)
	for neighbor := range currentNeighbors {
		nOfn := n.findNeighbors(neighbor, -1)
		for nn := range nOfn {
			if n.rand.Float32() < n.behavior.addNeighborOfNeighborProb() {
				toAdd[nn] = void{}
			}
		}
	}
	return toAdd
}
