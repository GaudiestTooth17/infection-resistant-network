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
	R0() float64
	ReportInfections(node int, numInfections uint)
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
	numNodesInfectedBy   []uint
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

// SetNumNodes MUST be called before actually using the disease. Unfortunately, because of the
// way the structs relate to each other, it isn't possible to specify the number of nodes when
// the disease is created.
func (d *basicDisease) SetNumNodes(n int) {
	d.numNodes = n
	d.numNodesInfectedBy = make([]uint, d.numNodes)
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

// R0 calculates the R0 of the disease
func (d *basicDisease) R0() float64 {
	numSpreaders := uint(0)
	numInfectedBySpreaders := uint(0)
	for _, numInfected := range d.numNodesInfectedBy {
		if numInfected > 0 {
			numSpreaders++
			numInfectedBySpreaders += numInfected
		}
	}
	// If no new nodes actually caught the disease, return 0 instead of NaN
	if numSpreaders == 0 {
		return 0
	}
	return float64(numInfectedBySpreaders) / float64(numSpreaders)
}

func (d *basicDisease) ReportInfections(node int, numInfections uint) {
	// fmt.Fprintf(os.Stderr, "%p\n", d)
	d.numNodesInfectedBy[node] += numInfections
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
