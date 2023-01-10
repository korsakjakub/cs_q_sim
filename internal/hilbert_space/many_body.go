package hilbert_space

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// Given an operator from 1-body Hilbert space, return the one-body operator from 'dim'-body Hilbert space, the one-body operator being in slot 'particle' <= 'dim'
func ManyBodyOperator(operator *mat.Dense, particle int, dim int) *mat.Dense {
	spin := float64(operator.RawMatrix().Cols-1) * 0.5
	if particle > dim {
		return operator
	}
	var m mat.Dense
	var n mat.Dense
	if particle > 0 {
		m = *Id(spin)
		for i := 0; i < particle-1; i += 1 {
			var temp mat.Dense
			temp.Kronecker(&m, Id(spin))
			m = temp
		}
		n.Kronecker(&m, operator)
	} else {
		n = *operator
	}
	if particle < dim-1 {
		for i := particle + 1; i < dim; i += 1 {
			var temp mat.Dense
			temp.Kronecker(&n, Id(spin))
			n = temp
		}
	}
	return &n
}

func ManyBodyVector(states string, dim int) {
	if len(states)%2 != 0 {
		panic("Odd length of input string")
	}
	for i, state := range states {
		if state == 'u' {
			receiver.Kronecker(mat.NewDense(dim, 1, []float64{1.0, 0.0}))
		} else if state == 'd' {
			s := mat.NewDense(dim, 1, []float64{0.0, 1.0})
		} else if state == 'p' {
			s := mat.NewDense(dim, 1, []float64{1.0 / math.Sqrt2, 1.0 / math.Sqrt2})
		} else if state == 'm' {
			s := mat.NewDense(dim, 1, []float64{1.0 / math.Sqrt2, -1.0 / math.Sqrt2})
		} else {
			panic("Unknown state")
		}

	}
}
