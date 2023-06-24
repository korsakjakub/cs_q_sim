package simulations

import (
	"fmt"
	"time"

	cs "github.com/korsakjakub/cs_q_sim/pkg/cs_q_sim"
)

/*
*
Given:
- List of Interactions
Restricted by:
- Class of geometries

Calculate:
- List of (R, \phi, \theta) with \phi being the azimutal angle - a degree of freedom
- Perhaps keep \beta constant
*
*/
func FindGeometryGivenInteractions(conf cs.Config) {
	var bath []cs.State
	conf.Physics.BathCount = len(conf.Physics.InteractionCoefficients) - 1
	bc := conf.Physics.BathCount

	start := time.Now()
	for i := 0; i < bc; i += 1 {
		bath = append(bath, cs.State{Angle: 0.0, Distance: 0.0})
	}

	s := &cs.System{
		CentralSpin:   cs.State{Angle: 0.0, Distance: 0.0},
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	if conf.Verbosity == "debug" {
		fmt.Println("Calculate interaction strength values...")
	}
	for j := 1; j <= bc; j += 1 {
		s.InteractionAt(j)
		s.DistanceGivenInteractionAt(j)
	}

	start_time := start.Format(time.RFC3339)

	elapsed_time := time.Since(start)
	r := cs.ResultsIO{
		Filename: start_time,
		Metadata: cs.Metadata{
			Date:           start_time,
			Simulation:     "Forces vs particle number",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsed_time.String(),
		},
		Values: struct {
			System cs.System "mapstructure:\"system\""
		}{
			System: *s,
		},
	}
	r.Write(conf.Files)
}
