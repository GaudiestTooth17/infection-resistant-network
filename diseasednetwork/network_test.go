package diseasednetwork

import (
	"testing"
)

func makeCompleteDiseasedNet(numNodes, numToInfect int) DiseasedNetwork {
	net := makeCompleteNetwork(numNodes)
	dis := NewBasicDisease(0, 0, 0)
	infectionStrat := InfectN{n: numToInfect}
	diseasedNet := NewDiseasedNetwork(dis, net, infectionStrat, NewSimpleBehavior(0, numNodes, 0.0, 0.0))
	return diseasedNet
}

func makeCircularDiseasedNet(numNodes, numToInfect int) DiseasedNetwork {
	net := makeCircularNetwork(numNodes)
	dis := NewBasicDisease(0, 0, 0)
	infectionStrat := InfectN{n: numToInfect}
	diseasedNet := NewDiseasedNetwork(dis, net, infectionStrat, NewSimpleBehavior(0, numNodes, 0.0, 0.0))
	return diseasedNet
}

func TestFindNeighbors(t *testing.T) {
	numNodes := 4
	diseasedNet := makeCompleteDiseasedNet(numNodes, 0)

	for node := 0; node < numNodes; node++ {
		numNeighbors := len(diseasedNet.findNeighbors(node, -1))
		if numNeighbors != numNodes-1 {
			t.Errorf("Expected %d neighbors for node %d, found %d", numNodes-1, node, numNeighbors)
		}
	}

	diseasedNet = makeCircularDiseasedNet(numNodes, 0)
	for node := 0; node < numNodes; node++ {
		numNeighbors := len(diseasedNet.findNeighbors(node, -1))
		if numNeighbors != 2 {
			t.Errorf("Expected node %d to have 2 neighbors, has %d neighbors", node, numNeighbors)
		}
	}
}

func TestState(t *testing.T) {
	numNodes := 10
	diseasedNet := makeCompleteDiseasedNet(numNodes, 0)
	if len(diseasedNet.FindNodesInState(StateS)) != numNodes {
		t.Errorf("Expected %d susceptible nodes before changing states, found %d",
			numNodes, len(diseasedNet.FindNodesInState(StateS)))
	}

	diseasedNet.changeState(1, StateE)
	diseasedNet.changeState(2, StateI)
	diseasedNet.changeState(3, StateR)

	susceptible := diseasedNet.FindNodesInState(StateS)
	exposed := diseasedNet.FindNodesInState(StateE)
	infected := diseasedNet.FindNodesInState(StateI)
	removed := diseasedNet.FindNodesInState(StateR)
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
