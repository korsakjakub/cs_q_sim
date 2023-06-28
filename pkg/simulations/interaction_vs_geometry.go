package simulations

import (
	"time"

	cs "github.com/korsakjakub/cs_q_sim/pkg/cs_q_sim"
	"gonum.org/v1/plot/plotter"
)

func Interactions(conf cs.Config) {
	var bath []cs.State
	bc := conf.Physics.BathCount

	start := time.Now()
	for i := 0; i < bc; i += 1 {
		bath = append(bath, cs.State{Angle: cs.PolarAngleCos(i, conf.Physics), Distance: conf.Physics.ConstantDistance})
	}

	start_time := start.Format(time.RFC3339)
	var xys plotter.XYs

	s := &cs.System{
		CentralSpin:   cs.State{Angle: 0.0, Distance: 0.0},
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	for j := 0; j <= bc; j += 1 {
		s.InteractionAt(j)
	}
	for j := 0; j < bc; j += 1 {
		xys = append(xys, plotter.XY{X: float64(j), Y: s.Bath[j].InteractionStrength * 1e-3})
	}

	elapsed_time := time.Since(start)
	r := cs.ResultsIO{
		Filename: start_time,
		Metadata: cs.Metadata{
			Date:           start_time,
			Simulation:     "Interaction strength",
			SimulationId:   "interactions",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsed_time.String(),
		},
		Values: struct {
			System cs.System "mapstructure:\"system\""
		}{
			System: *s,
		},
		XYs: []plotter.XYs{xys},
	}
	r.Write(conf.Files)
}
