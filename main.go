package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <startingMatrixName> <simConfName>\n", os.Args[0])
		return
	}
	startingMatrixName := os.Args[1]
	simConfName := os.Args[2]
	// genotypeConfName := os.Args[3]

	fitnessCalculator := readFitnessCalculator(simConfName, startingMatrixName)
	// genotypes := readGenotypeConf(genotypeConfName)
	// fmt.Printf("Fitness Calculator: %v\n", fitnessCalculator)
	// fmt.Printf("Genotypes: %v\n", genotypes)

	for i := 0; i < 1; i++ {
		timeStart := time.Now()
		fitness := fitnessCalculator.CalculateFitness()
		fmt.Printf("Trial %d: proportion of nodes still susceptible: %f (%v).\n", i, fitness, time.Now().Sub(timeStart))
	}
}
