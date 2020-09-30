package dynamicnet

import (
	"testing"
)

func TestStep(t *testing.T) {
	numNodes := 2000
	adjMat := makeCompleteNetwork(numNodes)
	iStrat := InfectN{n: 1}
	dis := NewBasicDisease(1, 1, 1.0)
	beh := NewSimpleBehavior(1, numNodes, 0.0, 0.0)
	net := NewDiseasedNetwork(dis, adjMat, iStrat, beh)

	numInfected := len(net.FindNodesInState(StateI))
	if numInfected != 1 {
		t.Errorf("Expected 1 infected node, found %d", numInfected)
	}

	net.Step()
	numExposed := len(net.FindNodesInState(StateE))
	if numExposed != numNodes-1 {
		t.Errorf("Expected %d exposed nodes after 1 step, found %d",
			numNodes-1, numExposed)
	}

	net.Step()
	numInfected = len(net.FindNodesInState(StateI))
	if numInfected != numNodes-1 {
		t.Errorf("Expected %d infected nodes after 2 steps, found %d",
			numNodes-1, numInfected)
	}

	net.Step()
	numRemoved := len(net.FindNodesInState(StateR))
	if numRemoved != numNodes {
		t.Errorf("Expected %d removed nodes after 3 steps, found %d", numNodes, numRemoved)
	}
}
