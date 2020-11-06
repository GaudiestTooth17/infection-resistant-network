package diseasednetwork

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
	TimeToI() int16
	TimeToR() int16
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

// TimeToI returns the time it takes to change from the Exposed State to the Infectious State
func (d basicDisease) TimeToI() int16 {
	return d.timeToI
}

// TimeToR returns the time it takes to change from the Infectious State to the Removed State
func (d basicDisease) TimeToR() int16 {
	return d.timeToR
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
