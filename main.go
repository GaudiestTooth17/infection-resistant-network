package main

import (
	"fmt"
	"os"
	"time"

	dsnet "github.com/GaudiestTooth17/infection-resistant-network/diseasednetwork"

	"github.com/GaudiestTooth17/infection-resistant-network/optimized"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <startingMatrixName>\n", os.Args[0])
		return
	}
	startingMatrixName := os.Args[1]
	// simConfName := os.Args[2]
	// genotypeConfName := os.Args[3]

	// fitnessCalculator := readFitnessCalculator(simConfName, startingMatrixName)
	network := readAdjacencyList(startingMatrixName)
	disease := dsnet.NewBasicDisease(4, 7, .02, dsnet.NewInfectN(10))
	fitnessCalculator := optimized.NewNetworkFitnessCalculator(network, 100, 100, disease)
	// genotypes := readGenotypeConf(genotypeConfName)
	// fmt.Printf("Fitness Calculator: %v\n", fitnessCalculator)
	// fmt.Printf("Genotypes: %v\n", genotypes)

	for i := 0; i < 1; i++ {
		timeStart := time.Now()
		fitness := fitnessCalculator.CalculateFitness()
		fmt.Printf("Trial %d: proportion of nodes still susceptible: %f (%v).\n", i, fitness, time.Now().Sub(timeStart))
	}
}
