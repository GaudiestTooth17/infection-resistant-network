package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/GaudiestTooth17/infection-resistant-network/dynamicnet"

	"github.com/GaudiestTooth17/infection-resistant-network/optimized"
)

// readFitnessCalculator reads a text file containing a text description of a fitness calculator
func readFitnessCalculator(fitnessCalcFilename, adjListFilename string) optimized.NetworkFitnessCalculator {
	adjacencyMatrix := readAdjacencyList(adjListFilename)
	fitnessCalcFile, err := os.Open(fitnessCalcFilename)
	if err != nil {
		panic(err)
	}
	defer fitnessCalcFile.Close()
	reader := bufio.NewReader(fitnessCalcFile)

	// there are four values that need to be retrieved
	// it will ignore any labels, but they should be used for clarity!
	// skip over the plaintext header containing the labels
	var line string
	for true {
		line, err = reader.ReadString('\n')
		if unicode.IsDigit(rune(line[0])) || err != nil {
			break
		}
		// do nothing...
	}
	// check to make sure err is still nil before proceeding
	if err != nil {
		panic(err)
	}

	var numTrials int
	var simLength int
	var disease dynamicnet.Disease
	var infectionStrategy dynamicnet.InitialInfectionStrategy
	numTrials, err = strconv.Atoi(line[:len(line)-1])
	if err != nil {
		panic(err)
	}
	for i := 0; i < 3; i++ {
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			panic(err)
		}
		line = line[:len(line)-1]
		if i == 0 {
			simLength, err = strconv.Atoi(line)
		} else if i == 1 {
			disease = parseDisease(line)
		} else if i == 2 {
			infectionStrategy = parseInfectionStrategy(line)
		}
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Printf("Problem with value %d\n", i)
			fmt.Println(line)
			panic(err)
		}
	}
	fmt.Printf("num trials: %d\n", numTrials)
	return optimized.NewNetworkFitnessCalculator(adjacencyMatrix, numTrials, simLength, infectionStrategy, disease)
}

// parseDisease parameters from line
func parseDisease(line string) dynamicnet.Disease {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		fmt.Printf("%v\n", fields)
		panic("Expected three values for disease description!")
	}
	timeToI, err := strconv.Atoi(fields[0])
	if err != nil {
		panic(err)
	} else if timeToI > math.MaxInt16 || timeToI < 0 {
		panic("timeToI must be in the range of a 16 bit int and non-negative.")
	}
	timeToR, err := strconv.Atoi(fields[1])
	if err != nil {
		panic(err)
	} else if timeToR > math.MaxInt16 || timeToR < 0 {
		panic("timeToR must be in the range of a 16 bit int and non-negative.")
	}
	infectionProbability, err := strconv.ParseFloat(fields[2], 32)
	if err != nil {
		panic(err)
	}
	return dynamicnet.NewBasicDisease(int16(timeToI), int16(timeToR), float32(infectionProbability))
}

// parseInfectionStrategy parameters from line
func parseInfectionStrategy(line string) dynamicnet.InitialInfectionStrategy {
	n, err := strconv.Atoi(line)
	if err != nil {
		panic(err)
	}
	return dynamicnet.NewInfectN(n)
}
