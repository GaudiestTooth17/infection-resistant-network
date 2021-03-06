package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/GaudiestTooth17/infection-resistant-network/evolution"
)

// readGenotypeConf - Reads a csv with genotype values in it.
// The column headers should be minConnections, maxConnections, removeInfectedNeighborProb, addNeighborOfNeighborProb.
// The headers won't actually be read, but the data will be assigned in that order.
func readGenotypeConf(filename string) []evolution.Float32Genotype {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	genotypes := make([]evolution.Float32Genotype, 0)

	for ok := fileScanner.Scan(); ok; ok = fileScanner.Scan() {
		line := fileScanner.Text()
		if !unicode.IsDigit(rune(line[0])) {
			continue
		}
		genotypes = append(genotypes, lineToGenotype(line))
	}

	return genotypes
}

// lineToGenotype - line should be of the form NUM,NUM,NUM,NUM
func lineToGenotype(line string) evolution.Float32Genotype {
	numbers := strings.Split(line, ",")
	if len(numbers) != 4 {
		panic("Incorrect number of fields!")
	}
	floats := make([]float32, 4)
	for i := range floats {
		float, err := strconv.ParseFloat(numbers[i], 32)
		if err != nil {
			panic(err)
		}
		floats[i] = float32(float)
	}
	return evolution.NewFloat32Genotype(floats)
}
