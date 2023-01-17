package hilbert_space

import (
	"gonum.org/v1/gonum/blas/cblas128"
	"gonum.org/v1/gonum/mat"
	"math/cmplx"
)

type StateVec cblas128.Vector

func NewKet(elements []complex128) *StateVec {
	return &StateVec{N: len(elements), Inc: 1, Data: elements}
}

func NewKetReal(elements []float64) *StateVec {
	data := make([]complex128, len(elements))
	for i, el := range elements {
		data[i] = complex(el, 0.0)
	}
	return &StateVec{N: len(elements), Inc: 1, Data: data}
}

func (u *StateVec) Dot(v *StateVec) complex128 {
	return cblas128.Dotc(cblas128.Vector(*u), cblas128.Vector(*v))
}

func (u *StateVec) Norm() float64 {
	return cblas128.Nrm2(cblas128.Vector(*u))
}

func (u *StateVec) Evolve(time float64, energies []complex128, eigenBasis []*StateVec) *StateVec {
	out := make([]complex128, len(eigenBasis))

	for k, _ := range u.Data { // iterate over slots of a vector
		for j, basisVector := range eigenBasis {
			out[k] += cmplx.Exp(-energies[j]*complex(0, time)) * basisVector.Dot(u) * basisVector.Data[k]
		}
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

func KetsFromCMatrix(mat mat.CMatrix) []*StateVec {
	rows, cols := mat.Dims()
	out := make([]*StateVec, cols)
	for col := 0; col < cols; col++ {
		tmp := make([]complex128, rows)
		for row := 0; row < rows; row++ {
			tmp[row] = mat.At(row, col)
		}
		out[col] = NewKet(tmp)
	}
	return out
}

func KetsFromMatrix(mat mat.Matrix) []*StateVec {
	rows, cols := mat.Dims()
	out := make([]*StateVec, cols)
	for col := 0; col < cols; col++ {
		tmp := make([]complex128, rows)
		for row := 0; row < rows; row++ {
			tmp[row] = complex(mat.At(row, col), 0.0)
		}
		out[col] = NewKet(tmp)
	}
	return out
}
