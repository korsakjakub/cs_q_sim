package hilbert_space

import (
	"math"
	"math/cmplx"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestNewKet(t *testing.T) {
	type args struct {
		elements []complex128
	}
	tests := []struct {
		name string
		args args
		want *StateVec
	}{
		{
			name: "Test creating new kets",
			args: args{
				elements: []complex128{complex(0.0, 1.0), complex(1.0, 0.0)},
			},
			want: &StateVec{N: 2, Inc: 1, Data: []complex128{complex(0.0, 1.0), complex(1.0, 0.0)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKet(tt.args.elements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateVec_Dot(t *testing.T) {
	type args struct {
		v *StateVec
	}
	tests := []struct {
		name string
		u    *StateVec
		args args
		want complex128
	}{
		{
			name: "orthogonal vectors",
			u:    NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 0.0)}),
			args: args{
				v: NewKet([]complex128{complex(0.0, 0.0), complex(0.1, 0.0)}),
			},
			want: complex(0.0, 0.0),
		},
		{
			name: "the same vectors",
			u:    NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 0.5)}),
			args: args{
				v: NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 0.5)}),
			},
			want: complex(1.25, 0.0),
		},
		{
			name: "different vectors",
			u:    NewKet([]complex128{complex(1.0, 0.0), complex(0.5, 0.0)}),
			args: args{
				v: NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 0.5)}),
			},
			want: complex(1.0, 0.25),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Dot(tt.args.v); got != tt.want {
				t.Errorf("StateVec.Dot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateVec_Norm(t *testing.T) {
	tests := []struct {
		name string
		u    *StateVec
		want float64
	}{
		{
			name: "Norm of unit vector",
			u:    NewKet([]complex128{complex(1/math.Sqrt2, 0.0), complex(1/math.Sqrt2, 0.0)}),
			want: 1.0,
		},
		{
			name: "Norm of complex unit vector",
			u:    NewKet([]complex128{complex(0.0, 1/math.Sqrt2), complex(1/math.Sqrt2, 0.0)}),
			want: 1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Norm(); math.Abs(got-tt.want) > 1e-6 {
				t.Errorf("StateVec.Norm() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			if got := tt.u.Evolve(tt.args.time, tt.args.energies, tt.args.eigenBasis); cmplx.Abs(got.Dot(tt.want))-1.0 > 1e-6 {
				t.Errorf("StateVec.At() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKetFromFloats(t *testing.T) {
	type args struct {
		elements []float64
	}
	tests := []struct {
		name string
		args args
		want *StateVec
	}{
		{
			name: "ket 0",
			args: args{
				elements: []float64{1.0, 0.0, 0.0, 0.0},
			},
			want: NewKet([]complex128{complex(1.0, 0.0), complex(0.0, 0.0)}),
		},
		{
			name: "ket +",
			args: args{
				elements: []float64{1.0 / math.Sqrt2, 0.0, 1.0 / math.Sqrt2, 0.0},
			},
			want: NewKet([]complex128{complex(1.0/math.Sqrt2, 0.0), complex(1.0/math.Sqrt2, 0.0)}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KetFromFloats(tt.args.elements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KetFromFloats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKetsFromMatrix(t *testing.T) {
	type args struct {
		mat mat.CMatrix
	}
	tests := []struct {
		name string
		args args
		want []*StateVec
	}{
		{
			name: "single ket",
			args: args{
				mat: mat.NewCDense(2, 1, []complex128{complex(1.0, 0.0), complex(2.0, 0.0)}),
			},
			want: []*StateVec{
				NewKet([]complex128{complex(1.0, 0.0), complex(2.0, 0.0)}),
			},
		},
		{
			name: "three kets",
			args: args{
				mat: mat.NewCDense(2, 3, []complex128{
					complex(1.0, 0.0), complex(3.0, 0.0), complex(5.0, 0.0),
					complex(2.0, 0.0), complex(4.0, 0.0), complex(6.0, 0.0)}),
			},
			want: []*StateVec{
				NewKet([]complex128{complex(1.0, 0.0), complex(2.0, 0.0)}),
				NewKet([]complex128{complex(3.0, 0.0), complex(4.0, 0.0)}),
				NewKet([]complex128{complex(5.0, 0.0), complex(6.0, 0.0)}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KetsFromMatrix(tt.args.mat)
			if len(got) != len(tt.want) {
				t.Errorf("Outputs dimensions mismatch. Want: %v, got %v", len(got), len(tt.want))
			}
			for i := 0; i < len(got); i++ {
				if !mat.CEqual(mat.NewCDense(len(got[i].Data), 1, got[i].Data), mat.NewCDense(len(tt.want[i].Data), 1, tt.want[i].Data)) {
					t.Errorf("KetsFromMatrix() = %v, want %v", got[i], tt.want[i])
				}
			}
		})
	}
}

func TestNewKetReal(t *testing.T) {
	type args struct {
		elements []float64
	}
	tests := []struct {
		name string
		args args
		want *StateVec
	}{
		{
			name: "simple",
			args: args{
				elements: []float64{1.0, 0.0, -1.0},
			},
			want: &StateVec{N: 3, Inc: 1, Data: []complex128{complex(1.0, 0.0), complex(0.0, 0.0), complex(-1.0, 0.0)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKetReal(tt.args.elements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKetReal() = %v, want %v", got, tt.want)
			}
		})
	}
}
