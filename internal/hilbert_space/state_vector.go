package hilbert_space

import (
	"errors"
	"math/cmplx"

	"gonum.org/v1/gonum/blas/cblas128"
	"gonum.org/v1/gonum/mat"
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

func NewZeroKet(dim int) *StateVec {
	out := make([]complex128, dim)
	for i, _ := range out {
		out[i] = complex(0.0, 0.0)
	}
	return NewKet(out)
}

func NewStdBasisKet(at, dim int) (*StateVec, error) {
	ket := NewZeroKet(dim)
	if at >= dim {
		return nil, errors.New("at > dim makes little sense")
	}
	ket.Data[at] = complex(1.0, 0.0)
	return ket, nil
}

func (u *StateVec) Scale(a complex128) *StateVec {
	out := u
	for i, e := range u.Data {
		out.Data[i] = a * e
	}
	return out
}

func (u *StateVec) Dot(v *StateVec) complex128 {
	return cblas128.Dotc(cblas128.Vector(*u), cblas128.Vector(*v))
}

func (u *StateVec) Norm() float64 {
	return cblas128.Nrm2(cblas128.Vector(*u))
}

func (u *StateVec) Normalize() {
	norm := complex(u.Norm(), 0.0)
	for i, e := range u.Data {
		u.Data[i] = e / norm
	}
}

func (u *StateVec) Sub(a *StateVec) error {
	if len(u.Data) != len(a.Data) {
		return errors.New("dimensions mismatch")
	}
	for i, e := range u.Data {
		u.Data[i] = e - a.Data[i]
	}
	return nil
}

// |Ψ(t)> = Σ_j exp(-i E_j t) * <E_j | Ψ(0) > * |E_j>
// Thus for k-th element we have
// (|Ψ(t)>)^k = Σ_j exp(-i E_j t) * <E_j | Ψ(0) > * (|E_j>)^k
func (u *StateVec) Evolve(time float64, energies []complex128, eigenBasis []*StateVec) *StateVec {
	out := make([]complex128, len(eigenBasis))

	for k := range u.Data { // iterate over slots of a vector
		for j, basisVector := range eigenBasis { // sum over eigenenergies
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

func CMatrixFromKets(kets []*StateVec) *mat.CDense {
	rows := len(kets)
	cols := len(kets[0].Data)
	m := mat.NewCDense(cols, rows, nil)
	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			m.Set(col, row, kets[row].Data[col])
		}
	}
	return m
}
