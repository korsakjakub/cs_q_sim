package simulations

import (
	"fmt"
	"time"

	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
)

/**
Given:
- List of Interactions
Restricted by:
- Class of geometries

Calculate:
- List of (R, \phi, \theta) with \phi being the azimutal angle - a degree of freedom
- Perhaps keep \beta constant
**/
func FindGeometryGivenInteractions(conf qs.Config) {
	cs := qs.State{Angle: 0.0, Distance: 0.0}
	var bath []qs.State
	conf.Physics.BathCount = len(conf.Physics.InitialKet) - 1
	bc := conf.Physics.BathCount

	start := time.Now()
	for i := 0; i < bc; i += 1 {
		bath = append(bath, qs.State{Angle: qs.PolarAngleCos(i, conf.Physics), Distance: 1e3})
	}

	s := &qs.System{
		CentralSpin:   cs,
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	if conf.Verbosity == "debug" {
		fmt.Println("Calculate interaction strength values...")
	}
	for j := 0; j <= bc; j += 1 {
		s.InteractionAt(j)
		s.DistanceGivenInteractionAt(j)
	}

	start_time := start.Format(time.RFC3339)

	elapsed_time := time.Since(start)
	r := qs.ResultsIO{
		Filename: start_time,
		Metadata: qs.Metadata{
			Date:           start_time,
			Simulation:     "Forces vs particle number",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsed_time.String(),
		},
		Values: struct {
			System       qs.System "mapstructure:\"system\""
			EigenValues  []string  "mapstructure:\"evalues\""
			EigenVectors []string  "mapstructure:\"evectors\""
		}{
			System: *s,
		},
	}
	r.Write(conf.Files)
}
