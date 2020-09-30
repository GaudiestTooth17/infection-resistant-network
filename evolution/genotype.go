package evolution

// Float32Genotype is used for applications with parameters that need to be tweaked
type Float32Genotype struct {
	floats []float32
}

// NewFloat32Genotype creates a new genotype from the float slice
func NewFloat32Genotype(floats []float32) Float32Genotype {
	return Float32Genotype{floats: floats}
}

// Len gives the number of floats in the genotype
func (f *Float32Genotype) Len() int {
	return len(f.floats)
}

// Get returns the float32 at i
func (f *Float32Genotype) Get(i int) float32 {
	return f.floats[i]
}
