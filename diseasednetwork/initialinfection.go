package diseasednetwork

import (
	"math/rand"
	"time"
)

// InitialInfectionStrategy is used to seed a DiseasedNetwork with some number of infected nodes
type InitialInfectionStrategy interface {
	apply(disease Disease)
}

// InfectN chooses n random nodes to infect
type InfectN struct {
	n int
}

// NewInfectN returns an instance of the InfectN infection strategy
func NewInfectN(n int) InfectN {
	return InfectN{n: n}
}

// apply the infection strategy to a DiseadedNetwork
func (i InfectN) apply(disease Disease) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	infectedNodes := make(map[int]bool)
	for len(infectedNodes) < i.n {
		nodeToInfect := rand.Intn(disease.NumNodes())
		if !infectedNodes[nodeToInfect] {
			infectedNodes[nodeToInfect] = true
			disease.SetState(nodeToInfect, StateI)
		}
	}
}
