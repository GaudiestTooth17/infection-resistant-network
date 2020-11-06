package diseasednetwork

// makeCompleteNetwork returns a complete network with the given number of nodes
func makeCompleteNetwork(numNodes int) Network {
	network := NewNetwork(numNodes)
	for i := 0; i < numNodes; i++ {
		for j := i + 1; j < numNodes; j++ {
			network.AddEdge(i, j, 1)
		}
	}
	return network
}

// makeCircularNetwork makes a network resembling a ring where each node has degree 2
func makeCircularNetwork(numNodes int) Network {
	network := NewNetwork(numNodes)
	for i := 0; i < numNodes; i++ {
		network.AddEdge(i, (i+1)%numNodes, 1)
	}
	return network
}
