package quantum_simulator

import (
	"math/cmplx"

	"gonum.org/v1/gonum/mat"
)

type Observable struct {
	mat.CDense
}

func (o *Observable) ExpectationValue(state StateVec) float64 {
	sum := complex(0.0, 0.0)
	for a, stateA := range state.Data {
		for b, stateB := range state.Data {
			sum += cmplx.Conj(stateA) * o.At(a, b) * stateB
		}
	}
	return real(sum)
}
