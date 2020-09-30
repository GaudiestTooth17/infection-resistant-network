package evolution

// PopulationManager searches the solution space given by its calculator
type PopulationManager struct {
	population []Float32Genotype
}

// FitnessCalculator takes a genotype and gives it a fitness rating
type FitnessCalculator interface {
	// CalculateFitness gives a rating between 0 and 1 for how good the provided
	// genotype/solution is
	CalculateFitness(genotype Float32Genotype) float32
}
