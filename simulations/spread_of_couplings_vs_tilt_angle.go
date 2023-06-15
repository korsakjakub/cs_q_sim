package simulations

import (
	"math"
	"time"

	cs "github.com/korsakjakub/cs_q_sim/pkg/cs_q_sim"
	"gonum.org/v1/plot/plotter"
)

func spread(states []cs.State) float64 {
	max := math.Abs(states[0].InteractionStrength)
	min := math.Abs(states[0].InteractionStrength)
	for i := range states {
		strength := math.Abs(states[i].InteractionStrength)
		if strength > max {
			max = strength
		} else if strength < min {
			min = strength
		}
	}
	return max - min
}

func prepareStates(conf cs.Config) []cs.State {
	var bath []cs.State
	bc := conf.Physics.BathCount

	for i := 0; i < bc; i += 1 {
		bath = append(bath, cs.State{Angle: cs.PolarAngleCos(i, conf.Physics), Distance: conf.Physics.ConstantDistance})
	}

	s := &cs.System{
		CentralSpin:   cs.State{Angle: 0.0, Distance: 0.0},
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	for j := 0; j <= bc; j += 1 {
		s.InteractionAt(j)
	}
	return s.Bath
}

func SpreadOfCouplingsVsTiltAngle(conf cs.Config) {
	start := time.Now()

	if len(conf.Physics.TiltAngleRange) != 2 {
		panic("TiltAngleRange should have length 2. (min, max)")
	}
	minTiltAngle := conf.Physics.TiltAngleRange[0] * math.Pi
	maxTiltAngle := conf.Physics.TiltAngleRange[1] * math.Pi

	tiltAngle := minTiltAngle
	var xys plotter.XYs

	for tiltAngle < maxTiltAngle {
		conf.Physics.TiltAngle = tiltAngle
		states := prepareStates(conf)
		spread := spread(states)
		if spread < 1e-8 {
			spread = 0.0
		}
		xys = append(xys, plotter.XY{X: tiltAngle / math.Pi, Y: spread * 1e-3})

		tiltAngle += conf.Physics.Dt
	}

	start_time := start.Format(time.RFC3339)

	elapsed_time := time.Since(start)
	r := cs.ResultsIO{
		Filename: start_time,
		Metadata: cs.Metadata{
			Date:           start_time,
			Simulation:     "Spread of couplings vs tilt angle",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsed_time.String(),
		},
		Values: struct {
			System cs.System "mapstructure:\"system\""
		}{
			cs.System{PhysicsConfig: conf.Physics},
		},
		XYs: []plotter.XYs{xys},
	}
	r.Write(conf.Files)
}
