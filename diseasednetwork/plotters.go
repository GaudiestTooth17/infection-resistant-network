package diseasednetwork

import "gonum.org/v1/plot/plotter"

// PlotMaker observes a DiseasedNetwork and saves the information it gathers to make a plot
type PlotMaker interface {
	// feedInformation should be called for every step information should be plotted for
	feedInformation(network *DiseasedNetwork)
	// Points returns all the data points that the Plotter gathered
	Points() plotter.XYs
	// PlotName gives the name of the plot the points belong to.
	// This is useful for either labeling an entire plot or a single line on the plot.
	PlotName() string
}

type r0Plotter struct {
	points plotter.XYs
	name   string
}

// NewR0Plotter initializes an r0Plotter
func NewR0Plotter(name string, numSteps int) PlotMaker {
	return &r0Plotter{name: name, points: make(plotter.XYs, numSteps)}
}

func (p *r0Plotter) feedInformation(network *DiseasedNetwork) {
	p.points[network.stepNum] = plotter.XY{X: float64(network.stepNum), Y: network.R0(0)}
}

func (p *r0Plotter) Points() plotter.XYs {
	return p.points
}

func (p *r0Plotter) PlotName() string {
	return p.name
}
