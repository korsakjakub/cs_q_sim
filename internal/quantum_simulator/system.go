package quantum_simulator

import (
	"math"

	hs "github.com/korsakjakub/cs_q_sim/internal/hilbert_space"
	"gonum.org/v1/gonum/mat"
)

const e0 = 8.854e-12

type System struct {
	CentralSpin   State
	Bath          []State
	PhysicsConfig PhysicsConfig
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

func (s *System) DistanceGivenInteractionAt(j int) float64 {
	rj := math.Pow(math.Abs(s.PhysicsConfig.BathDipoleMoment*s.PhysicsConfig.AtomDipoleMoment/(4*math.Pi*e0)*0.5/s.Bath[j-1].InteractionStrength*(1.0-3.0*math.Pow(s.Bath[j-1].Angle, 2))), 1.0/3.0)
	s.Bath[j-1].Distance = rj
	return rj
}

func PolarAngleCos(j int, conf PhysicsConfig) float64 {
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
	} else if conf.Geometry == "icosahedron" {
		phi := 0.5 + math.Sqrt(5.0)*0.5
		a := math.Sqrt(phi*phi + 1.0)
		aphi := a * phi
		v := []point{{0.0, a, aphi}, {0.0, -a, aphi}, {0.0, a, -aphi}, {0.0, -a, -aphi},
			{a, aphi, 0.0}, {-a, aphi, 0.0}, {a, -aphi, 0.0}, {-a, -aphi, 0.0},
			{aphi, 0.0, a}, {-aphi, 0.0, a}, {aphi, 0.0, -a}, {-aphi, 0.0, -a}}
		return v[j].y*math.Sin(conf.TiltAngle) + v[j].z*math.Cos(conf.TiltAngle)
	} else if conf.Geometry == "sphere" { // https://stackoverflow.com/questions/9600801/evenly-distributing-n-points-on-a-sphere
		bc := float64(conf.BathCount)
		phi := math.Acos(1.0 - 2.0*(float64(j)+0.5)/float64(bc))
		theta := math.Pi * (1.0 + math.Sqrt(5.0)) * bc

		y := math.Sin(theta) * math.Sin(phi)
		z := math.Cos(phi)

		return y*math.Sin(conf.TiltAngle) + z*math.Cos(conf.TiltAngle)
	}
	return 0.0
}

// Given an index j, return force between the j-th bath molecule and the central spin
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
	sm := hs.Sm(spin)
	sp := hs.Sp(spin)

	f := s.InteractionAt(j)
	h := hs.ManyBodyOperator(sp, 0, dim)
	h.Mul(h, hs.ManyBodyOperator(sm, j, dim))
	h2 := hs.ManyBodyOperator(sm, 0, dim)
	h2.Mul(h2, hs.ManyBodyOperator(sp, j, dim))
	h.Add(h, h2)
	h.Scale(f, h)
	return mat.NewSymDense(h.RawMatrix().Cols, h.RawMatrix().Data)
}

// Given values of magnetic fields b0, and b, return the magnetic term of the hamiltonian
func (s *System) hamiltonianMagneticTerm(b0, b float64) *mat.SymDense {
	dim := len(s.Bath) + 1
	spin := s.PhysicsConfig.Spin
	sz := hs.Sz(spin)

	h := hs.ManyBodyOperator(sz, 0, dim)
	h.Scale(b0, h)

	for j := 1; j < dim; j++ {
		h2 := hs.ManyBodyOperator(sz, j, dim)
		h2.Scale(b, h2)
		h.Add(h, h2)
	}
	return mat.NewSymDense(h.RawMatrix().Cols, h.RawMatrix().Data)
}

// Given values of magnetic fields b0, and b, return the whole hamiltonian H_XX
func (s *System) Hamiltonian(b0, b float64) *mat.SymDense {
	spinDim := int(2.0*s.PhysicsConfig.Spin + 1.0)
	bc := len(s.Bath)
	dim := int(math.Pow(float64(spinDim), float64(bc)+1.0)) // bc is the BathCount and the total amount of objects in our system is BathCount + 1
	h := mat.NewSymDense(dim, nil)

	for j := 0; j <= bc; j += 1 {
		h.AddSym(h, s.hamiltonianHeisenbergTermAt(j))
	}
	h.AddSym(h, s.hamiltonianMagneticTerm(b0, b))

	return h
}

// Given a hamiltinian matrix, return its eigenvectors and eigenvalues
func (s *System) Diagonalize(hamiltonian *mat.SymDense) ([]float64, *mat.Dense) {
	var eig mat.EigenSym
	if err := eig.Factorize(hamiltonian, true); !err {
		panic("cannot diagonalize")
	}
	dim, _ := hamiltonian.Caps()
	eigenVectors := mat.NewDense(dim, dim, nil)
	eig.VectorsTo(eigenVectors)
	eigenValues := eig.Values(nil)

	// orto := NewOrtho(ComplexToFloats(eig.Values(nil)), hs.KetsFromCMatrix(evec))
	// orto.Orthonormalize()
	// return orto.OrthoToEigen()
	if !isDiagonalizedProperly(hamiltonian, eigenValues, eigenVectors) {
		panic("could not diagonalize properly")
	}
	return eigenValues, eigenVectors
}

func isDiagonalizedProperly(matrix *mat.SymDense, eig []float64, evec *mat.Dense) bool {
	for i := 0; i < evec.RawMatrix().Rows; i++ {
		left := mat.NewVecDense(len(eig), nil)
		vec := evec.ColView(i)
		left.MulVec(matrix, vec)
		right := mat.NewVecDense(len(eig), nil)
		right.ScaleVec(eig[i], vec)

		if !mat.EqualApprox(left, right, 1e-8) {
			return false
		}
	}
	return true
}

func RealPart(m *mat.CDense) *mat.Dense {
	r, c := m.Dims()
	out := mat.NewDense(r, c, nil)
	for i, el := range m.RawCMatrix().Data {
		out.RawMatrix().Data[i] = real(el)
	}
	return out
}

func ComplexToFloats(m []complex128) []float64 {
	out := make([]float64, len(m))
	for i, el := range m {
		out[i] = real(el)
	}
	return out
}
