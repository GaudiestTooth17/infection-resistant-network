package optimized

import (
	"fmt"
	"time"

	"github.com/GaudiestTooth17/infection-resistant-network/diseasednetwork"
	dsnet "github.com/GaudiestTooth17/infection-resistant-network/diseasednetwork"
	"github.com/GaudiestTooth17/infection-resistant-network/dynamicnet"
	"github.com/GaudiestTooth17/infection-resistant-network/evolution"
)

// NetworkFitnessCalculator implements manager.FitnessCalculator and
// measures the fitness of an agent behavior on a network
type NetworkFitnessCalculator struct {
	network   diseasednetwork.Network
	numTrials int
	simLength int
	disease   diseasednetwork.Disease
	r0        float64
}

// NewNetworkFitnessCalculator creates a NetworkFitnessCalculator with the provided values
func NewNetworkFitnessCalculator(network dsnet.Network, numTrials, simLength int, disease dsnet.Disease) NetworkFitnessCalculator {
	return NetworkFitnessCalculator{
		network:   network,
		numTrials: numTrials,
		simLength: simLength,
		disease:   disease,
		r0:        -1,
	}
}

// CalculateFitness - Calculate how fit the parameters are as agent behaviors for a DiseasedNetwork
func (n NetworkFitnessCalculator) CalculateFitness() float32 {
	// TODO: add R0 calculation

	trialFitnesses := make([]float32, n.numTrials)
	fitnessChannel := make(chan FitnessData)
	for trial := 0; trial < n.numTrials; trial++ {
		network := dsnet.NewDiseasedNetwork([]dsnet.Disease{n.disease.MakeCopy()}, n.network)
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

// R0 of the disease that was given to the fitness calculator
func (n NetworkFitnessCalculator) R0() float64 {
	return n.r0
}

// CalcAndOutput sequentially calculates the fitness of a network
// and prints the change in states to the screen
func (n *NetworkFitnessCalculator) CalcAndOutput() float32 {
	// run simulations
	network := dsnet.NewDiseasedNetwork([]dsnet.Disease{n.disease.MakeCopy()}, n.network)
	printStates(network.GetNodeStates(0))

	// run simulation
	for len(network.FindNodesInState(dsnet.StateE, 0))+len(network.FindNodesInState(dsnet.StateI, 0)) > 0 {
		network.Step()
		printStates(network.GetNodeStates(0))
	}
	n.r0 = network.R0(0)

	fmt.Println("end")
	return rateNetwork(network)
}

// For each node with a different state, prints "<node> <state>\n" to stdout. Finishes with a newline.
func printStates(states []uint8) {
	for node, state := range states {
		fmt.Printf("%d %d\n", node, state)
	}
	fmt.Println()
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
	susceptibleNodes := len(network.FindNodesInState(diseasednetwork.StateS, 0))
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
