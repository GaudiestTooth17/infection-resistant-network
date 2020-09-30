package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/GaudiestTooth17/infection-resistant-network/dynamicnet"
)

func readAdjacencyList(fileName string) dynamicnet.Network {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileReader := bufio.NewReader(file)

	//set up adjMatrix
	line, err := fileReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	n, err := strconv.Atoi(line[:len(line)-1])
	if err != nil {
		panic(err)
	}
	adjMatrix := makeAdjacencyMatrix(n)

	//populate adjMatrix
	for true {
		line, err = fileReader.ReadString('\n')
		// This is the exit condition
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil && !errors.Is(err, io.EOF) {
			panic(err)
		}
		i, j := lineToCoordinate(line)
		adjMatrix[i][j] = 1
	}
	return dynamicnet.NewAdjacencyMatrix(adjMatrix)
}

func makeAdjacencyMatrix(n int) [][]uint8 {
	adjMatrix := make([][]uint8, n)
	for i := 0; i < n; i++ {
		adjMatrix[i] = make([]uint8, n)
	}
	return adjMatrix
}

func lineToCoordinate(line string) (i, j int) {
	coordinates := strings.Fields(line)
	if len(coordinates) != 2 {
		panic("Error: " + line)
	}
	i, err := strconv.Atoi(coordinates[0])
	if err != nil {
		panic(err)
	}
	j, err = strconv.Atoi(coordinates[1])
	if err != nil {
		panic(err)
	}
	return i, j
}
