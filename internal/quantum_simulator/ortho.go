package quantum_simulator

import (
	"math"

	hs "github.com/korsakjakub/cs_q_sim/internal/hilbert_space"
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
			for i, _ := range v {
				for j := 0; j < i; j++ {
					u := v[j].Scale(v[i].Dot(v[j]) / complex(math.Pow(v[j].Norm(), 2), 0.0))
					err := v[i].Sub(u)
					if err != nil {
						return
					}
				}
			}
			for i, _ := range v {
				v[i].Normalize()
			}
		}
	}
	/*
		for k, v := range o {
			if len(v) > 1 {
				var qr mat.QR
				u := hs.CMatrixFromKets(v)
				rows, cols := u.Dims()
				qr.Factorize(RealPart(u))
				q := mat.NewDense(rows, cols, nil)
				qr.QTo(q)
				o[k] = hs.KetsFromMatrix(q)
				fmt.Println(o[k])
			}
		}
	*/
}
