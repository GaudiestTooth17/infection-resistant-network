package diseasednetwork

import (
	"testing"
)

// TestStep runs a few steps and tests the state of the network after each one
// numNodes can be adjusted
func TestStep(t *testing.T) {
	numNodes := 500
	adjMat := makeCompleteNetwork(numNodes)
	dis := NewBasicDisease(1, 1, 1.0, InfectN{n: 1})
	net := NewDiseasedNetwork([]Disease{dis}, adjMat)

	numInfected := len(dis.FindNodesInState(StateI))
	if numInfected != 1 {
		t.Errorf("Expected 1 infected node, found %d", numInfected)
	}

	net.Step()
	numExposed := len(dis.FindNodesInState(StateE))
	if numExposed != numNodes-1 {
		t.Errorf("Expected %d exposed nodes after 1 step, found %d",
			numNodes-1, numExposed)
	}

	net.Step()
	numInfected = len(dis.FindNodesInState(StateI))
	if numInfected != numNodes-1 {
		t.Errorf("Expected %d infected nodes after 2 steps, found %d",
			numNodes-1, numInfected)
	}

	net.Step()
	numRemoved := len(dis.FindNodesInState(StateR))
	if numRemoved != numNodes {
		t.Errorf("Expected %d removed nodes after 3 steps, found %d", numNodes, numRemoved)
	}
}
