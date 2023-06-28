package cs_q_sim

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

const e0 = 8.854e-12

type System struct {
	CentralSpin   State
	Bath          []State
	PhysicsConfig PhysicsConfig
	DownSpins     int
}

type State struct {
	Angle               float64
	Distance            float64
	InteractionStrength float64
}

type point struct {
	x float64
	y float64
	z float64
}

type Eigen struct {
	EigenValues  []float64
	EigenVectors *mat.Dense
}

func (s *System) DistanceGivenInteractionAt(j int) float64 {
	rj := math.Pow(math.Abs(s.PhysicsConfig.BathDipoleMoment*s.PhysicsConfig.AtomDipoleMoment/(4*math.Pi*e0)*0.5/s.Bath[j-1].InteractionStrength*(1.0-3.0*math.Pow(s.Bath[j-1].Angle, 2))), 1.0/3.0)
	s.Bath[j-1].Distance = rj
	return rj
}

func PolarAngleCos(j int, conf PhysicsConfig) float64 {
	conf.TiltAngle *= math.Pi
	if conf.Geometry == "ring" {
		return math.Cos(float64(2*j)*math.Pi/float64(conf.BathCount)) * math.Sin(conf.TiltAngle)
	} else if conf.Geometry == "cube" && j < 8 {
		a := 1 / math.Sqrt(3.0)
		v := []point{{a, a, a}, {-a, a, a}, {a, -a, a}, {-a, -a, a}, {a, a, -a}, {-a, a, -a}, {a, -a, -a}, {-a, -a, -a}}
		return v[j].y*math.Sin(conf.TiltAngle) + v[j].z*math.Cos(conf.TiltAngle)
	} else if conf.Geometry == "dodecahedron" && j < 20 {
		a := 1 / math.Sqrt(3.0)
		phi := (0.5 + math.Sqrt(5.0)*0.5) * a
		iphi := 1.0 / phi * a
		v := []point{{a, a, a}, {-a, a, a}, {a, -a, a}, {-a, -a, a},
			{a, a, -a}, {-a, a, -a}, {a, -a, -a}, {-a, -a, -a},
			{0.0, phi, iphi}, {0.0, -phi, iphi}, {0.0, phi, -iphi}, {0.0, -phi, -iphi},
			{iphi, 0.0, phi}, {-iphi, 0.0, phi}, {iphi, 0.0, -phi}, {-iphi, 0.0, -phi},
			{phi, iphi, 0.0}, {-phi, iphi, 0.0}, {phi, -iphi, 0.0}, {-phi, -iphi, 0.0}}
		return v[j].y*math.Sin(conf.TiltAngle) + v[j].z*math.Cos(conf.TiltAngle)
	} else if conf.Geometry == "icosahedron" && j < 12 {
		/*
			vertices calculated with mathematica
			https://www.wolframcloud.com/obj/76badea4-ada5-4dc5-a415-8d6ea89de353
		*/
		v := []point{{0.0, 0.0, -1.0}, {0.0, 0.0, 1.0}, {-0.894427, 0.0, -0.447214}, {0.894427, 0.0, 0.447214}, {0.723607, -0.525731, -0.447214}, {0.723607, 0.525731, -0.447214}, {-0.723607, -0.525731, 0.447214}, {-0.723607, 0.525731, 0.447214}, {-0.276393, -0.850651, -0.447214}, {-0.276393, 0.850651, -0.447214}, {0.276393, -0.850651, 0.447214}, {0.276393, 0.850651, 0.447214}}
		return v[j].y*math.Sin(conf.TiltAngle) + v[j].z*math.Cos(conf.TiltAngle)
	} else if conf.Geometry == "sphere" { // https://stackoverflow.com/questions/9600801/evenly-distributing-n-points-on-a-sphere
		bc := float64(conf.BathCount)
		phi := math.Acos(1.0 - 2.0*(float64(j)+0.5)/bc)
		theta := math.Pi * (1.0 + math.Sqrt(5.0)) * bc

		y := math.Sin(theta) * math.Sin(phi)
		z := math.Cos(phi)

		return y*math.Sin(conf.TiltAngle) + z*math.Cos(conf.TiltAngle)
	}
	return 0.0
}

// InteractionAt returns the interaction strength between the j-th bath molecule and the central spin, given an index j
func (s *System) InteractionAt(j int) float64 {
	if len(s.PhysicsConfig.InteractionCoefficients) > 0 {
		cj := s.PhysicsConfig.InteractionCoefficients[j]
		s.Bath[j-1].InteractionStrength = cj
		return cj
	} else {
		if j == 0 {
			return 0.0
		}

		var k float64
		if s.PhysicsConfig.Units == "atomic" {
			k = 149.42785955012954
		} else { // SI
			k = 1 / (4 * math.Pi * e0)
		}

		// Bath has indices 0:BathCount-1, and j has a range of 0:BathCount -> for j = 0 we mean the central spin which is not a part of the Bath.
		// Therefore we pick Bath[j-1] instead of Bath[j]
		c := k * (s.PhysicsConfig.BathDipoleMoment * s.PhysicsConfig.AtomDipoleMoment) / math.Pow(math.Abs(s.Bath[j-1].Distance), 3) *
			0.5 * (1.0 - 3.0*math.Pow(s.Bath[j-1].Angle, 2))

		// assign force value to bath state
		s.Bath[j-1].InteractionStrength = c
		return c
	}
}

