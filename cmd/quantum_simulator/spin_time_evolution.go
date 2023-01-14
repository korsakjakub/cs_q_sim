package main

import (
	"math"
	"time"

	au "github.com/korsakjakub/cs_q_sim/internal/analysis_utilities"
	hs "github.com/korsakjakub/cs_q_sim/internal/hilbert_space"
	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/plotter"
)

func spin_time_evolution(conf qs.Config) {
	cs := qs.State{Angle: 0.0, Distance: 0.0}
	var bath []qs.State
	bc := len(conf.Physics.SpinEvolutionConfig.InitialKet) - 1
	timeRange := conf.Physics.SpinEvolutionConfig.TimeRange
	spin := conf.Physics.Spin
	initialKet := hs.NewKetReal(hs.ManyBodyVector(conf.Physics.SpinEvolutionConfig.InitialKet, int(2*spin+1)))
	observable := hs.Observable{Dense: mat.Dense(*hs.ManyBodyOperator(hs.Sz(spin), 0, bc+1))}

	start := time.Now()
	for i := 0; i < bc; i += 1 {
		bath = append(bath, qs.State{Angle: float64(i) * math.Pi / float64(bc), Distance: 1e3})
	}

	s := &qs.System{
		CentralSpin:   cs,
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	b := conf.Physics.SpinEvolutionConfig.MagneticField
	b0 := 1.002 * b
	diagJob := qs.DiagonalizationInput{Hamiltonian: s.Hamiltonian(b0, b), B: b}
	diagOuts := make(chan qs.DiagonalizationResults)

	go s.Diagonalize(diagJob, diagOuts)

	diag := <-diagOuts
	close(diagOuts)

	var xys plotter.XYs
	for t := 0; t < timeRange; t += 1 {
		xys = append(xys, plotter.XY{X: float64(t), Y: observable.ExpectationValue(initialKet.Evolve(float64(t), diag.EigenValues, hs.KetsFromMatrix(diag.EigenVectors)))})
	}

	elapsed_time := time.Since(start)
	start_time := start.Format(time.RFC3339)

	au.PlotBasic(xys, "spin-"+start_time+".png", conf.Files)
	r := qs.ResultsIO{
		Filename: start_time,
		Metadata: qs.Metadata{
			Date:           start_time,
			Simulation:     "Central spin exp. val. evoluiton",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsed_time.String(),
		},
		Config: conf.Physics,
		XYs:    xys,
	}
	r.Write(conf.Files)
}
