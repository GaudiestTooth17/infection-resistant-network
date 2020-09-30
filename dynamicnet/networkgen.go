package dynamicnet

// makeCompleteNetwork returns a complete network with the given number of nodes
func makeCompleteNetwork(numNodes int) Network {
	data := make([][]uint8, numNodes)
	for i := 0; i < numNodes; i++ {
		data[i] = make([]uint8, numNodes)
		for j := 0; j < numNodes; j++ {
			if i != j {
				data[i][j] = 1
			}
		}
	}
	return adjacencyMatrix{data: data}
}

// makeCircularNetwork makes a network resembling a ring where each node has degree 2
func makeCircularNetwork(numNodes int) Network {
	data := make([][]uint8, numNodes)
	for i := 0; i < numNodes; i++ {
		data[i] = make([]uint8, numNodes)
	}
	for i := 0; i < numNodes; i++ {
		data[i][(i+1)%numNodes] = 1
		data[(i+1)%numNodes][i] = 1
	}
	return adjacencyMatrix{data: data}
}
