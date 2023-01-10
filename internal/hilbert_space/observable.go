package hilbert_space

import (
	"math/cmplx"

	"gonum.org/v1/gonum/mat"
)

type Observable struct {
	mat.Dense
}

func (o *Observable) ExpectationValue(state *StateVec) float64 {
	sum := complex(0.0, 0.0)
	for a, stateA := range state.Data {
		for b, stateB := range state.Data {
			sum += cmplx.Conj(stateA) * complex(o.At(a, b), 0.0) * stateB
		}
	}
	return real(sum)
}
