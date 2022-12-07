package quantum_simulator

import (
	"gonum.org/v1/gonum/mat"
)

// Given an operator from 1-body Hilbert space, return the one-body operator from 'dim'-body Hilbert space, the one-body operator being in slot 'particle' <= 'dim'
func manyBody(operator *mat.Dense, particle int, dim int) *mat.Dense {
	if particle > dim {
		return operator
	}
	var m mat.Dense
	var n mat.Dense
	if particle > 0 {
		m = *Id(0.5)
		for i := 0; i < particle-1; i += 1 {
			var temp mat.Dense
			temp.Kronecker(&m, Id(0.5))
			m = temp
		}
		n.Kronecker(&m, operator)
	} else {
		n = *operator
	}
	if particle < dim-1 {
		for i := particle + 1; i < dim; i += 1 {
			var temp mat.Dense
			temp.Kronecker(&n, Id(0.5))
			n = temp
		}
	}
	return &n
}
