package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GaudiestTooth17/infection-resistant-network/optimized"
)

func main() {
	if len(os.Args) == 3 {
		runWithVis()
	} else if len(os.Args) == 5 {
		runBatchAndGraphR0s()
	} else {
		fmt.Printf("Usage: %s <disease-file> <matrix-file> [num-sims] [sim-length]\n", os.Args[0])
		return
	}
}

// runWithVis will run one simulation and print the node states to stdout so that
// graph-visualizer can be used to graphically inspect the simulation
func runWithVis() {
	diseaseName := os.Args[1]
	matrixName := os.Args[2]
	network := readAdjacencyList(matrixName)
	disease := readDisease(diseaseName)
	fitnessCalculator := optimized.NewNetworkFitnessCalculator(network, 100, 100, disease)

	timeStart := time.Now()
	fitness := fitnessCalculator.CalcAndOutput()
	fmt.Fprintf(os.Stderr, "Proportion of nodes still susceptible: %f R0: %f (%v).\n",
		fitness, fitnessCalculator.R0(), time.Now().Sub(timeStart))
}

// runBatch runs a batch of simulations and reports the average number of nodes
// that were left susceptible at the end of each
func runBatch() {
	diseaseName := os.Args[1]
	matrixName := os.Args[2]
	disease := readDisease(diseaseName)
	network := readAdjacencyList(matrixName)

	numSims, err := strconv.Atoi(os.Args[3])
	check(err)
	simLength, err := strconv.Atoi(os.Args[4])
	check(err)

	fitnessCalculator := optimized.NewNetworkFitnessCalculator(network, numSims, simLength, disease)
	timeStart := time.Now()
	fitness := fitnessCalculator.CalculateFitness()
	fmt.Printf("Proportion of nodes still susceptible: %f (%v).\n",
		fitness, time.Now().Sub(timeStart))
}

// runBatchAndGraphR0s runs a batch of simulations and reports the average R0 of the disease
// it also saves a png showing a plot of R0 with respect to time step.
func runBatchAndGraphR0s() {
	diseaseName := os.Args[1]
	networkName := os.Args[2]
	disease := readDisease(diseaseName)
	network := readAdjacencyList(networkName)

	numSims, err := strconv.Atoi(os.Args[3])
	check(err)
	simLength, err := strconv.Atoi(os.Args[4])
	check(err)

	fitnessCalculator := optimized.NewNetworkFitnessCalculator(network, numSims, simLength, disease)
	timeStart := time.Now()
	plotName := "R0s from " + noExt(diseaseName) + " on " + noExt(networkName)
	averageR0 := fitnessCalculator.GraphAverageR0(plotName)
	fmt.Printf("Average R0: %f (%v).\n", averageR0, time.Now().Sub(timeStart))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// noExt semoves the file extension from a string or does nothing if there is no file extension
func noExt(str string) string {
	return strings.Split(str, ".")[0]
}
