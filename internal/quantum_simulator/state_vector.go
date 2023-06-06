package quantum_simulator

import (
	"fmt"
	"math/cmplx"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

// Evolve returns the state vector at some time t > 0
// |Ψ(t)> = Σ_j exp(-i E_j t) * <E_j | Ψ(0) > * |E_j>
// Thus for k-th element we have
// (|Ψ(t)>)^k = Σ_j exp(-i E_j t) * <E_j | Ψ(0) > * (|E_j>)^k
func Evolve(initialVector *mat.VecDense, time float64, energies []float64, eigenBasis *mat.Dense) []complex128 {
	dim := eigenBasis.ColView(0).Len()
	out := make([]complex128, dim)

	for k := range initialVector.RawVector().Data { // iterate over slots of a vector
		for j := 0; j < eigenBasis.RowView(0).Len(); j++ { // sum over energies
			basisVector := eigenBasis.ColView(j)
			out[k] += cmplx.Exp(complex(energies[j], 0)*complex(0, time)) * complex(mat.Dot(basisVector, initialVector)*basisVector.AtVec(k), 0.0)
		}
	}
	return out
}

func GetInitialBasis(particlesCount, downCount int) [][]int {

	countOnes := func(num int) int {
		count := 0
		for num != 0 {
			count += num & 1
			num >>= 1
		}
		return count
	}

	var arrays [][]int
	maxNum := (1 << particlesCount) - 1 // Maximum value for n bits

	for i := 0; i <= maxNum; i++ {
		if countOnes(i) == downCount {
			array := make([]int, particlesCount)
			binaryStr := strconv.FormatInt(int64(i), 2)
			binaryStr = fmt.Sprintf("%0*s", particlesCount, binaryStr) // Pad with leading zeros
			for j, bit := range binaryStr {
				array[j], _ = strconv.Atoi(string(bit))
			}
			arrays = append(arrays, array)
		}
	}

	return arrays
}
