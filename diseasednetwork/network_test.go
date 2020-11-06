package diseasednetwork

import (
	"testing"
)

func makeCompleteDiseasedNet(numNodes, numToInfect int) (DiseasedNetwork, Disease) {
	net := makeCompleteNetwork(numNodes)
	dis := NewBasicDisease(0, 0, 0, InfectN{n: numToInfect})
	diseasedNet := NewDiseasedNetwork([]Disease{dis}, net)
	return diseasedNet, dis
}

func makeCircularDiseasedNet(numNodes, numToInfect int) (DiseasedNetwork, Disease) {
	net := makeCircularNetwork(numNodes)
	dis := NewBasicDisease(0, 0, 0, InfectN{n: numToInfect})
	diseasedNet := NewDiseasedNetwork([]Disease{dis}, net)
	return diseasedNet, dis
}

func TestFindNeighbors(t *testing.T) {
	numNodes := 4
	diseasedNet, _ := makeCompleteDiseasedNet(numNodes, 0)

	for node := 0; node < numNodes; node++ {
		numNeighbors := len(diseasedNet.findNeighbors(node, -1, 0))
		if numNeighbors != numNodes-1 {
			t.Errorf("Expected %d neighbors for node %d, found %d", numNodes-1, node, numNeighbors)
		}
	}

	diseasedNet, _ = makeCircularDiseasedNet(numNodes, 0)
	for node := 0; node < numNodes; node++ {
		numNeighbors := len(diseasedNet.findNeighbors(node, -1, 0))
		if numNeighbors != 2 {
			t.Errorf("Expected node %d to have 2 neighbors, has %d neighbors", node, numNeighbors)
		}
	}
}

func TestState(t *testing.T) {
	numNodes := 10
	_, dis := makeCompleteDiseasedNet(numNodes, 0)
	if len(dis.FindNodesInState(StateS)) != numNodes {
		t.Errorf("Expected %d susceptible nodes before changing states, found %d",
			numNodes, len(dis.FindNodesInState(StateS)))
	}

	dis.SetState(1, StateE)
	dis.SetState(2, StateI)
	dis.SetState(3, StateR)

	susceptible := dis.FindNodesInState(StateS)
	exposed := dis.FindNodesInState(StateE)
	infected := dis.FindNodesInState(StateI)
	removed := dis.FindNodesInState(StateR)
	if len(susceptible) != 7 {
		t.Errorf("Expected 7 susceptible nodes, found %d", len(susceptible))
	}
	if len(exposed) != 1 {
		t.Errorf("Expected 1 exposed node, found %d", len(exposed))
	}
	if len(infected) != 1 {
		t.Errorf("Expected 1 infected node, found %d", len(infected))
	}
	if len(removed) != 1 {
		t.Errorf("Expected 1 removed node, found %d", len(removed))
	}
}
