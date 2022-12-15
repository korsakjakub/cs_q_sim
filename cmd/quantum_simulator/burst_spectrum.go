package main

import (
	"math"
	"math/cmplx"
	"strconv"
	"sync"
	"time"

	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	"gonum.org/v1/plot/plotter"
)

func burstSpectrum(conf qs.Config) {
	cs := qs.State{Angle: 0.0, Distance: 0.0}
	var bath []qs.State
	bc, err := strconv.Atoi(conf.Physics.BathCount)
	if err != nil {
		panic(err)
	}
	fieldRange, err := strconv.Atoi(conf.Physics.FieldRange)
	if err != nil {
		panic(err)
	}
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
	results := make(chan qs.Results, fieldRange)
	var jobs []qs.Input

	for i := 0.0; i < float64(fieldRange); i += 1.0 {
		b := i * 1e3
		b0 := 1.0002 * b
		jobs = append(jobs, qs.Input{Hamiltonian: s.Hamiltonian(b0, b), B: b})
	}

	var wg sync.WaitGroup
	wg.Add(len(jobs))

	for _, job := range jobs {
		go func(j qs.Input) {
			defer wg.Done()
			s.DiagonalizeBurst(j, results)
		}(job)
	}

	wg.Wait()

	for i := 0; i < fieldRange; i += 1 {
		v := <-results
		for _, ev := range v.EigenValues {
			xys = append(xys, plotter.XY{X: v.B, Y: cmplx.Abs(ev)})
		}
	}
	close(results)
	elapsed_time := time.Since(start)
	start_time := start.Format(time.RFC3339)

	qs.Plot_spectrum_mag_field(xys, start_time+".png", conf.Files)
	r := qs.ResultsIO{
		Filename: start_time,
		Metadata: qs.Metadata{
			Date:           start_time,
			Simulation:     "burst spectrum vs. mag. field",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsed_time.String(),
		},
		Config: conf.Physics,
		XYs:    xys,
	}
	r.Write(conf.Files)
}
