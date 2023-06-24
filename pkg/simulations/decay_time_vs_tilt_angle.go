package simulations

import (
	"math"
	"time"

	cs "github.com/korsakjakub/cs_q_sim/pkg/cs_q_sim"
	"gonum.org/v1/plot/plotter"
)

func DecayTimeVsTiltAngle(conf cs.Config) {
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
			spread = 1e-8
		}
		decayTime := 1 / (spread * 1e-3)
		xys = append(xys, plotter.XY{X: tiltAngle / math.Pi, Y: decayTime})

		tiltAngle += conf.Physics.Dt
	}

	start_time := start.Format(time.RFC3339)

	elapsed_time := time.Since(start)
	r := cs.ResultsIO{
		Filename: start_time,
		Metadata: cs.Metadata{
			Date:           start_time,
			Simulation:     "Decay time vs tilt angle",
			SimulationId:   conf.Simulation,
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
