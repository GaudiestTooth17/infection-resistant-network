package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/GaudiestTooth17/infection-resistant-network/diseasednetwork"
)

func readDisease(fileName string) diseasednetwork.Disease {
	diseaseFile, err := os.Open(fileName)
	check(err)
	defer diseaseFile.Close()
	reader := bufio.NewReader(diseaseFile)

	diseaseLine, err := reader.ReadString('\n')
	check(err)

	return parseDisease(diseaseLine)
}

// parseDisease parameters from line
func parseDisease(line string) diseasednetwork.Disease {
	fields := strings.Fields(line)
	if len(fields) != 4 {
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
	} else if infectionProbability < 0 || infectionProbability > 1.0 {
		panic("infectionProbability must be at least 0 and at most 1.")
	}

	numberToInfectAtStart, err := strconv.Atoi(fields[3])
	return diseasednetwork.NewBasicDisease(int16(timeToI), int16(timeToR), float32(infectionProbability),
		diseasednetwork.NewInfectN(numberToInfectAtStart))
}
