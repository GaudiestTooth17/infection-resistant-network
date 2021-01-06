package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/GaudiestTooth17/infection-resistant-network/diseasednetwork"
)

func readAdjacencyList(fileName string) diseasednetwork.Network {
	// read a file with the first line being the number of nodes
	// the next lines are edges with from node and to node
	// then there may or may not be a blank line so that node coordinates can be recorded
	// However, the coordinates are not read.
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileReader := bufio.NewReader(file)

	// initialize network with number of nodes
	line, err := fileReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	n, err := strconv.Atoi(line[:len(line)-1])
	if err != nil {
		panic(err)
	}
	network := diseasednetwork.NewNetwork(n)

	//populate network
	for true {
		line, err = fileReader.ReadString('\n')
		// exit on EOF or if the line is length 1 (or 0) (as in it is only \n)
		if (err != nil && errors.Is(err, io.EOF)) || len(line) < 2 {
			break
		} else if err != nil && !errors.Is(err, io.EOF) {
			panic(err)
		}
		i, j := lineToCoordinate(line)
		network.AddEdge(i, j, 1)
	}
	return network
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
