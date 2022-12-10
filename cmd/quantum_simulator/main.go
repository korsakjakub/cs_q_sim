package main

import (
	"math"
	"math/cmplx"
	"os"
	"strconv"

	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	"gonum.org/v1/gonum/mat"
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
	for i := 0; i < bc; i += 1 {
		bath = append(bath, qs.State{Angle: float64(i) * math.Pi / float64(bc), Distance: 1e3})
	}
	s := &qs.System{
		CentralSpin:   cs,
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	var xys plotter.XYs
	for b := 0.0; b < 1e6; b += 1e4 {
		b0 := 1.0002 * b
		eigenVectors := make(chan *mat.CDense, 100)
		eigenValues := make(chan complex128, 100)

		go s.Diagonalize(s.Hamiltonian(b0, b), eigenVectors, eigenValues)

		var eval []complex128
		for {
			v, ok := <-eigenValues
			if !ok {
				break
			}
			eval = append(eval, v)
			xys = append(xys, plotter.XY{X: b, Y: cmplx.Abs(v)})
		}
	}
	qs.Plot_spectrum_mag_field(xys, "first_plot.png", conf.Files)
}

func main() {
	conf = qs.LoadConfig([]string{os.Getenv("CONFIG_PATH")}, os.Getenv("CONFIG_NAME"), os.Getenv("CONFIG_TYPE"))
	// conf = qs.LoadConfig([]string{"../../config"})
	spectrum_vs_b(conf)
}
