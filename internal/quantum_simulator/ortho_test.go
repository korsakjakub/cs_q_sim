package quantum_simulator

import (
	"math"
	"reflect"
	"testing"

	hs "github.com/korsakjakub/cs_q_sim/internal/hilbert_space"
	"gonum.org/v1/gonum/mat"
)

func TestNewOrtho(t *testing.T) {
	type args struct {
		eigenvalues  []float64
		eigenvectors []*hs.StateVec
	}
	tests := []struct {
		name string
		args args
		want *Ortho
	}{
		{
			name: "single",
			args: args{
				eigenvalues: []float64{1.0},
				eigenvectors: []*hs.StateVec{
					{N: 1, Inc: 1, Data: []complex128{1.0}},
				},
			},
			want: &Ortho{1.0: []*hs.StateVec{{N: 1, Inc: 1, Data: []complex128{1.0}}}},
		},
		{
			name: "double",
			args: args{
				eigenvalues: []float64{1.0, 1.0},
				eigenvectors: []*hs.StateVec{
					{N: 1, Inc: 1, Data: []complex128{1.0}},
					{N: 1, Inc: 1, Data: []complex128{0.0}},
				},
			},
			want: &Ortho{1.0: []*hs.StateVec{{N: 1, Inc: 1, Data: []complex128{1.0}}, {N: 1, Inc: 1, Data: []complex128{0.0}}}},
		},
		{
			name: "two degenerate energies, four vectors",
			args: args{
				eigenvalues: []float64{1.0, 1.0, 2.0, 2.0},
				eigenvectors: []*hs.StateVec{
					{N: 1, Inc: 1, Data: []complex128{1.0}}, {N: 1, Inc: 1, Data: []complex128{0.0}},
					{N: 1, Inc: 1, Data: []complex128{1.0}}, {N: 1, Inc: 1, Data: []complex128{0.0}},
				},
			},
			want: &Ortho{
				1.0: []*hs.StateVec{{N: 1, Inc: 1, Data: []complex128{1.0}}, {N: 1, Inc: 1, Data: []complex128{0.0}}},
				2.0: []*hs.StateVec{{N: 1, Inc: 1, Data: []complex128{1.0}}, {N: 1, Inc: 1, Data: []complex128{0.0}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrtho(tt.args.eigenvalues, tt.args.eigenvectors); !reflect.DeepEqual(got, *tt.want) {
				t.Errorf("NewOrtho() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrtho_Orthonormalize(t *testing.T) {
	tests := []struct {
		name string
		o    Ortho
		want []*hs.StateVec
	}{
		{
			name: "two",
			o: NewOrtho([]float64{1.0, 1.0}, []*hs.StateVec{
				{N: 2.0, Inc: 1.0, Data: []complex128{1.0, 0.0}},
				{N: 2.0, Inc: 1.0, Data: []complex128{0.6, 0.8}}},
			),
			want: []*hs.StateVec{
				{N: 2.0, Inc: 1.0, Data: []complex128{1.0, 0.0}},
				{N: 2.0, Inc: 1.0, Data: []complex128{0.0, 1.0}},
			},
		},
		{
			name: "3d",
			o: NewOrtho([]float64{1.0, 1.0, 1.0}, []*hs.StateVec{
				{N: 3.0, Inc: 1.0, Data: []complex128{1.0, -1.0, 1.0}},
				{N: 3.0, Inc: 1.0, Data: []complex128{1.0, 0.0, 1.0}},
				{N: 3.0, Inc: 1.0, Data: []complex128{1.0, 1.0, 2.0}},
			}),
			want: []*hs.StateVec{
				{N: 3.0, Inc: 1.0, Data: []complex128{complex(-1.0/math.Sqrt(3.0), 0.0), complex(-1.0/math.Sqrt(3.0), 0.0), complex(-1.0/math.Sqrt(3.0), 0.0)}},
				{N: 3.0, Inc: 1.0, Data: []complex128{complex(1.0/math.Sqrt(6.0), 0.0), complex(2.0/math.Sqrt(6.0), 0.0), complex(1.0/math.Sqrt(6.0), 0.0)}},
				{N: 3.0, Inc: 1.0, Data: []complex128{complex(-1.0/math.Sqrt(2.0), 0.0), 0.0, complex(1.0/math.Sqrt(2.0), 0.0)}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.Orthonormalize()
			for _, vecs := range tt.o {
				got := mat.NewDense(len(vecs[0].Data), len(vecs), nil)
				got.Mul(RealPart(hs.CMatrixFromKets(vecs)).T(), RealPart(hs.CMatrixFromKets(vecs)))
				d, _ := got.Dims()
				if !mat.EqualApprox(got, hs.Id(0.5*(float64(d)-1.0)), 1e-14) {
					t.Errorf("Orthonormalize() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
