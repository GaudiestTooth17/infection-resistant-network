package dynamicnet

// There is currently no dynamic behavior in networks, so this test has no bearing

// import (
// 	"testing"

// 	"github.com/GaudiestTooth17/infection-resistant-network/diseasednetwork"
// 	"github.com/GaudiestTooth17/infection-resistant-network/networkgenerator"
// )

// const (
// 	stateS = diseasednetwork.StateS
// 	stateE = diseasednetwork.StateE
// 	stateI = diseasednetwork.StateI
// 	stateR = diseasednetwork.StateR
// )

// func TestImmuneNetwork(t *testing.T) {
// 	numNodes := 100
// 	net := diseasednetwork.NewDiseasedNetwork(
// 		[]diseasednetwork.Disease{diseasednetwork.NewBasicDisease(1, 1, 1.0, diseasednetwork.NewInfectN(1))},
// 		networkgenerator.MakeCompleteNetwork(numNodes))

// 	numInfected := len(net.FindNodesInState(stateI, 0))
// 	if numInfected != 1 {
// 		t.Errorf("Expected 1 infected node to start with, found %d", numInfected)
// 	}

// 	for i := 0; i < 20; i++ {
// 		net.Step()
// 		numSusceptible := len(net.FindNodesInState(stateS, 0))
// 		if numSusceptible != numNodes-1 {
// 			numExposed := len(net.FindNodesInState(stateE, 0))
// 			numInfected := len(net.FindNodesInState(stateI, 0))
// 			numRecovered := len(net.FindNodesInState(stateR, 0))
// 			t.Errorf("(step %d) Expected %d susceptible nodes, found %d.\nAlso found %d exposed, %d infected, %d recovered.\n",
// 				i, numNodes-1, numSusceptible, numExposed, numInfected, numRecovered)
// 		}
// 	}
// }
