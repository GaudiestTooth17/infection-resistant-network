package optimized

import (
	"fmt"

	"github.com/GaudiestTooth17/infection-resistant-network/dynamicnet"
	"github.com/GaudiestTooth17/infection-resistant-network/evolution"
)

// NetworkFitnessCalculator implements manager.FitnessCalculator and
// measures the fitness of an agent behavior on a network
type NetworkFitnessCalculator struct {
	network           dynamicnet.Network
	numTrials         int
	simLength         int
	disease           dynamicnet.Disease
	infectionStrategy dynamicnet.InitialInfectionStrategy
}

// NewNetworkFitnessCalculator creates a NetworkFitnessCalculator with the provided values
func NewNetworkFitnessCalculator(network dynamicnet.Network, numTrials, simLength int, infStrat dynamicnet.InitialInfectionStrategy, disease dynamicnet.Disease) NetworkFitnessCalculator {
	return NetworkFitnessCalculator{
		network:           network,
		numTrials:         numTrials,
		simLength:         simLength,
		infectionStrategy: infStrat,
		disease:           disease,
	}
}

// CalculateFitness - Calculate how fit the parameters are as agent behaviors for a DiseasedNetwork
func (n NetworkFitnessCalculator) CalculateFitness(genotype evolution.Float32Genotype) float32 {
	behavior := genotypeToAgentBehavior(genotype)

	trialFitnesses := make([]float32, n.numTrials)
	for trial := 0; trial < n.numTrials; trial++ {
		fmt.Printf("trial %d\n", trial)
		network := dynamicnet.NewDiseasedNetwork(n.disease, n.network, n.infectionStrategy, behavior)
		for step := 0; step < n.simLength; step++ {
			fmt.Printf("step %d\n", step)
			network.Step()
		}
		fmt.Println()
		trialFitnesses[trial] = float32(network.NumNodes()-len(network.FindNodesInState(dynamicnet.StateI))) / float32(network.NumNodes())
	}

	totalFitness := float32(0)
	for _, fitness := range trialFitnesses {
		totalFitness += fitness / float32(n.numTrials)
	}
	return totalFitness
}

// genotypeToAgentBehavior converts a Float32Genotype to an AgentBehavior
func genotypeToAgentBehavior(genotype evolution.Float32Genotype) dynamicnet.AgentBehavior {
	return dynamicnet.NewSimpleBehavior(int(genotype.Get(0)), int(genotype.Get(1)),
		genotype.Get(2), genotype.Get(3))
}
