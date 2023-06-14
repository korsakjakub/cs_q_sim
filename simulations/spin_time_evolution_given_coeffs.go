package simulations

import (
	"fmt"
	"math"
	"time"

	cs "github.com/korsakjakub/cs_q_sim/pkg/cs_q_sim"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/plotter"
)

func SpinTimeEvolutionSelectedCoeffs(conf cs.Config) {
	var bath []cs.State
	conf.Physics.BathCount = len(conf.Physics.InitialKet) - 1
	bc := conf.Physics.BathCount
	timeRange := conf.Physics.TimeRange
	spin := conf.Physics.Spin
	initialKet := mat.NewVecDense(len(conf.Physics.InitialKet), cs.ManyBodyVector(conf.Physics.InitialKet, int(2*spin+1)))
	downSpins := downSpins(conf.Physics.InitialKet)
	observables := prepareObservables(conf.Physics, downSpins)

	if conf.Verbosity == "debug" {
		fmt.Println("Calculating initial states...")
	}
	start := time.Now()
	for i := 0; i < bc; i += 1 {
		bath = append(bath, cs.State{Angle: cs.PolarAngleCos(i, conf.Physics), Distance: 1e3})
	}

	s := &cs.System{
		CentralSpin:   cs.State{Angle: 0.0, Distance: 0.0},
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	if conf.Verbosity == "debug" {
		fmt.Println("Preparing the Hamiltonian...")
	}
	b := conf.Physics.BathMagneticField
	b0 := conf.Physics.CentralMagneticField

	if conf.Verbosity == "debug" {
		fmt.Println("Diagonalizing...")
	}
	eigen := s.Diagonalize(s.Hamiltonian(b0, b))

	gramMatrix := cs.Grammian(initialKet, eigen.EigenVectors)

	start_time := start.Format(time.RFC3339)

	if conf.Verbosity == "debug" {
		fmt.Println("Calculating time evolution...")
	}
	xyss := make([]plotter.XYs, len(observables))
	for i, observable := range observables {
		var xys plotter.XYs
		for t := 0; t < timeRange; t += 1 {
			time := conf.Physics.Dt * float64(t)
			if conf.Verbosity == "debug" {
				fmt.Printf("t= %.4f\t(%.2f%%)\n", time, 100.0*float64(t)/float64(timeRange))
			}
			xys = append(xys, plotter.XY{X: time / (2.0 * math.Pi), Y: observable.ExpectationValue(cs.Evolve(initialKet, time, eigen.EigenValues, eigen.EigenVectors, gramMatrix))})
		}
		xyss[i] = xys
	}

	if conf.Verbosity == "debug" {
		fmt.Println("Wrapping up...")
	}
	elapsed_time := time.Since(start)
	r := cs.ResultsIO{
		Filename: start_time,
		Metadata: cs.Metadata{
			Date:           start_time,
			Simulation:     "Central spin expectation value time evolution",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsed_time.String(),
		},
		Values: struct {
			System cs.System "mapstructure:\"system\""
		}{
			System: *s,
		},
		XYs: xyss,
	}
	r.Write(conf.Files)
}
