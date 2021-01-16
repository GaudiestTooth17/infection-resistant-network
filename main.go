package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/GaudiestTooth17/infection-resistant-network/optimized"
)

func main() {
	if len(os.Args) == 3 {
		runWithVis()
	} else if len(os.Args) == 5 {
		runBatch()
	} else {
		fmt.Printf("Usage: %s <disease-file> <matrix-file> [num-sims] [sim-length]\n", os.Args[0])
		return
	}
}

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

func check(err error) {
	if err != nil {
		panic(err)
	}
}
