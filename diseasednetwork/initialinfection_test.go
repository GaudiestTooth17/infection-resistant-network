package diseasednetwork

import (
	"testing"
)

func TestInfectN(t *testing.T) {
	numNodes := 1000
	for i := range [5]int{1, 25, 100, 500, 1000} {
		net := makeCircularDiseasedNet(numNodes, i)
		numInfected := len(net.FindNodesInState(StateI))
		if numInfected != i {
			t.Errorf("Expected %d infected, found %d", i, numInfected)
		}
	}
}
