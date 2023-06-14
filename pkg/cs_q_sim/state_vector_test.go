package cs_q_sim

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/mat"
)

/*
func TestStateVec_Evolve(t *testing.T) {
	type args struct {
		time       float64
		energies   []complex128
		eigenBasis []*StateVec
	}
	tests := []struct {
		name string
		u    *StateVec
		args args
		want *StateVec
	}{
		{
			name: "Trivial example",
			u:    NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 0.0)}),
			args: args{
				time:       0.0,
				energies:   []complex128{complex(1.0, 0.0), complex(1.0, 0.0)},
				eigenBasis: []*StateVec{NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 1.0)}), NewKet([]complex128{complex(0.0, 0.0), complex(1.0, 0.0)})},
			},
			want: NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 0.0)}),
		},
		{
			name: "Large trivial",
			u:    NewKetReal(ManyBodyVector("uu", 2)),
			args: args{
				time:       0.0,
				energies:   []complex128{complex(1.0, 0.0), complex(1.0, 0.0), complex(1.0, 0.0), complex(1.0, 0.0)},
				eigenBasis: []*StateVec{NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 1.0), complex(0.0, 0.0), complex(0.0, 0.0)}), NewKet([]complex128{complex(0.0, 0.0), complex(1.0, 0.0), complex(0.0, 0.0), complex(0.0, 0.0)}), NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 1.0), complex(0.0, 0.0), complex(0.0, 0.0)}), NewKet([]complex128{complex(0.0, 0.0), complex(1.0, 0.0), complex(0.0, 0.0), complex(0.0, 0.0)})},
			},
			want: &StateVec{N: 4, Inc: 1, Data: []complex128{
				complex(2.0, 0.0), complex(0.0, 2.0), complex(0.0, 0.0), complex(0.0, 0.0),
			}},
		},
		{
			name: "another example",
			u:    NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 0.0)}),
			args: args{
				time:       1.0,
				energies:   []complex128{complex(math.Pi, 0.0), complex(math.Pi, 0.0)},
				eigenBasis: []*StateVec{NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 0.0)}), NewKet([]complex128{complex(0.0, 0.0), complex(1.0, 0.0)})},
			},
			want: NewKet([]complex128{complex(-1.0, 0.0), complex(0.0, 0.0)}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Evolve(tt.args.time, tt.args.energies, tt.args.eigenBasis); cmplx.Abs(got.Dot(tt.want))-math.Pow(got.Norm(), 2) > 1e-6 {
				t.Errorf("StateVec.At() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/

func TestBasisIndices(t *testing.T) {
	type args struct {
		particlesCount int
		downCount      int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "4 and 0",
			args: args{
				particlesCount: 4,
				downCount:      0,
			},
			want: []int{0},
		},
		{
			name: "4 and 1",
			args: args{
				particlesCount: 4,
				downCount:      1,
			},
			want: []int{1, 2, 4, 8},
		},
		{
			name: "4 and 2",
			args: args{
				particlesCount: 4,
				downCount:      2,
			},
			want: []int{3, 5, 6, 9, 10, 12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BasisIndices(tt.args.particlesCount, tt.args.downCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BasisIndices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestrictMatrixToSubspace(t *testing.T) {
	type args struct {
		matrix  *mat.Dense
		indices []int
	}
	tests := []struct {
		name string
		args args
		want *mat.Dense
	}{
		{
			name: "3x3 -> 2x2",
			args: args{
				matrix: mat.NewDense(3, 3, []float64{
					1, 2, 3,
					4, 5, 6,
					7, 8, 9,
				}),
				indices: []int{0, 2},
			},
			want: mat.NewDense(2, 2, []float64{
				1, 3,
				7, 9,
			}),
		},
		{
			name: "5x5 -> 3x3",
			args: args{
				matrix: mat.NewDense(5, 5, []float64{
					1, 2, 3, 4, 5,
					6, 7, 8, 9, 10,
					11, 12, 13, 14, 15,
					16, 17, 18, 19, 20,
					21, 22, 23, 24, 25,
				}),
				indices: []int{1, 2, 4},
			},
			want: mat.NewDense(3, 3, []float64{
				7, 8, 10,
				12, 13, 15,
				22, 23, 25,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RestrictMatrixToSubspace(tt.args.matrix, tt.args.indices); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RestrictMatrixToSubspace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammian(t *testing.T) {
	type args struct {
		v *mat.VecDense
		m *mat.Dense
	}
	tests := []struct {
		name string
		args args
		want *mat.Dense
	}{
		{
			name: "4x1, 4x4",
			args: args{
				v: mat.NewVecDense(4, []float64{1, 2, 3, 4}),
				m: mat.NewDense(4, 4, []float64{
					2, 5, 7, 1, 3, 9, 6, 2, 4, 8, 1, 3, 5, 3, 2, 6,
				}),
			},
			want: mat.NewDense(1, 4, []float64{40, 59, 30, 38}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Grammian(tt.args.v, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Grammian() = %v, want %v", got, tt.want)
			}
		})
	}
}
