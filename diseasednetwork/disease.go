package diseasednetwork

// todo: add a rating function to basicDisease

// represent disease states
const (
	StateS = iota
	StateE = iota
	StateI = iota
	StateR = iota
)

// NewBasicDisease creates a new Disease with the given values
func NewBasicDisease(timeToI, timeToR int16, infectionProbability float32,
	infectionStrategy InitialInfectionStrategy) Disease {
	return &basicDisease{
		timeToI:              timeToI,
		timeToR:              timeToR,
		infectionProbability: infectionProbability,
		nodeState:            make(map[int]uint8),
		timeInState:          make(map[int]int16),
		infStrat:             infectionStrategy,
		numNodes:             -1,
	}
}

// Disease represents an infection that spreads across a network
type Disease interface {
	InfectionProbability() float32
	State(node int) uint8
	SetState(node int, state uint8)
	ResetTimeInState(node int)
	IncTimeInState(node int)
	TimeInState(node int) int16
	InitialInfection() InitialInfectionStrategy
	SetNumNodes(n int)
	NumNodes() int
	FindNodesInState(state int) map[int]Void
	MakeCopy() Disease
	updateStates()
	Rate() float64
}

// basicDisease is a simple disease with constant values
type basicDisease struct {
	timeToI              int16
	timeToR              int16
	infectionProbability float32
	nodeState            map[int]uint8
	timeInState          map[int]int16
	infStrat             InitialInfectionStrategy
	numNodes             int
}

// InfectionProbability returns the probability that in one time step a node will infect
// its neighbor
func (d *basicDisease) InfectionProbability() float32 {
	return d.infectionProbability
}

func (d *basicDisease) State(node int) uint8 {
	return d.nodeState[node]
}

func (d *basicDisease) SetState(node int, state uint8) {
	d.nodeState[node] = state
	d.timeInState[node] = 0
}

func (d *basicDisease) ResetTimeInState(node int) {
	d.timeInState[node] = 0
}

func (d *basicDisease) IncTimeInState(node int) {
	d.timeInState[node]++
}

func (d *basicDisease) TimeInState(node int) int16 {
	return d.timeInState[node]
}

func (d *basicDisease) InitialInfection() InitialInfectionStrategy {
	return d.infStrat
}

// NumNodes returns the number of nodes that the disease infects. This is used for the
// initial infection strategy. numNodes starts at -1 and if it isn't set to a nonnegative value
// this function will panic when called.
func (d *basicDisease) NumNodes() int {
	if d.numNodes < 0 {
		panic("Disease has negative number of nodes!")
	}
	return d.numNodes
}

func (d *basicDisease) SetNumNodes(n int) {
	d.numNodes = n
}

// FindNodesInState finds all the nodes in the network with the given state
func (d *basicDisease) FindNodesInState(state int) map[int]Void {
	s := uint8(state)
	nodes := make(map[int]Void)
	for node, st := range d.nodeState {
		if s == st {
			nodes[node] = Void{}
		}
	}
	return nodes
}

// updateStates updates the states and the time in state for all the nodes
func (d *basicDisease) updateStates() {
	exposedNodes := d.FindNodesInState(StateE)
	infectedNodes := d.FindNodesInState(StateI)
	for node := range exposedNodes {
		if d.TimeInState(node) == d.timeToI {
			d.SetState(node, StateI)
		}
	}
	for node := range infectedNodes {
		if d.TimeInState(node) == d.timeToR {
			d.SetState(node, StateR)
		}
	}
}

func (d *basicDisease) Rate() float64 {
	susceptibleNodes := len(d.FindNodesInState(StateS))
	// exposedNodes := len(network.FindNodesInState(diseasednetwork.StateE))
	// infectedNodes := len(network.FindNodesInState(diseasednetwork.StateI))
	// removedNodes := len(network.FindNodesInState(diseasednetwork.StateR))
	// fmt.Printf("%d S, %d E, %d I, %d R\n", susceptibleNodes, infectedNodes, exposedNodes, removedNodes)
	totalNodes := d.NumNodes()
	return float64(susceptibleNodes) / float64(totalNodes)
}

/*
type basicDisease struct {
	timeToI              int16
	timeToR              int16
	infectionProbability float32
	nodeState            map[int]uint8
	timeInState          map[int]int16
	infStrat             InitialInfectionStrategy
	numNodes             int
}
*/

func (d *basicDisease) MakeCopy() Disease {
	nodeState := make(map[int]uint8)
	for node, state := range d.nodeState {
		nodeState[node] = state
	}
	timeInState := make(map[int]int16)
	for node, time := range d.timeInState {
		timeInState[node] = time
	}
	return &basicDisease{
		timeToI:              d.timeToI,
		timeToR:              d.timeToR,
		infectionProbability: d.infectionProbability,
		nodeState:            nodeState,
		timeInState:          timeInState,
		infStrat:             d.infStrat,
		numNodes:             d.numNodes,
	}
}
