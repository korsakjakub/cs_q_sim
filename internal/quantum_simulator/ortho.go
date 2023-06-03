package quantum_simulator

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

type Ortho []Eigen

type Eigen struct {
	Eigenvalue   float64
	EigenVectors []*mat.VecDense
}

func NewOrtho(eigenvalues []float64, eigenvectors []*mat.VecDense) Ortho {
	o := Ortho{}
	sl := make([]struct {
		eval float64
		evec *mat.VecDense
	}, len(eigenvalues))

	for i, e := range eigenvalues {
		sl[i] = struct {
			eval float64
			evec *mat.VecDense
		}{
			eval: e,
			evec: eigenvectors[i],
		}
	}
	sort.Slice(sl, func(i, j int) bool {
		return sl[i].eval < sl[j].eval
	})
	iOffset := 0
	for i := 0; i < len(sl); i++ {
		if len(sl) == 1 {
			o = append(o, Eigen{Eigenvalue: sl[0].eval, EigenVectors: []*mat.VecDense{sl[0].evec}})
			break
		}
		if i == len(sl)-1 {
			o = append(o, Eigen{Eigenvalue: sl[i].eval, EigenVectors: []*mat.VecDense{sl[i].evec}})
			break
		}
		if math.Abs(sl[i].eval-sl[i+1].eval) > 1e-10 {
			o = append(o, Eigen{Eigenvalue: sl[i].eval, EigenVectors: []*mat.VecDense{sl[i].evec}})
			continue
		}
		evecs := []*mat.VecDense{sl[i].evec}

		for j := 1; j < len(sl)-i; j++ {
			if math.Abs(sl[i].eval-sl[i+j].eval) < 1e-10 {
				evecs = append(evecs, sl[i+j].evec)
				iOffset++
			} else {
				break
			}
		}
		o = append(o, Eigen{Eigenvalue: sl[i].eval, EigenVectors: evecs})
		i += iOffset
		iOffset = 0
	}
	return o
}

/*
func (o Ortho) Orthonormalize() {
	for _, eigen := range o {
		if len(eigen.EigenVectors) > 1 {
			for i := range eigen.EigenVectors {
				for j := 0; j < i; j++ {
					u := eigen.EigenVectors[j].Scale(eigen.EigenVectors[i].Dot(eigen.EigenVectors[j]) / complex(math.Pow(eigen.EigenVectors[j].Norm(), 2), 0.0))
					err := eigen.EigenVectors[i].Sub(u)
					if err != nil {
						panic(err)
					}
				}
			}
			for i := range eigen.EigenVectors {
				eigen.EigenVectors[i].Normalize()
			}
		}
	}
}

func (o *Ortho) OrthoToEigen() ([]complex128, *mat.CDense) {
	var eigenvalues []complex128
	var kets []*hs.StateVec
	for _, eigen := range *o {
		degen_num := len(eigen.EigenVectors)
		for i := 0; i < degen_num; i++ {
			eigenvalues = append(eigenvalues, complex(eigen.Eigenvalue, 0.0))
		}
		kets = append(kets, eigen.EigenVectors...)
	}
	eigenvectors := hs.CMatrixFromKets(kets)
	return eigenvalues, eigenvectors
}
*/
