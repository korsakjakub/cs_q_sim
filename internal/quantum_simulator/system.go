package quantum_simulator

import (
	"math"
	"strconv"
)

const e0 = 8.854e-12

type System struct {
	CentralSpin State
	Bath        []State
}

type State struct {
	Angle            float64
	Distance         float64
	MagneticField    float64
	InteractionForce float64
}

func (s *System) forces(pc PhysicsConfig) []float64 {

	bc, err := strconv.Atoi(pc.BathCount)

	if err != nil {
		parse(err)
	}
	muMol, err := strconv.ParseFloat(pc.MoleculeMass, 64)
	if err != nil {
		parse(err)
	}
	muAtom, err := strconv.ParseFloat(pc.AtomMass, 64)
	if err != nil {
		parse(err)
	}

	var C = make([]float64, bc)
	for j := 0; j < bc; j += 1 {
		C[j] = (muMol * muAtom) / (4 * math.Pi * e0 * math.Pow(math.Abs(s.Bath[j].Distance), 3)) *
			0.5 * (1.0 - 3.0*math.Pow(math.Cos(s.Bath[j].Angle), 2))
	}
	return C
}
