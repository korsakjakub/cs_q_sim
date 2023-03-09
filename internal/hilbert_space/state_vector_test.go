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

func TestCMatrixFromKets(t *testing.T) {
	type args struct {
		kets []*StateVec
	}
	tests := []struct {
		name string
		args args
		want *mat.CDense
	}{
		{
			name: "one",
			args: args{
				kets: []*StateVec{{
					N:    2,
					Inc:  1,
					Data: []complex128{1.0, 0.0},
				}},
			},
			want: mat.NewCDense(2, 1, []complex128{1.0, 0.0}),
		},
		{
			name: "two",
			args: args{
				kets: []*StateVec{
					{N: 2, Inc: 1, Data: []complex128{1.0, 3.0}},
					{N: 2, Inc: 1, Data: []complex128{2.0, 4.0}},
				},
			},
			want: mat.NewCDense(2, 2, []complex128{1.0, 2.0, 3.0, 4.0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CMatrixFromKets(tt.args.kets); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CMatrixFromKets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateVec_Normalize(t *testing.T) {
	tests := []struct {
		name string
		u    *StateVec
		want *StateVec
	}{
		{
			name: "number",
			u:    NewKet([]complex128{5.0}),
			want: NewKet([]complex128{1.0}),
		},
		{
			name: "3D",
			u:    NewKet([]complex128{2.0, 3.0, 6.0}),
			want: NewKet([]complex128{2.0 / 7.0, 3.0 / 7.0, 6.0 / 7.0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.Normalize()
			if got := tt.u; !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("Normalize() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}

func TestStateVec_Sub(t *testing.T) {
	type args struct {
		a *StateVec
	}
	tests := []struct {
		name    string
		u       *StateVec
		args    args
		want    *StateVec
		wantErr bool
	}{
		{
			name:    "2D",
			u:       NewKet([]complex128{1.0, 2.0}),
			args:    args{a: NewKet([]complex128{3.0, 5.0})},
			want:    NewKet([]complex128{-2.0, -3.0}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Sub(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Sub() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.u.Data; !reflect.DeepEqual(got, tt.want.Data) {
				t.Errorf("Sub() error = %v, want %v", got, tt.want.Data)
			}
		})
	}
}

func TestNewZeroKet(t *testing.T) {
	type args struct {
		dim int
	}
	tests := []struct {
		name string
		args args
		want *StateVec
	}{
		{
			name: "d=2",
			args: args{dim: 2},
			want: &StateVec{N: 2, Inc: 1, Data: []complex128{0.0, 0.0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewZeroKet(tt.args.dim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewZeroKet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStdBasisKet(t *testing.T) {
	type args struct {
		at  int
		dim int
	}
	tests := []struct {
		name    string
		args    args
		want    *StateVec
		wantErr bool
	}{
		{
			name: "2D",
			args: args{
				at:  1,
				dim: 2,
			},
			want:    NewKet([]complex128{0.0, 1.0}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStdBasisKet(tt.args.at, tt.args.dim)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStdBasisKet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStdBasisKet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateVec_Scale(t *testing.T) {
	type args struct {
		a complex128
	}
	tests := []struct {
		name string
		u    *StateVec
		args args
		want *StateVec
	}{
		{
			name: "scale by 2",
			u:    NewKet([]complex128{1.0, 2.0}),
			args: args{
				a: complex(2.0, 0.0),
			},
			want: NewKet([]complex128{2.0, 4.0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.Scale(tt.args.a)
			if got := tt.u; !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("Scale() got = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}
