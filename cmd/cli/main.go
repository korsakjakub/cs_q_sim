package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	cs "github.com/korsakjakub/cs_q_sim/pkg/cs_q_sim"
	sim "github.com/korsakjakub/cs_q_sim/pkg/simulations"
	"github.com/spf13/pflag"
)

var conf cs.Config

func printHeader(name string) {
	fmt.Printf("Starting the simulation: %s...\n", name)
}

func main() {
	var configFiles []string
	pflag.StringSliceVarP(&configFiles, "values", "f", []string{}, "specify values in a YAML file or a URL (can specify multiple)")
	pflag.Parse()
	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			panic(err)
		}
	}
	fmt.Println(configFiles)
	conf = cs.LoadConfig(configFiles)
	if conf.Verbosity == "debug" {
		fmt.Println("Config:")
		spew.Dump(conf)
	}

	if err := cs.Validate(conf.Files, []string{
		"FigDir",
		"OutputsDir",
		"ResultsConfig",
	}); err != nil {
		panic(err)
	}

	switch conf.Simulation {
	case "spin-evolution":
		printHeader("spin evolution")
		if err := cs.Validate(conf.Physics, []string{
			"BathDipoleMoment",
			"AtomDipoleMoment",
			"Spin",
			"TiltAngle",
			"ConstantDistance",
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
	case "spread-of-couplings":
		printHeader("spread of couplings")
		if err := cs.Validate(conf.Physics, []string{
			"BathDipoleMoment",
			"AtomDipoleMoment",
			"BathCount",
			"Spin",
			"TiltAngleRange",
			"ConstantDistance",
			"Geometry",
			"BathMagneticField",
			"CentralMagneticField",
			"Dt",
		}); err != nil {
			panic(err)
		}
		sim.SpreadOfCouplingsVsTiltAngle(conf)
	case "spread-of-couplings-selected-coeffs":
		if err := cs.Validate(conf.Physics, []string{
			"Spin",
			"InteractionCoefficients",
			"TiltAngleRange",
			"Dt",
		}); err != nil {
			panic(err)
		}
		printHeader("spread of couplings for selected coefficients")
		sim.SpreadOfCouplingsVsTiltAngle(conf)
	case "decay-time":
		printHeader("decay time")
		if err := cs.Validate(conf.Physics, []string{
			"BathDipoleMoment",
			"AtomDipoleMoment",
			"BathCount",
			"Spin",
			"TiltAngleRange",
			"ConstantDistance",
			"Geometry",
			"BathMagneticField",
			"CentralMagneticField",
			"Dt",
		}); err != nil {
			panic(err)
		}
		sim.DecayTimeVsTiltAngle(conf)
	case "spectrum":
		printHeader("spectrum")
		sim.Spectrum(conf)
	case "interactions":
		if err := cs.Validate(conf.Physics, []string{
			"BathDipoleMoment",
			"AtomDipoleMoment",
			"BathCount",
			"ConstantDistance",
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
	case "spin-evolution-selected-coeffs":
		if err := cs.Validate(conf.Physics, []string{
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
		sim.SpinTimeEvolution(conf)
	case "find-geometry-given-interactions":
		s := []string{
			"BathDipoleMoment",
			"AtomDipoleMoment",
			"Spin",
			"InteractionCoefficients",
			"BathMagneticField",
			"CentralMagneticField",
		}
		if err := cs.Validate(conf.Physics, s); err != nil {
			panic(err)
		}
		printHeader("Find geometry for specified interactions")
		sim.FindGeometryGivenInteractions(conf)
	}
}
