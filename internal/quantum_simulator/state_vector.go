package quantum_simulator

import (
	"math/cmplx"

	"gonum.org/v1/gonum/blas/cblas128"
)

type StateVec cblas128.Vector

func NewKet(elements []complex128) *StateVec {
	return &StateVec{N: len(elements), Inc: 1, Data: elements}
}

func (u *StateVec) Dot(v *StateVec) complex128 {
	return cblas128.Dotc(cblas128.Vector(*u), cblas128.Vector(*v))
}

func (u *StateVec) Norm() float64 {
	return cblas128.Nrm2(cblas128.Vector(*u))
}

func (u *StateVec) Evolve(time float64, energies []float64, eigenBasis []StateVec) *StateVec {
	out := make([]complex128, len(eigenBasis))

	for i, basisVector := range eigenBasis {
		out[i] = cmplx.Exp(-complex(0, energies[i]*time)) * basisVector.Dot(u)
	}
	return NewKet(out)
}
