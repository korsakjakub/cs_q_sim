package quantum_simulator

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestSystem_forceAt(t *testing.T) {
	type fields struct {
		CentralSpin   State
		Bath          []State
		PhysicsConfig PhysicsConfig
	}
	type args struct {
		j int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name:   "central",
			fields: fields{CentralSpin: State{0.0, 1.0}, Bath: []State{{0.0, 1.0}}, PhysicsConfig: PhysicsConfig{MoleculeMass: "1.1e-10", AtomMass: "1.0", BathCount: "1"}},
			args:   args{0},
			want:   0.0,
		},
		{
			name:   "test",
			fields: fields{CentralSpin: State{0.0, 1.0}, Bath: []State{{0.0, 1.0}}, PhysicsConfig: PhysicsConfig{MoleculeMass: "1.1e-10", AtomMass: "1.0", BathCount: "1"}},
			args:   args{1},
			want:   -0.98865,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				CentralSpin:   tt.fields.CentralSpin,
				Bath:          tt.fields.Bath,
				PhysicsConfig: tt.fields.PhysicsConfig,
			}
			if got := s.forceAt(tt.args.j); math.Abs(got-tt.want) > 1e-4 {
				t.Errorf("System.forceAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystem_hamiltonianHeisenbergTermAt(t *testing.T) {
	type fields struct {
		CentralSpin   State
		Bath          []State
		PhysicsConfig PhysicsConfig
	}
	type args struct {
		j int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mat.Dense
	}{
		{
			name:   "first",
			fields: fields{CentralSpin: State{0.0, 1.0}, Bath: []State{{0.0, 1.0}}, PhysicsConfig: PhysicsConfig{MoleculeMass: "1.1e-10", AtomMass: "1.0", BathCount: "1", Spin: "0.5"}},
			args:   args{1},
			want: mat.NewDense(4, 4, []float64{
				0.0, 0.0, 0.0, 0.0,
				0.0, 0.0, -0.98865, 0.0,
				0.0, -0.98865, 0.0, 0.0,
				0.0, 0.0, 0.0, 0.0,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				CentralSpin:   tt.fields.CentralSpin,
				Bath:          tt.fields.Bath,
				PhysicsConfig: tt.fields.PhysicsConfig,
			}
			if got := s.hamiltonianHeisenbergTermAt(tt.args.j); !mat.EqualApprox(got, tt.want, 1e-4) {
				t.Errorf("System.hamiltonianHeisenbergTermAt() = %v, want %v", mat.Formatted(got), mat.Formatted(tt.want))
			}
		})
	}
}

func TestSystem_hamiltonianMagneticTerm(t *testing.T) {
	type fields struct {
		CentralSpin   State
		Bath          []State
		PhysicsConfig PhysicsConfig
	}
	type args struct {
		b0 float64
		b  float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mat.Dense
	}{
		{
			name:   "test",
			fields: fields{CentralSpin: State{0.0, 1.0}, Bath: []State{{0.0, 1.0}}, PhysicsConfig: PhysicsConfig{MoleculeMass: "1.1e-10", AtomMass: "1.0", BathCount: "1", Spin: "0.5"}},
			args:   args{b0: 1.0, b: 3.0},
			want: mat.NewDense(4, 4, []float64{
				-1.0, 0.0, 0.0, 0.0,
				0.0, -1.0, 0.0, 0.0,
				0.0, 0.0, 1.0, 0.0,
				0.0, 0.0, 0.0, 1.0,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				CentralSpin:   tt.fields.CentralSpin,
				Bath:          tt.fields.Bath,
				PhysicsConfig: tt.fields.PhysicsConfig,
			}
			if got := s.hamiltonianMagneticTerm(tt.args.b0, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("System.hamiltonianMagneticTerm() = %v, want %v", mat.Formatted(got), mat.Formatted(tt.want))
			}
		})
	}
}

func TestSystem_hamiltonian(t *testing.T) {
	type fields struct {
		CentralSpin   State
		Bath          []State
		PhysicsConfig PhysicsConfig
	}
	type args struct {
		b0 float64
		b  float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mat.Dense
	}{
		{
			name:   "test",
			fields: fields{CentralSpin: State{0.0, 1.0}, Bath: []State{{0.0, 1.0}}, PhysicsConfig: PhysicsConfig{MoleculeMass: "1.1e-10", AtomMass: "1.0", BathCount: "1", Spin: "0.5"}},
			args:   args{b0: 1.0, b: 3.0},
			want: mat.NewDense(4, 4, []float64{
				-1.0, 0.0, 0.0, 0.0,
				0.0, -1.0, -0.98865, 0.0,
				0.0, -0.98865, 1.0, 0.0,
				0.0, 0.0, 0.0, 1.0,
			}),
		},
		{
			name:   "bigger_system",
			fields: fields{CentralSpin: State{0.0, 1.0}, Bath: []State{{0.0, 1.0}, {0.0, 2.0}}, PhysicsConfig: PhysicsConfig{MoleculeMass: "1.1e-10", AtomMass: "1.0", BathCount: "2", Spin: "0.5"}},
			args:   args{b0: 1.0, b: 3.0},
			want: mat.NewDense(8, 8, []float64{
				-1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
				0.0, -1.0, 0.0, 0.0, -0.12358, 0.0, 0.0, 0.0,
				0.0, 0.0, -1.0, 0.0, -0.98865, 0.0, 0.0, 0.0,
				0.0, 0.0, 0.0, -1.0, 0.0, -0.98865, -0.12358, 0.0,
				0.0, -0.12358, -0.98865, 0.0, 1.0, 0.0, 0.0, 0.0,
				0.0, 0.0, 0.0, -0.98865, 0.0, 1.0, 0.0, 0.0,
				0.0, 0.0, 0.0, -0.1235, 0.0, 0.0, 1.0, 0.0,
				0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				CentralSpin:   tt.fields.CentralSpin,
				Bath:          tt.fields.Bath,
				PhysicsConfig: tt.fields.PhysicsConfig,
			}
			if got := s.hamiltonian(tt.args.b0, tt.args.b); !mat.EqualApprox(got, tt.want, 1e-4) {
				t.Errorf("System.hamiltonian() = \n%v, want \n%v", mat.Formatted(got, mat.Prefix("   "), mat.Squeeze()), mat.Formatted(tt.want))
			}
		})
	}
}

func TestSystem_diagonalize(t *testing.T) {
	type fields struct {
		CentralSpin   State
		Bath          []State
		PhysicsConfig PhysicsConfig
	}
	type args struct {
		hamiltonian  *mat.Dense
		eigenVectors chan *mat.CDense
		eigenValues  chan complex128
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mat.CDense
		want1  []complex128
	}{
		{
			name:   "2-body",
			fields: fields{CentralSpin: State{0.0, 1.0}, Bath: []State{{0.0, 1.0}}, PhysicsConfig: PhysicsConfig{MoleculeMass: "1.1e-10", AtomMass: "1.0", BathCount: "1", Spin: "0.5"}},
			args: args{
				hamiltonian: mat.NewDense(4, 4, []float64{
					-1.0, 0.0, 0.0, 0.0,
					0.0, -1.0, -0.98865, 0.0,
					0.0, -0.98865, 1.0, 0.0,
					0.0, 0.0, 0.0, 1.0,
				}),
				eigenVectors: make(chan *mat.CDense),
				eigenValues:  make(chan complex128),
			},
			want: mat.NewCDense(4, 4, []complex128{
				0.0, 0.0, 1.0, 0.0,
				-0.38, 0.924967, 0.0, 0.0,
				0.924967, 0.38, 0.0, 0.0,
				0.0, 0.0, 0.0, 1.0},
			),
			want1: []complex128{1.40621, -1.40621, -1.0, 1.0},
		},
		{
			name:   "3-body",
			fields: fields{CentralSpin: State{0.0, 1.0}, Bath: []State{{0.0, 1.0}, {0.0, 2.0}}, PhysicsConfig: PhysicsConfig{MoleculeMass: "1.1e-10", AtomMass: "1.0", BathCount: "2", Spin: "0.5"}},
			args: args{
				hamiltonian: mat.NewDense(8, 8, []float64{
					-1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
					0.0, -1.0, 0.0, 0.0, -0.12358, 0.0, 0.0, 0.0,
					0.0, 0.0, -1.0, 0.0, -0.98865, 0.0, 0.0, 0.0,
					0.0, 0.0, 0.0, -1.0, 0.0, -0.98865, -0.12358, 0.0,
					0.0, -0.12358, -0.98865, 0.0, 1.0, 0.0, 0.0, 0.0,
					0.0, 0.0, 0.0, -0.98865, 0.0, 1.0, 0.0, 0.0,
					0.0, 0.0, 0.0, -0.1235, 0.0, 0.0, 1.0, 0.0,
					0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0,
				}),
				eigenVectors: make(chan *mat.CDense),
				eigenValues:  make(chan complex128),
			},
			want: mat.NewCDense(8, 8, []complex128{
				0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0,
				0.0, 0.0, 0.0, 0.047360608952921095, 0.9922780311236191, 0.11463542937958777, 0.0, 0.0,
				0.0, 0.0, 0.0, 0.3788887040079745, -0.12403349930334975, 0.9170927112488266, 0.0, 0.0,
				-0.9242307292434829, 0.0, -0.38183771185436777, 0.0, 0.0, 0.0, 0.0, 0.0,
				0.0, 0.0, 0.0, -0.9242295833258997, 0.0, 0.38183723928558116, 0.0, 0.0,
				-0.3788897239387569, -0.12403349930335, 0.9171016481084843, 0.0, 0.0, 0.0, 0.0, 0.0,
				-0.047330077283605464, 0.9922780311236191, 0.11456233605562893, 0.0, 0.0, 0.0, 0.0, 0.0,
				0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0,
			}),
			want1: []complex128{-1.41163, 1.41163, 1.41163, -1.41163, -1.0, 1.0, -1.0, 1.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				CentralSpin:   tt.fields.CentralSpin,
				Bath:          tt.fields.Bath,
				PhysicsConfig: tt.fields.PhysicsConfig,
			}
			go s.diagonalize(tt.args.hamiltonian, tt.args.eigenVectors, tt.args.eigenValues)

			evec := <-tt.args.eigenVectors

			var eval []complex128
			for {
				v, ok := <-tt.args.eigenValues
				if !ok {
					break
				}
				eval = append(eval, v)
			}
			if evec.RawCMatrix().Cols != tt.want.RawCMatrix().Cols || evec.RawCMatrix().Rows != tt.want.RawCMatrix().Rows {
				t.Errorf("Dims got = %v, %v, Dims want %v, %v", evec.RawCMatrix().Cols, evec.RawCMatrix().Rows, tt.want.RawCMatrix().Cols, tt.want.RawCMatrix().Rows)
			}
			if !mat.CEqualApprox(evec, tt.want, 1e-4) {
				t.Errorf("System.diagonalize() got = %v, want %v", evec, tt.want)
			}
			if len(eval) != len(tt.want1) {
				t.Errorf("Dims got = %v, Dims want %v", len(eval), len(tt.want1))
			}
			sort.Slice(eval, func(i, j int) bool {
				return real(eval[i]) < real(eval[j])
			})
			sort.Slice(tt.want1, func(i, j int) bool {
				return real(tt.want1[i]) < real(tt.want1[j])
			})
			if !mat.CEqualApprox(mat.NewCDense(1, len(eval), eval), mat.NewCDense(1, len(tt.want1), tt.want1), 1e-4) {
				t.Errorf("System.diagonalize() got1 = %v, want %v", eval, tt.want1)
			}
		})
	}
}
