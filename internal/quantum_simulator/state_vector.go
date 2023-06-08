package quantum_simulator

import (
	"math/cmplx"

	"gonum.org/v1/gonum/mat"
)

// Evolve returns the state vector at some time t > 0
// |Ψ(t)> = Σ_j exp(-i E_j t) * <E_j | Ψ(0) > * |E_j>
// Thus for k-th element we have
// (|Ψ(t)>)^k = Σ_j exp(-i E_j t) * <E_j | Ψ(0) > * (|E_j>)^k
func Evolve(initialVector *mat.VecDense, time float64, energies []float64, eigenBasis *mat.Dense, grammian *mat.Dense) []complex128 {
	dim := eigenBasis.ColView(0).Len()
	out := make([]complex128, dim)

	for k := range initialVector.RawVector().Data { // iterate over slots of a vector
		for j := 0; j < eigenBasis.RowView(0).Len(); j++ { // sum over energies
			basisVector := eigenBasis.ColView(j)
			out[k] += cmplx.Exp(complex(energies[j], 0)*complex(0, time)) * complex(grammian.At(0, j)*basisVector.AtVec(k), 0.0)
		}
	}
	return out
}

func countOnes(num int) int {
	count := 0
	for num != 0 {
		count += num & 1
		num >>= 1
	}
	return count
}

func Grammian(v *mat.VecDense, m *mat.Dense) *mat.Dense {
	var gram mat.Dense
	if v.Len() != m.RawMatrix().Cols {
		panic("dimensions mismatch")
	}
	gram.Mul(v.T(), m)
	return &gram
}

/*
	BasisIndices generates a list of indices of vectors that have a specific number of "downspins"
	The length of such basis is particlesCount choose downCount
*/
func BasisIndices(particlesCount, downCount int) []int {
	var indices []int
	maxNum := (1 << particlesCount) - 1 // Maximum value for n bits

	for i := 0; i <= maxNum; i++ {
		if countOnes(i) == downCount {
			indices = append(indices, i)
		}
	}
	return indices
}

/*
	RestrictMatrixToSubspace takes a matrix, and a list of indices and returns a matrix that has only blocks intersecting from the list of indices.

	Example:
	matrix = 1, 2, 3, 4,
			 5, 6, 7, 8,
			 9, 10, 11, 12,
			 13, 14, 15, 16
	indices = {0, 2}

	yields:
			 1, 3,
			 9, 11
*/
func RestrictMatrixToSubspace(matrix *mat.Dense, indices []int) *mat.Dense {
	dim := len(indices)

	selectedRows := make([]mat.Vector, dim)
	selectedCols := make([]mat.Vector, dim)
	for i, index := range indices {
		selectedRows[i] = matrix.RowView(index)
		selectedCols[i] = matrix.ColView(index)
	}

	h := mat.NewDense(dim, dim, nil)
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			h.Set(i, j, selectedRows[i].AtVec(indices[j]))
		}
	}
	return h
}
