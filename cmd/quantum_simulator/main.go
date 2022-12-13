package main

import (
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"time"

	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	"gonum.org/v1/plot/plotter"
)

var conf qs.Config

func spectrum_vs_b(conf qs.Config) {
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
	jobs := make(chan qs.Input, fieldRange)
	results := make(chan qs.Results, fieldRange)

	go s.Diagonalize(jobs, results)

	for i := 0.0; i < float64(fieldRange); i += 1.0 {
		b := i * 1e3
		b0 := 1.0002 * b
		jobs <- qs.Input{Hamiltonian: s.Hamiltonian(b0, b), B: b}
	}

	close(jobs)

	for i := 0; i < fieldRange; i += 1 {
		v := <-results
		for _, ev := range v.EigenValues {
			xys = append(xys, plotter.XY{X: v.B, Y: cmplx.Abs(ev)})
		}
	}
	close(results)
	qs.Plot_spectrum_mag_field(xys, "first_plot.png", conf.Files)
	elapsed_time := time.Since(start)
	start_time := start.Format(time.RFC3339)

	r := qs.ResultsIO{
		Filename: start_time,
		Metadata: qs.Metadata{
			Date:           start_time,
			Simulation:     "spectrum vs. mag. field",
			Cpu:            "bsf",
			Ram:            "32GB",
			CompletionTime: elapsed_time.String(),
		},
		XYs: xys,
	}
	r.Write(conf.Files)
}

func main() {
	conf = qs.LoadConfig([]string{os.Getenv("CONFIG_PATH")}, os.Getenv("CONFIG_NAME"), os.Getenv("CONFIG_TYPE"))
	// conf = qs.LoadConfig([]string{"config"})
	spectrum_vs_b(conf)
}
