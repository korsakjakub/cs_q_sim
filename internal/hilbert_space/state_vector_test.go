package hilbert_space

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestStateVec_Evolve(t *testing.T) {
	type args struct {
		time       float64
		energies   []float64
		eigenBasis *mat.Dense
	}
	tests := []struct {
		name string
		u    *mat.VecDense
		args args
		want *mat.VecDense
	}{
		{
			name: "Trivial example",
			u:    mat.NewVecDense(2, []float64{1.0, 0.0}),
			args: args{
				time:       0.0,
				energies:   []float64{1.0, 1.0},
				eigenBasis: mat.NewDense(2, 2, []float64{1.0, 0.0, 0.0, 1.0}),
			},
			want: mat.NewVecDense(2, []float64{1.0, 0.0}),
		},
		{
			name: "Large trivial",
			u:    mat.NewVecDense(4, ManyBodyVector("uu", 2)),
			args: args{
				time:       0.0,
				energies:   []float64{1.0, 1.0, 1.0, 1.0},
				eigenBasis: mat.NewDense(4, 4, []float64{1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0}),
			},
			want: mat.NewVecDense(4, []float64{2.0, 0.0, 0.0, 0.0}),
		},
		{
			name: "another example",
			u:    mat.NewVecDense(2, []float64{1.0, 0.0}),
			args: args{
				time:       1.0,
				energies:   []float64{math.Pi, math.Pi},
				eigenBasis: mat.NewDense(2, 2, []float64{1.0, 0.0, 0.0, 1.0}),
			},
			want: mat.NewVecDense(2, []float64{-1.0, 0.0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Evolve(tt.u, tt.args.time, tt.args.energies, tt.args.eigenBasis); mat.Dot(got, tt.want)-math.Pow(got.Norm(2), 2) > 1e-8 {
				t.Errorf("StateVec.At() = %v, want %v", got, tt.want)
			}
		})
	}
}
