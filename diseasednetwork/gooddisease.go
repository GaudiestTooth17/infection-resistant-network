package diseasednetwork

// todo: add some sort of tracking mechanism to keep track of how many nodes were benefited

// NewGoodDisease makes a positive "disease"
func NewGoodDisease(timeToR, timeToS int16, infProb float32, infStrat InitialInfectionStrategy) Disease {
	return &goodDisease{
		timeToR:              timeToR,
		timeToS:              timeToS,
		infectionProbability: infProb,
		nodeState:            make(map[int]uint8),
		timeInState:          make(map[int]int16),
		infStrat:             infStrat,
		numNodes:             -1,
	}
}

// goodDisease represents a desirable "infection" spreading across the network
// this could be information dispersion or some other sort of positive interaction
// It follows the SIRS pattern
type goodDisease struct {
	timeToR              int16
	timeToS              int16
	infectionProbability float32
	nodeState            map[int]uint8
	timeInState          map[int]int16
	infStrat             InitialInfectionStrategy
	numNodes             int
	numInfected          int
}

// InfectionProbability returns the probability that in one time step a node will infect
// its neighbor
func (d *goodDisease) InfectionProbability() float32 {
	return d.infectionProbability
}

func (d *goodDisease) State(node int) uint8 {
	return d.nodeState[node]
}

func (d *goodDisease) SetState(node int, state uint8) {
	if state == uint8(StateI) {
		d.numInfected++
	} else if state == uint8(StateR) {
		d.numInfected--
	}
	d.nodeState[node] = state
	d.timeInState[node] = 0
}

func (d *goodDisease) ResetTimeInState(node int) {
	d.timeInState[node] = 0
}

func (d *goodDisease) IncTimeInState(node int) {
	d.timeInState[node]++
}

func (d *goodDisease) TimeInState(node int) int16 {
	return d.timeInState[node]
}

func (d *goodDisease) InitialInfection() InitialInfectionStrategy {
	return d.infStrat
}

// NumNodes returns the number of nodes that the disease infects. This is used for the
// initial infection strategy. numNodes starts at -1 and if it isn't set to a nonnegative value
// this function will panic when called.
func (d *goodDisease) NumNodes() int {
	if d.numNodes < 0 {
		panic("Disease has negative number of nodes!")
	}
	return d.numNodes
}

func (d *goodDisease) SetNumNodes(n int) {
	d.numNodes = n
}

// FindNodesInState finds all the nodes in the network with the given state
func (d *goodDisease) FindNodesInState(state int) map[int]Void {
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
func (d *goodDisease) updateStates() {
	infectedNodes := d.FindNodesInState(StateI)
	recoveredNodes := d.FindNodesInState(StateR)
	for node := range infectedNodes {
		if d.TimeInState(node) == d.timeToR {
			d.SetState(node, StateR)
		}
	}
	for node := range recoveredNodes {
		if d.TimeInState(node) == d.timeToS {
			d.SetState(node, StateS)
		}
	}
}

func (d *goodDisease) Rate() float64 {
	return float64(d.numInfected) / float64(d.numNodes)
}

func (d *goodDisease) MakeCopy() Disease {
	nodeState := make(map[int]uint8)
	for node, state := range d.nodeState {
		nodeState[node] = state
	}
	timeInState := make(map[int]int16)
	for node, time := range d.timeInState {
		timeInState[node] = time
	}
	return &goodDisease{
		timeToS:              d.timeToS,
		timeToR:              d.timeToR,
		infectionProbability: d.infectionProbability,
		nodeState:            nodeState,
		timeInState:          timeInState,
		infStrat:             d.infStrat,
		numNodes:             d.numNodes,
	}
}
