package main

import (
	"fmt"
	"os"
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

	for i, genotype := range genotypes {
		fitness := fitnessCalculator.CalculateFitness(genotype)
		fmt.Printf("Genotype %d has fitness %f\n", i, fitness)
	}
}
