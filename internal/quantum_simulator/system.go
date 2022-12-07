package quantum_simulator

import (
	"math"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

const e0 = 8.854e-12

type System struct {
	CentralSpin   State
	Bath          []State
	PhysicsConfig PhysicsConfig
}

type State struct {
	Angle    float64
	Distance float64
}

func (s *System) forceAt(j int) float64 {
	if j == 0 {
		return 0.0
	}

	muMol, err := strconv.ParseFloat(s.PhysicsConfig.MoleculeMass, 64)
	if err != nil {
		parse(err)
	}
	muAtom, err := strconv.ParseFloat(s.PhysicsConfig.AtomMass, 64)
	if err != nil {
		parse(err)
	}

	c := (muMol * muAtom) / (4 * math.Pi * e0 * math.Pow(math.Abs(s.Bath[j-1].Distance), 3)) *
		0.5 * (1.0 - 3.0*math.Pow(math.Cos(s.Bath[j-1].Angle), 2))
	return c
}

func (s *System) hamiltonianHeisenbergTermAt(j int) *mat.Dense {
	bc, err := strconv.Atoi(s.PhysicsConfig.BathCount)
	if err != nil {
		parse(err)
	}
	spin, err := strconv.ParseFloat(s.PhysicsConfig.Spin, 64)
	if err != nil {
		parse(err)
	}

	dim := bc + 1
	f := s.forceAt(j)
	h := manyBody(Sp(spin), 0, dim)
	h.Mul(h, manyBody(Sm(spin), j, dim))
	h2 := manyBody(Sm(spin), 0, dim)
	h2.Mul(h2, manyBody(Sp(spin), j, dim))
	h.Add(h, h2)
	h.Apply(func(i, j int, v float64) float64 { return v * f }, h)

	return h
}

func (s *System) hamiltonianMagneticTerm(b0, b float64) *mat.Dense {
	bc, err := strconv.Atoi(s.PhysicsConfig.BathCount)
	if err != nil {
		parse(err)
	}
	spin, err := strconv.ParseFloat(s.PhysicsConfig.Spin, 64)
	if err != nil {
		parse(err)
	}

	h := manyBody(Sz(spin), 0, bc+1)
	h.Apply(func(i, j int, v float64) float64 { return v * (b0 - b) }, h)
	return h
}

func (s *System) hamiltonian(b0, b float64) *mat.Dense {
	spin, err := strconv.ParseFloat(s.PhysicsConfig.Spin, 64)
	if err != nil {
		parse(err)
	}
	spin_dim := int(2.0*spin + 1.0)
	bc, err := strconv.Atoi(s.PhysicsConfig.BathCount)
	if err != nil {
		parse(err)
	}
	dim := int(math.Pow(float64(spin_dim), float64(bc)+1.0))
	h := mat.NewDense(dim, dim, nil)
	for j := 0; j <= bc; j += 1 {
		h.Add(h, s.hamiltonianHeisenbergTermAt(j))
	}
	h.Add(h, s.hamiltonianMagneticTerm(b0, b))

	return h
}
