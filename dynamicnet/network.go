package dynamicnet

// Network represents an undirected graph/network
type Network interface {
	// edgeWeight returns the weight of the edge between n1 and n2 or 0 if there is no edge
	edgeWeight(n1, n2 int) uint8
	// numNodes gives the number of nodes in the network
	numNodes() int
	// remove the edge between n1 and n2
	removeEdge(n1, n2 int)
	// add (or update) edge between n1 and n2 with the provided weight
	addEdge(n1, n2 int, weight uint8)
}

// adjacencyMatrix uses a 2d n by n slice to store edge weights
type adjacencyMatrix struct {
	data [][]uint8
}

// NewAdjacencyMatrix returns a network from the adjacency matrix
func NewAdjacencyMatrix(matrix [][]uint8) Network {
	return adjacencyMatrix{data: matrix}
}

func (a adjacencyMatrix) edgeWeight(n1, n2 int) uint8 {
	return a.data[n1][n2]
}

func (a adjacencyMatrix) numNodes() int {
	return len(a.data)
}

func (a adjacencyMatrix) removeEdge(n1, n2 int) {
	a.data[n1][n2] = 0
	a.data[n2][n1] = 0
}

func (a adjacencyMatrix) addEdge(n1, n2 int, weight uint8) {
	a.data[n1][n2] = weight
	a.data[n2][n1] = weight
}
