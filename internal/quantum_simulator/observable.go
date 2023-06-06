package quantum_simulator

import (
	"gonum.org/v1/gonum/mat"
)

type Observable struct {
	mat.Dense
}

func (o *Observable) ExpectationValue(state []complex128) float64 {
	oTimesStateReal := mat.NewVecDense(len(state), nil)
	oTimesStateImag := mat.NewVecDense(len(state), nil)
	realStateData := make([]float64, len(state))
	imagStateData := make([]float64, len(state))
	for i, c := range state {
		realStateData[i] = real(c)
		imagStateData[i] = imag(c)
	}

	realState := mat.NewVecDense(len(state), realStateData)
	imagState := mat.NewVecDense(len(state), imagStateData)

	oTimesStateReal.MulVec(o, realState)
	oTimesStateImag.MulVec(o, imagState)
	expectationValue := 0.0
	expectationValue += mat.Dot(realState, oTimesStateReal)
	expectationValue += mat.Dot(imagState, oTimesStateImag)
	return expectationValue
}
