package simulations

import (
	"math"
	"sync"
	"time"

	cs "github.com/korsakjakub/cs_q_sim/pkg/cs_q_sim"
	"gonum.org/v1/plot/plotter"
)

func Spectrum(conf cs.Config) {
	var bath []cs.State
	bc := conf.Physics.BathCount
	fieldRange := conf.Physics.MagneticFieldRange
	start := time.Now()
	for i := 0; i < bc; i += 1 {
		bath = append(bath, cs.State{Angle: float64(i) * math.Pi / float64(bc), Distance: 1e3})
	}

	s := &cs.System{
		CentralSpin:   cs.State{Angle: 0.0, Distance: 0.0},
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	var xys plotter.XYs
	results := make(chan spectrumOutput, fieldRange)
	var jobs []spectrumInput

	for i := 0.0; i < float64(fieldRange); i += 1.0 {
		b := i * 1e3
		b0 := 1.0002 * b
		jobs = append(jobs, spectrumInput{s.Hamiltonian(b0, b), b})
	}

	var wg sync.WaitGroup
	wg.Add(len(jobs))

	for _, job := range jobs {
		go func(j spectrumInput) {
			defer wg.Done()
			eigen := s.Diagonalize(j.hamiltonian)
			results <- spectrumOutput{eigenValues: eigen.EigenValues, eigenVectors: eigen.EigenVectors, magneticField: j.magneticField}
		}(job)
	}

	wg.Wait()

	for i := 0; i < fieldRange; i += 1 {
		v := <-results
		for _, ev := range v.eigenValues {
			xys = append(xys, plotter.XY{X: v.magneticField, Y: ev})
		}
	}
	close(results)
	elapsed_time := time.Since(start)
	start_time := start.Format(time.RFC3339)

	r := cs.ResultsIO{
		Filename: start_time,
		Metadata: cs.Metadata{
			Date:           start_time,
			Simulation:     "burst spectrum vs. mag. field",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsed_time.String(),
		},
		Values: struct {
			System cs.System "mapstructure:\"system\""
		}{System: *s},
		XYs: []plotter.XYs{xys},
	}
	r.Write(conf.Files)
}
