package quantum_simulator

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

func delta(a, b float64) float64 {
	if math.Abs(a-b) < 1e-4 {
		return 1.0
	}
	return 0.0
}

func elSz(row, col float64) float64 {
	return delta(row, col) * col
}

func elSp(row, col, spin float64) float64 {
	return delta(row, col+1) * math.Sqrt(spin*(spin+1.0)-col*row)
}

func elSm(row, col, spin float64) float64 {
	return delta(row+1, col) * math.Sqrt(spin*(spin+1.0)-col*row)
}

func Id(spin float64) *mat.Dense {
	dim := int(2.0*spin + 1.0)
	data := mat.NewDense(dim, dim, nil)
	for i := 0; i < dim; i++ {
		data.Set(i, i, 1)
	}
	return data
}

func Sm(spin float64) *mat.Dense {
	dim := int(2.0*spin + 1.0)
	data := make([]float64, dim*dim)
	i := 0
	for m := spin; m >= -spin; m -= 1.0 {
		for n := spin; n >= -spin; n -= 1.0 {
			data[i] = elSm(m, n, spin)
			i += 1
		}
	}
	return mat.NewDense(dim, dim, data)
}

func Sp(spin float64) *mat.Dense {
	dim := int(2.0*spin + 1.0)
	data := make([]float64, dim*dim)
	i := 0
	for m := spin; m >= -spin; m -= 1.0 {
		for n := spin; n >= -spin; n -= 1.0 {
			data[i] = elSp(m, n, spin)
			i += 1
		}
	}
	return mat.NewDense(dim, dim, data)
}

func Sz(spin float64) *mat.Dense {
	dim := int(2.0*spin + 1.0)
	data := make([]float64, dim*dim)
	i := 0
	for m := spin; m >= -spin; m -= 1.0 {
		for n := spin; n >= -spin; n -= 1.0 {
			data[i] = elSz(m, n)
			i += 1
		}
	}
	return mat.NewDense(dim, dim, data)
}
