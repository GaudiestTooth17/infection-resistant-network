package diseasednetwork

// represent disease states
const (
	StateS = iota
	StateE = iota
	StateI = iota
	StateR = iota
)

// NewBasicDisease creates a new Disease with the given values
func NewBasicDisease(timeToI, timeToR int16, infectionProbability float32) Disease {
	return basicDisease{timeToI: timeToI, timeToR: timeToR,
		infectionProbability: infectionProbability}
}

// Disease represents an infection that spreads across a network
type Disease interface {
	TimeToI() int16
	TimeToR() int16
	InfectionProbability() float32
}

// basicDisease is a simple disease with constant values
type basicDisease struct {
	timeToI              int16
	timeToR              int16
	infectionProbability float32
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
func (d basicDisease) InfectionProbability() float32 {
	return d.infectionProbability
}
