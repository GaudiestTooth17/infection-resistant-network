package dynamicnet

import "testing"

func TestImmuneNetwork(t *testing.T) {
	numNodes := 100
	net := NewDiseasedNetwork(NewBasicDisease(1, 1, 1.0),
		makeCompleteNetwork(numNodes), InfectN{n: 1},
		NewSimpleBehavior(0, numNodes, 1.0, 0.0))

	numInfected := len(net.FindNodesInState(StateI))
	if numInfected != 1 {
		t.Errorf("Expected 1 infected node to start with, found %d", numInfected)
	}

	for i := 0; i < 20; i++ {
		net.Step()
		numSusceptible := len(net.FindNodesInState(StateS))
		if numSusceptible != numNodes-1 {
			numExposed := len(net.FindNodesInState(StateE))
			numInfected := len(net.FindNodesInState(StateI))
			numRecovered := len(net.FindNodesInState(StateR))
			t.Errorf("(step %d) Expected %d susceptible nodes, found %d.\nAlso found %d exposed, %d infected, %d recovered.\n",
				i, numNodes-1, numSusceptible, numExposed, numInfected, numRecovered)
		}
	}
}
