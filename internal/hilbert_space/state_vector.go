package hilbert_space

import (
	"math/cmplx"

	"gonum.org/v1/gonum/blas/cblas128"
	"gonum.org/v1/gonum/mat"
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

func (u *StateVec) Evolve(time float64, energies []complex128, eigenBasis []*StateVec) *StateVec {
	out := make([]complex128, len(eigenBasis))

	for i, basisVector := range eigenBasis {
		out[i] = cmplx.Exp(-energies[i]*complex(0, time)) * basisVector.Dot(u)
	}
	return NewKet(out)
}

func KetFromFloats(elements []float64) *StateVec {
	if len(elements)%2 != 0 {
		panic("Odd size of input vector")
	}
	cElements := make([]complex128, len(elements)/2)
	for i := range cElements {
		cElements[i] = complex(elements[2*i], elements[2*i+1])
	}
	return NewKet(cElements)
}

func KetsFromMatrix(mat mat.CMatrix) []*StateVec {
	rows, cols := mat.Dims()
	out := make([]*StateVec, rows)
	for row := 0; row < rows; row++ {
		tmp := make([]complex128, cols)
		for col := 0; col < cols; col++ {
			tmp[col] = mat.At(row, col)
		}
		out[row] = NewKet(tmp)
	}
	return out
}
