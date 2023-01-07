package quantum_simulator

import (
	"math"

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

// Given an index j, return force between the j-th bath molecule and the central spin
func (s *System) forceAt(j int) float64 {
	if j == 0 {
		return 0.0
	}

	// Bath has indices 0:BathCount-1, and j has a range of 0:BathCount -> for j = 0 we mean the central spin which is not a part of the Bath.
	// Therefore we pick Bath[j-1] instead of Bath[j]
	c := (s.PhysicsConfig.MoleculeMass * s.PhysicsConfig.AtomMass) / (4 * math.Pi * e0 * math.Pow(math.Abs(s.Bath[j-1].Distance), 3)) *
		0.5 * (1.0 - 3.0*math.Pow(math.Cos(s.Bath[j-1].Angle), 2))
	return c
}

// Given an index j, return the Heisenberg term (0, j - interaction) of the hamiltonian
func (s *System) hamiltonianHeisenbergTermAt(j int) *mat.Dense {
	spin := s.PhysicsConfig.Spin
	dim := s.PhysicsConfig.BathCount + 1 // bc is the BathCount and the total amount of objects in our system is BathCount + 1
	sm := Sm(spin)
	sp := Sp(spin)

	f := s.forceAt(j)
	h := manyBody(sp, 0, dim)
	h.Mul(h, manyBody(sm, j, dim))
	h2 := manyBody(sm, 0, dim)
	h2.Mul(h2, manyBody(sp, j, dim))
	h.Add(h, h2)
	h.Scale(f, h)
	return h
}

// Given values of magnetic fields b0, and b, return the magnetic term of the hamiltonian
func (s *System) hamiltonianMagneticTerm(b0, b float64) *mat.Dense {
	var h mat.Dense
	h.Scale(b0-b, manyBody(Sz(s.PhysicsConfig.Spin), 0, s.PhysicsConfig.BathCount+1))
	return &h
}

// Given values of magnetic fields b0, and b, return the whole hamiltonian H_XX
func (s *System) Hamiltonian(b0, b float64) *mat.Dense {
	spinDim := int(2.0*s.PhysicsConfig.Spin + 1.0)
	bc := s.PhysicsConfig.BathCount
	dim := int(math.Pow(float64(spinDim), float64(bc)+1.0)) // bc is the BathCount and the total amount of objects in our system is BathCount + 1
	h := mat.NewDense(dim, dim, nil)

	for j := 0; j <= bc; j += 1 {
		h.Add(h, s.hamiltonianHeisenbergTermAt(j))
	}
	h.Add(h, s.hamiltonianMagneticTerm(b0, b))

	return h
}

type Results struct {
	EigenVectors *mat.CDense
	EigenValues  []complex128
	B            float64
}

type Input struct {
	Hamiltonian *mat.Dense
	B           float64
}

// Given a hamiltinian matrix, return its eigenvectors and eigenvalues
func (s *System) Diagonalize(input Input, results chan<- Results) {
	var eig mat.Eigen
	if err := eig.Factorize(input.Hamiltonian, mat.EigenRight); !err {
		panic("cannot diagonalize")
	}
	dim, _ := input.Hamiltonian.Caps()
	evec := mat.NewCDense(dim, dim, nil)
	eig.VectorsTo(evec)

	results <- Results{
		EigenVectors: evec,
		EigenValues:  eig.Values(nil),
		B:            input.B,
	}
}