// Given an index j, return the Heisenberg term (0, j - interaction) of the hamiltonian
func (s *System) hamiltonianHeisenbergTermAt(j int) *mat.SymDense {
	spin := s.PhysicsConfig.Spin
	dim := len(s.Bath) + 1 // bc is the BathCount and the total amount of objects in our system is BathCount + 1
	sm := Sm(spin)
	sp := Sp(spin)

	f := s.InteractionAt(j)
	h := ManyBodyOperator(sp, 0, dim)
	h.Mul(h, ManyBodyOperator(sm, j, dim))
	h2 := ManyBodyOperator(sm, 0, dim)
	h2.Mul(h2, ManyBodyOperator(sp, j, dim))
	h.Add(h, h2)
	h.Scale(f, h)
	return mat.NewSymDense(h.RawMatrix().Cols, h.RawMatrix().Data)
}

// Given values of magnetic fields b0, and b, return the magnetic term of the hamiltonian
func (s *System) hamiltonianMagneticTerm(b0, b float64) *mat.SymDense {
	dim := len(s.Bath) + 1
	spin := s.PhysicsConfig.Spin
	sz := Sz(spin)

	h := ManyBodyOperator(sz, 0, dim)
	h.Scale(b0, h)

	for j := 1; j < dim; j++ {
		h2 := ManyBodyOperator(sz, j, dim)
		h2.Scale(b, h2)
		h.Add(h, h2)
	}
	return mat.NewSymDense(h.RawMatrix().Cols, h.RawMatrix().Data)
}

// Hamiltonian returns the whole H_XX given values of magnetic fields b0, and b
func (s *System) Hamiltonian(b0, b float64) *mat.SymDense {
	spinDim := int(2.0*s.PhysicsConfig.Spin + 1.0)
	bc := len(s.Bath)
	dim := int(math.Pow(float64(spinDim), float64(bc)+1.0)) // bc is the BathCount and the total amount of objects in our system is BathCount + 1
	h := mat.NewSymDense(dim, nil)

	for j := 0; j <= bc; j += 1 {
		h.AddSym(h, s.hamiltonianHeisenbergTermAt(j))
	}
	h.AddSym(h, s.hamiltonianMagneticTerm(b0, b))

	if !mat.EqualApprox(h, h.T(), 1e-8) {
		panic("Hamiltonian is not symmetric.")
	}

	return h
}

func (s *System) HamiltonianInBase(b0, b float64, indices []int) *mat.SymDense {
	if len(indices) < 2 {
		panic("Only 1 dimension remained")
	}
	fullHamiltonian := s.Hamiltonian(b0, b)
	hDim := fullHamiltonian.SymmetricDim()
	dim := len(indices)

	dFullHamiltonian := mat.NewDense(hDim, hDim, fullHamiltonian.RawSymmetric().Data)

	h := RestrictMatrixToSubspace(dFullHamiltonian, indices)

	return mat.NewSymDense(dim, h.RawMatrix().Data)
}

// Diagonalize returns eigenvectors and eigenvalues given a hamiltonian matrix
func (s *System) Diagonalize(hamiltonian *mat.SymDense) Eigen {
	var eig mat.EigenSym
	if err := eig.Factorize(hamiltonian, true); !err {
		panic("cannot diagonalize")
	}
	dim, _ := hamiltonian.Caps()
	eigenVectors := mat.NewDense(dim, dim, nil)
	eig.VectorsTo(eigenVectors)
	eigenValues := eig.Values(nil)

	isDiagonalizedProperly := func(matrix *mat.SymDense, eig []float64, evec *mat.Dense) error {
		for i := 0; i < evec.RawMatrix().Rows; i++ {
			left := mat.NewVecDense(len(eig), nil)
			vec := evec.ColView(i)
			left.MulVec(matrix, vec)
			right := mat.NewVecDense(len(eig), nil)
			right.ScaleVec(eig[i], vec)

			if !mat.EqualApprox(left, right, 1e-8) {
				dVec := mat.NewVecDense(len(eig), nil)
				dVec.SubVec(left, right)
				return fmt.Errorf("the restriction A v = k v is not satisfied\nd = %f", dVec.Norm(2))
			}
		}
		if math.Abs(math.Abs(mat.Det(evec))-1.0) > 1e-8 {
			return fmt.Errorf("the determinant of Eigenvector matrix isn't +-1. Det = %v", mat.Det(evec))
		}
		return nil
	}

	if err := isDiagonalizedProperly(hamiltonian, eigenValues, eigenVectors); err != nil {
		fmt.Println(err)
	}
	return Eigen{eigenValues, eigenVectors}
}
