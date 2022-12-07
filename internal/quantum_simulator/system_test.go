package quantum_simulator

import (
	"math"
	"reflect"
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
