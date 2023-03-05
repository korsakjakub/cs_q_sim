package quantum_simulator

import (
	"math"

	hs "github.com/korsakjakub/cs_q_sim/internal/hilbert_space"
	"gonum.org/v1/gonum/mat"
)

type Ortho map[float64][]*hs.StateVec

func NewOrtho(eigenvalues []float64, eigenvectors []*hs.StateVec) Ortho {
	o := Ortho{}
	for i, e := range eigenvalues {
		if val, ok := o[e]; ok {
			o[e] = append(val, eigenvectors[i])
		} else {
			o[e] = []*hs.StateVec{eigenvectors[i]}
		}
	}
	return o
}

func (o Ortho) Orthonormalize() {
	for _, v := range o {
		if len(v) > 1 {
			for i := range v {
				for j := 0; j < i; j++ {
					u := v[j].Scale(v[i].Dot(v[j]) / complex(math.Pow(v[j].Norm(), 2), 0.0))
					err := v[i].Sub(u)
					if err != nil {
						return
					}
				}
			}
			for i := range v {
				v[i].Normalize()
			}
		}
	}
}

func (o *Ortho) OrthoToEigen() ([]complex128, *mat.CDense) {
	var eigenvalues []complex128
	var kets []*hs.StateVec
	for k, v := range *o {
		degen_num := len(v)
		for i := 0; i < degen_num; i++ {
			eigenvalues = append(eigenvalues, complex(k, 0.0))
		}
		kets = append(kets, v...)
	}
	eigenvectors := hs.CMatrixFromKets(kets)
	return eigenvalues, eigenvectors
}
