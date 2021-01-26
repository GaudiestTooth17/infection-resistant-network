package optimized

import (
	"fmt"
	"time"

	"github.com/GaudiestTooth17/infection-resistant-network/diseasednetwork"
	dsnet "github.com/GaudiestTooth17/infection-resistant-network/diseasednetwork"
	"github.com/GaudiestTooth17/infection-resistant-network/dynamicnet"
	"github.com/GaudiestTooth17/infection-resistant-network/evolution"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// NetworkFitnessCalculator implements manager.FitnessCalculator and
// measures the fitness of an agent behavior on a network
type NetworkFitnessCalculator struct {
	network   dsnet.Network
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
	trialFitnesses := make([]float32, n.numTrials)
	fitnessChannel := make(chan FitnessData)
	for trial := 0; trial < n.numTrials; trial++ {
		network := dsnet.NewDiseasedNetwork(&n.network,
			[]dsnet.Disease{n.disease.MakeCopy()}, []dsnet.PlotMaker{})
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

// GraphAverageR0 runs a batch of simulations and then graphs R0 at each of the steps.
// It returns the average of all entries in allR0s.
func (n NetworkFitnessCalculator) GraphAverageR0(plotName string) float64 {
	allR0s := make([]plotter.XYs, n.numTrials)
	r0Channel := make(chan R0Data)
	for trial := 0; trial < n.numTrials; trial++ {
		network := dsnet.NewDiseasedNetwork(&n.network,
			[]dsnet.Disease{n.disease.MakeCopy()},
			[]dsnet.PlotMaker{dsnet.NewR0Plotter("r0", n.simLength)})
		go r0Async(r0Channel, trial, network, n.simLength)
	}
	for i := 0; i < n.numTrials; i++ {
		r0data := <-r0Channel
		allR0s[r0data.trialNumber] = r0data.plotPoints
	}

	// TODO: make graph show interquartile ranges
	// save the plot
	p, err := plot.New()
	check(err)

	p.Title.Text = plotName
	p.X.Label.Text = "Time Step"
	p.Y.Label.Text = "R0"
	points := make(plotter.XYs, n.simLength)
	// Set the points to be time step to average R0
	for j := 0; j < n.simLength; j++ {
		// calculate average R0
		avgR0 := float64(0)
		for i := 0; i < n.numTrials; i++ {
			avgR0 += allR0s[i][j].Y / float64(n.numTrials)
		}
		// add average R0 to points
		points[j].X = float64(j)
		points[j].Y = avgR0
	}

	// add the points to the plot and save the plot
	err = plotutil.AddLines(p, plotName, points)
	check(err)
	err = p.Save(8*vg.Inch, 8*vg.Inch, plotName+".png")
	check(err)

	nonZeroEntries := float64(0)
	avgR0 := float64(0)
	for _, point := range points {
		if point.Y > 0 {
			nonZeroEntries++
			avgR0 += point.Y
		}
	}
	return avgR0 / nonZeroEntries
}

// R0 of the disease that was given to the fitness calculator
func (n NetworkFitnessCalculator) R0() float64 {
	return n.r0
}

// CalcAndOutput sequentially calculates the fitness of a network
// and prints the change in states to the screen
func (n *NetworkFitnessCalculator) CalcAndOutput() float32 {
	// run simulations
	network := dsnet.NewDiseasedNetwork(&n.network,
		[]dsnet.Disease{n.disease.MakeCopy()}, []dsnet.PlotMaker{})
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

func calcAsync(outChan chan<- FitnessData, trialNumber int, network dsnet.DiseasedNetwork, numSteps int) {
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

// R0Data is used to send data about R0 over a channel
type R0Data struct {
	trialNumber int
	plotPoints  plotter.XYs
	elapsedTime time.Duration
}

func r0Async(outChan chan<- R0Data, trialNumber int, network dsnet.DiseasedNetwork, numSteps int) {
	duration := time.Duration(0)
	for i := 0; i < numSteps; i++ {
		duration += network.Step()
	}
	// Watch out! This assumes that the r0 PlotMaker is the first in the slice.
	outChan <- R0Data{trialNumber: trialNumber, plotPoints: network.PlotMakers[0].Points(),
		elapsedTime: duration}
}

// genotypeToAgentBehavior converts a Float32Genotype to an AgentBehavior
func genotypeToAgentBehavior(genotype evolution.Float32Genotype) dynamicnet.AgentBehavior {
	return dynamicnet.NewSimpleBehavior(int(genotype.Get(0)), int(genotype.Get(1)),
		genotype.Get(2), genotype.Get(3))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
