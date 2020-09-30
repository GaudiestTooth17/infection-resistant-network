package dynamicnet

// AgentBehavior defines the way agents in the network react to the spreading infection
type AgentBehavior interface {
	removeInfectedNeighborProb() float32
	minConnections() int
	maxConnections() int
	addNeighborOfNeighborProb() float32
}

type simpleBehavior struct {
	minConn        int
	maxConn        int
	removeInfNProb float32
	addNofNProb    float32
}

// NewSimpleBehavior creates a behavior with unvarying values
func NewSimpleBehavior(minConnections, maxConnections int,
	removeInfectedNeighborProb, addNeighborOfNeighborProb float32) AgentBehavior {
	return simpleBehavior{
		minConn:        minConnections,
		maxConn:        maxConnections,
		removeInfNProb: removeInfectedNeighborProb,
		addNofNProb:    addNeighborOfNeighborProb,
	}
}

func (s simpleBehavior) removeInfectedNeighborProb() float32 {
	return s.removeInfNProb
}

func (s simpleBehavior) minConnections() int {
	return s.minConn
}

func (s simpleBehavior) maxConnections() int {
	return s.maxConn
}

func (s simpleBehavior) addNeighborOfNeighborProb() float32 {
	return s.addNofNProb
}
