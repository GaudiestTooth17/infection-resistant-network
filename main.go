package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <startingMatrixName> <calculatorConfName> <genotypeConfName>\n", os.Args[0])
		return
	}
	startingMatrixName := os.Args[1]
	calculatorConfName := os.Args[2]
	genotypeConfName := os.Args[3]

	fitnessCalculator := readFitnessCalculator(calculatorConfName, startingMatrixName)
	genotypes := readGenotypeConf(genotypeConfName)
	// fmt.Printf("Fitness Calculator: %v\n", fitnessCalculator)
	fmt.Printf("Genotypes: %v\n", genotypes)

	for i := range genotypes {
		timeStart := time.Now()
		fitness := fitnessCalculator.CalculateFitness()
		fmt.Printf("Genotype %d has fitness %f (%v).\n", i, fitness, time.Now().Sub(timeStart))
	}
}
