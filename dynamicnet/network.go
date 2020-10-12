package dynamicnet

// Network represents an undirected graph/network
type Network struct {
	// entries in neighbors are stored with the smaller node coming first, then the larger one
	neighbors map[int]map[int]uint8
}

// NewNetwork returns a Network with the given number of nodes
func NewNetwork(numNodes int) Network {
	neighbors := make(map[int]map[int]uint8)
	for i := 0; i < numNodes; i++ {
		neighbors[i] = make(map[int]uint8)
	}
	return Network{neighbors: neighbors}
}

// MakeCopy copies the network so that different memory is used for the internal node structure
func (n *Network) MakeCopy() Network {
	neighborsCopy := make(map[int]map[int]uint8)
	for node, adjacentNodes := range n.neighbors {
		neighborsCopy[node] = make(map[int]uint8)
		for adjacentNode, weight := range adjacentNodes {
			neighborsCopy[node][adjacentNode] = weight
		}
	}
	return Network{neighbors: neighborsCopy}
}

// NeighborsOf returns the neighbors of the given node
func (n Network) NeighborsOf(node int) map[int]uint8 {
	return n.neighbors[node]
}

// NumNodes returns the number of nodes in the network
func (n Network) NumNodes() int {
	return len(n.neighbors)
}

// EdgeWeight returns the weight of the edge from node1 to node2
func (n Network) EdgeWeight(node1, node2 int) uint8 {
	return n.neighbors[node1][node2]
}

// AddEdge adds an edge between node1 and node2 with the given weight
func (n *Network) AddEdge(node1, node2 int, weight uint8) {
	n.neighbors[node1][node2] = weight
	n.neighbors[node2][node1] = weight
}

func (n *Network) removeEdge(node1, node2 int) {
	delete(n.neighbors[node1], node2)
	delete(n.neighbors[node2], node1)
}
