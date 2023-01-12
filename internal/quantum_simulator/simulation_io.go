package quantum_simulator

import (
	"gonum.org/v1/gonum/mat"
)

type DiagonalizationInput struct {
	Hamiltonian *mat.Dense
	B           float64
}

type DiagonalizationResults struct {
	EigenVectors *mat.CDense
	EigenValues  []complex128
	B            float64
}

type SpinEvolutionInput struct {
	Time     float64
	DiagOuts DiagonalizationResults
}

type SpinEvolutionResults struct {
	ExpectedValue float64
	Time          float64
	MagneticField float64
}
