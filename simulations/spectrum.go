package simulations

import (
	"math"
	"sync"
	"time"

	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/plotter"
)

type spectrumInput struct {
	hamiltonian   *mat.SymDense
	magneticField float64
}

type spectrumOutput struct {
	eigenValues   []float64
	eigenVectors  *mat.Dense
	magneticField float64
}

func Spectrum(conf qs.Config) {
	cs := qs.State{Angle: 0.0, Distance: 0.0}
	var bath []qs.State
	bc := conf.Physics.BathCount
	fieldRange := conf.Physics.MagneticFieldRange
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

	r := qs.ResultsIO{
		Filename: start_time,
		Metadata: qs.Metadata{
			Date:           start_time,
			Simulation:     "burst spectrum vs. mag. field",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsed_time.String(),
		},
		Values: struct {
			System       qs.System "mapstructure:\"system\""
			EigenValues  []string  "mapstructure:\"evalues\""
			EigenVectors []string  "mapstructure:\"evectors\""
		}{System: *s},
		XYs: []plotter.XYs{xys},
	}
	r.Write(conf.Files)
}
