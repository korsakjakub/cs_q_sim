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

// Given an index j, return force between the j-th bath molecule and the central spin
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

	// Bath has indices 0:BathCount-1, and j has a range of 0:BathCount -> for j = 0 we mean the central spin which is not a part of the Bath.
	// Therefore we pick Bath[j-1] instead of Bath[j]
	c := (muMol * muAtom) / (4 * math.Pi * e0 * math.Pow(math.Abs(s.Bath[j-1].Distance), 3)) *
		0.5 * (1.0 - 3.0*math.Pow(math.Cos(s.Bath[j-1].Angle), 2))
	return c
}

// Given an index j, return the Heisenberg term (0, j - interaction) of the hamiltonian
func (s *System) hamiltonianHeisenbergTermAt(j int) *mat.Dense {
	bc, err := strconv.Atoi(s.PhysicsConfig.BathCount)
	if err != nil {
		parse(err)
	}
	spin, err := strconv.ParseFloat(s.PhysicsConfig.Spin, 64)
	if err != nil {
		parse(err)
	}

	dim := bc + 1 // bc is the BathCount and the total amount of objects in our system is BathCount + 1
	f := s.forceAt(j)
	h := manyBody(Sp(spin), 0, dim)
	h.Mul(h, manyBody(Sm(spin), j, dim))
	h2 := manyBody(Sm(spin), 0, dim)
	h2.Mul(h2, manyBody(Sp(spin), j, dim))
	h.Add(h, h2)
	h.Apply(func(i, j int, v float64) float64 { return v * f }, h) // Multiplication by a scalar f

	return h
}

// Given values of magnetic fields b0, and b, return the magnetic term of the hamiltonian
func (s *System) hamiltonianMagneticTerm(b0, b float64) *mat.Dense {
	bc, err := strconv.Atoi(s.PhysicsConfig.BathCount)
	if err != nil {
		parse(err)
	}
	spin, err := strconv.ParseFloat(s.PhysicsConfig.Spin, 64)
	if err != nil {
		parse(err)
	}

	h := manyBody(Sz(spin), 0, bc+1) // bc is the BathCount and the total amount of objects in our system is BathCount + 1
	h.Apply(func(i, j int, v float64) float64 { return v * (b0 - b) }, h)
	return h
}

// Given values of magnetic fields b0, and b, return the whole hamiltonian H_XX
func (s *System) Hamiltonian(b0, b float64) *mat.Dense {
	spin, err := strconv.ParseFloat(s.PhysicsConfig.Spin, 64)
	if err != nil {
		parse(err)
	}
	spin_dim := int(2.0*spin + 1.0)
	bc, err := strconv.Atoi(s.PhysicsConfig.BathCount)
	if err != nil {
		parse(err)
	}
	dim := int(math.Pow(float64(spin_dim), float64(bc)+1.0)) // bc is the BathCount and the total amount of objects in our system is BathCount + 1
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
func (s *System) Diagonalize(input <-chan Input, results chan<- Results) {
	var eig mat.Eigen
	for n := range input {
		if err := eig.Factorize(n.Hamiltonian, mat.EigenRight); !err {
			panic("cannot diagonalize")
		}
		dim, _ := n.Hamiltonian.Caps()
		evec := mat.NewCDense(dim, dim, nil)
		eig.VectorsTo(evec)

		results <- Results{
			EigenVectors: evec,
			EigenValues:  eig.Values(nil),
			B:            n.B,
		}
	}
}

// Given a hamiltinian matrix, return its eigenvectors and eigenvalues
func (s *System) DiagonalizeBurst(input Input, results chan<- Results) {
	n := input
	var eig mat.Eigen
	if err := eig.Factorize(n.Hamiltonian, mat.EigenRight); !err {
		panic("cannot diagonalize")
	}
	dim, _ := n.Hamiltonian.Caps()
	evec := mat.NewCDense(dim, dim, nil)
	eig.VectorsTo(evec)

	results <- Results{
		EigenVectors: evec,
		EigenValues:  eig.Values(nil),
		B:            n.B,
	}
}
