package quantum_simulator

import (
	"gonum.org/v1/gonum/mat"
)

type Observable struct {
	mat.Dense
}

func (o *Observable) ExpectationValue(state *mat.VecDense) float64 {
	oTimesState := mat.NewVecDense(state.Len(), nil)
	oTimesState.MulVec(o, state)
	expectationValue := mat.Dot(state, oTimesState)
	return expectationValue
}
