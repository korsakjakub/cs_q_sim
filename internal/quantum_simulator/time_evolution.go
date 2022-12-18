package quantum_simulator

import "gonum.org/v1/gonum/mat"

func Unitary(eigen mat.Eigen, time float64) *mat.CDense {
	var eigenvalues []complex128
	var eigenvectors *mat.CDense
	eigen.Values(eigenvalues)
	eigen.VectorsTo(eigenvectors)
	output := mat.NewCDense(1, 1, nil)
	return output
}
