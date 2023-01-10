package main

import (
	"math"
	"math/cmplx"
	"sync"
	"time"

	au "github.com/korsakjakub/cs_q_sim/internal/analysis_utilities"
	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	"gonum.org/v1/plot/plotter"
)

func spin_time_evolution(conf qs.Config) {
	/*
		Given:
		- initial condition -> psi(t=0)
		- hamiltonian <- system

		Do:
		- calculate \ket{\psi(t)} = \sum_j e^(-i E_j t) \ket{E_j}\bra{E_j} \ket{\psi(0)}
		- plot ||psi(t)||
	*/
	cs := qs.State{Angle: 0.0, Distance: 0.0}
	var bath []qs.State
	bc := conf.Physics.BathCount
	timeRange := conf.Physics.SpinEvolutionConfig.TimeRange
	start := time.Now()
	for i := 0; i < bc; i += 1 {
		bath = append(bath, qs.State{Angle: float64(i) * math.Pi / float64(bc), Distance: 1e3})
	}

	s := &qs.System{
		CentralSpin:   cs,
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	var xys plotter.XYs
	results := make(chan qs.Results, timeRange)
	var jobs []qs.Input

	for i := 0.0; i < float64(timeRange); i += 1.0 {
		b := conf.Physics.SpinEvolutionConfig.MagneticField
		b0 := 1.0002 * b
		jobs = append(jobs, qs.Input{Hamiltonian: s.Hamiltonian(b0, b), B: b})
	}

	var wg sync.WaitGroup
	wg.Add(len(jobs))

	for _, job := range jobs {
		go func(j qs.Input) {
			defer wg.Done()
			s.Diagonalize(j, results)
		}(job)
	}

	wg.Wait()

	for i := 0; i < timeRange; i += 1 {
		v := <-results
		for _, ev := range v.EigenValues {
			xys = append(xys, plotter.XY{X: v.B, Y: cmplx.Abs(ev)})
		}
	}
	close(results)
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
