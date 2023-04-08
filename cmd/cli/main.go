package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	sim "github.com/korsakjakub/cs_q_sim/simulations"
	pflag "github.com/spf13/pflag"
)

var conf qs.Config

func printHeader(name string) {
	fmt.Printf("Starting the simulation: %s...\n", name)
}

func main() {
	var configFiles []string
	pflag.StringSliceVarP(&configFiles, "values", "f", []string{}, "specify values in a YAML file or a URL (can specify multiple)")
	pflag.Parse()
	conf = qs.LoadConfig(configFiles)
	if conf.Verbosity == "debug" {
		fmt.Println("Config:")
		spew.Dump(conf)
	}

	if err := qs.Validate(conf.Files, []string{
		"FigDir",
		"OutputsDir",
		"ResultsConfig",
	}); err != nil {
		panic(err)
	}

	switch conf.Simulation {
	case "spin-evolution":
		printHeader("spin evolution")
		if err := qs.Validate(conf.Physics, []string{
				"BathDipoleMoment", 
				"AtomDipoleMoment",
				"Spin",
				"TiltAngle",
				"Geometry",
				"BathMagneticField",
				"CentralMagneticField",
				"TimeRange",
				"Dt",
				"InitialKet",
				"ObservablesConfig",
			}); err != nil {
			panic(err)
		}
		sim.SpinTimeEvolution(conf)
	case "spectrum":
		printHeader("spectrum")
		sim.Spectrum(conf)
	case "interactions":
		if err := qs.Validate(conf.Physics, []string{
				"BathDipoleMoment", 
				"AtomDipoleMoment",
				"Spin",
				"TiltAngle",
				"Geometry",
				"BathMagneticField",
				"CentralMagneticField",
			}); err != nil {
			panic(err)
		}
		printHeader("interactions")
		sim.Interactions(conf)
	case "time-evolution-selected-coeffs":
		if err := qs.Validate(conf.Physics, []string{
				"Spin",
				"InteractionCoefficients",
				"BathMagneticField",
				"CentralMagneticField",
				"TimeRange",
				"Dt",
				"InitialKet",
				"ObservablesConfig",
			}); err != nil {
			panic(err)
		}
		printHeader("spin evolution for selected coefficients")
		sim.SpinTimeEvolutionSelectedCoeffs(conf)
	}
}

