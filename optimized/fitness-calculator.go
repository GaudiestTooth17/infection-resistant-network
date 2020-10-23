package optimized

import (
	"time"

	"github.com/GaudiestTooth17/infection-resistant-network/diseasednetwork"
	"github.com/GaudiestTooth17/infection-resistant-network/dynamicnet"
	"github.com/GaudiestTooth17/infection-resistant-network/evolution"
)

// NetworkFitnessCalculator implements manager.FitnessCalculator and
// measures the fitness of an agent behavior on a network
type NetworkFitnessCalculator struct {
	network           diseasednetwork.Network
	numTrials         int
	simLength         int
	disease           diseasednetwork.Disease
	infectionStrategy diseasednetwork.InitialInfectionStrategy
}

// NewNetworkFitnessCalculator creates a NetworkFitnessCalculator with the provided values
func NewNetworkFitnessCalculator(network diseasednetwork.Network, numTrials, simLength int, infStrat diseasednetwork.InitialInfectionStrategy, disease diseasednetwork.Disease) NetworkFitnessCalculator {
	return NetworkFitnessCalculator{
		network:           network,
		numTrials:         numTrials,
		simLength:         simLength,
		infectionStrategy: infStrat,
		disease:           disease,
	}
}

// CalculateFitness - Calculate how fit the parameters are as agent behaviors for a DiseasedNetwork
func (n NetworkFitnessCalculator) CalculateFitness() float32 {
	// consider a fitness function that rewards spreading a positive infection while penalizing spreading a negative infection
	// fewer disconnected components is a plus, infected nodes is a minus
	// the key is to preserve the good a network serves
	// in the future, add ability for agents to change behavior over time

	trialFitnesses := make([]float32, n.numTrials)
	fitnessChannel := make(chan FitnessData)
	for trial := 0; trial < n.numTrials; trial++ {
		network := diseasednetwork.NewDiseasedNetwork(n.disease, n.network, n.infectionStrategy)
		go calcAsync(fitnessChannel, trial, network, n.simLength)
	}
	for i := 0; i < n.numTrials; i++ {
		fData := <-fitnessChannel
		trialFitnesses[fData.trialNumber] = fData.fitness
	}

	totalFitness := float32(0)
	for _, fitness := range trialFitnesses {
		totalFitness += fitness / float32(n.numTrials)
	}
	return totalFitness
}

// FitnessData conveys data about a fitness calculation over a channel
type FitnessData struct {
	trialNumber int
	fitness     float32
	elapsedTime time.Duration
}

func calcAsync(outChan chan<- FitnessData, trialNumber int, network diseasednetwork.DiseasedNetwork, numSteps int) {
	totalDuration := time.Duration(0)
	for i := 0; i < numSteps; i++ {
		totalDuration += network.Step()
	}
	fitness := rateNetwork(network)
	outChan <- FitnessData{trialNumber: trialNumber, fitness: fitness, elapsedTime: totalDuration}
}

func rateNetwork(network diseasednetwork.DiseasedNetwork) float32 {
	susceptibleNodes := len(network.FindNodesInState(diseasednetwork.StateS))
	// exposedNodes := len(network.FindNodesInState(diseasednetwork.StateE))
	// infectedNodes := len(network.FindNodesInState(diseasednetwork.StateI))
	// removedNodes := len(network.FindNodesInState(diseasednetwork.StateR))
	// fmt.Printf("%d S, %d E, %d I, %d R\n", susceptibleNodes, infectedNodes, exposedNodes, removedNodes)
	totalNodes := network.NumNodes()
	return float32(susceptibleNodes) / float32(totalNodes)
}

// genotypeToAgentBehavior converts a Float32Genotype to an AgentBehavior
func genotypeToAgentBehavior(genotype evolution.Float32Genotype) dynamicnet.AgentBehavior {
	return dynamicnet.NewSimpleBehavior(int(genotype.Get(0)), int(genotype.Get(1)),
		genotype.Get(2), genotype.Get(3))
}
