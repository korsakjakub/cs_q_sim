package cs_q_sim

import (
	"math"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestSystem_InteractionAt(t *testing.T) {
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
			fields: fields{CentralSpin: State{0.0, 1.0, 0.0}, Bath: []State{{0.0, 1.0, 0.0}}, PhysicsConfig: PhysicsConfig{BathDipoleMoment: 1.1e-10, AtomDipoleMoment: 1.0}},
			args:   args{0},
			want:   0.0,
		},
		{
			name:   "test",
			fields: fields{CentralSpin: State{0.0, 1.0, 0.0}, Bath: []State{{0.0, 1.0, 0.0}}, PhysicsConfig: PhysicsConfig{BathDipoleMoment: 1.1e-10, AtomDipoleMoment: 1.0}},
			args:   args{1},
			want:   0.98865,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				CentralSpin:   tt.fields.CentralSpin,
				Bath:          tt.fields.Bath,
				PhysicsConfig: tt.fields.PhysicsConfig,
			}
			if got := s.InteractionAt(tt.args.j); math.Abs(got-tt.want) > 1e-4 {
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
			fields: fields{CentralSpin: State{0.0, 1.0, 0.0}, Bath: []State{{0.0, 1.0, 0.0}}, PhysicsConfig: PhysicsConfig{BathDipoleMoment: 1.1e-10, AtomDipoleMoment: 1.0, Spin: 0.5}},
			args:   args{1},
			want: mat.NewDense(4, 4, []float64{
				0.0, 0.0, 0.0, 0.0,
				0.0, 0.0, 2 * 0.4943258, 0.0,
				0.0, 2 * 0.4943258, 0.0, 0.0,
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
		want   *mat.SymDense
	}{
		{
			name:   "test",
			fields: fields{CentralSpin: State{0.0, 1.0, 0.0}, Bath: []State{{0.0, 1.0, 0.0}}, PhysicsConfig: PhysicsConfig{BathDipoleMoment: 1.1e-10, AtomDipoleMoment: 1.0, Spin: 0.5}},
			args:   args{b0: 1.0, b: 3.0},
			want: mat.NewSymDense(4, []float64{
				2.0, 0.0, 0.0, 0.0,
				0.0, -1.0, 0.0, 0.0,
				0.0, 0.0, 1.0, 0.0,
				0.0, 0.0, 0.0, -2.0,
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
			if got := s.hamiltonianMagneticTerm(tt.args.b0, tt.args.b); !mat.EqualApprox(got, tt.want, 1e-8) {
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
			fields: fields{CentralSpin: State{0.0, 1.0, 0.0}, Bath: []State{{0.0, 1.0, 0.0}}, PhysicsConfig: PhysicsConfig{BathDipoleMoment: 1.1e-10, AtomDipoleMoment: 1.0, Spin: 0.5}},
			args:   args{b0: 1.0, b: 3.0},
			want: mat.NewDense(4, 4, []float64{
				2.0, 0.0, 0.0, 0.0,
				0.0, -1.0, 2 * 0.4943258, 0.0,
				0.0, 2 * 0.4943258, 1.0, 0.0,
				0.0, 0.0, 0.0, -2.0,
			}),
		},
		{
			name:   "bigger_system",
			fields: fields{CentralSpin: State{0.0, 1.0, 0.0}, Bath: []State{{0.0, 1.0, 0.0}, {0.0, 2.0, 0.0}}, PhysicsConfig: PhysicsConfig{BathDipoleMoment: 1.1e-10, AtomDipoleMoment: 1.0, Spin: 0.5}},
			args:   args{b0: 1.0, b: 3.0},
			want: mat.NewDense(8, 8, []float64{
				3.5, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
				0.0, 0.5, 0.0, 0.0, 2 * 0.0617907, 0.0, 0.0, 0.0,
				0.0, 0.0, 0.5, 0.0, 2 * 0.4943258, 0.0, 0.0, 0.0,
				0.0, 0.0, 0.0, -2.5, 0.0, 2 * 0.4943258, 2 * 0.0617907, 0.0,
				0.0, 2 * 0.0617907, 2 * 0.4943258, 0.0, 2.5, 0.0, 0.0, 0.0,
				0.0, 0.0, 0.0, 2 * 0.4943258, 0.0, -0.5, 0.0, 0.0,
				0.0, 0.0, 0.0, 2 * 0.0617907, 0.0, 0.0, -0.5, 0.0,
				0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, -3.5,
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
			if got := s.Hamiltonian(tt.args.b0, tt.args.b); !mat.EqualApprox(got, tt.want, 1e-4) {
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
		hamiltonian *mat.SymDense
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "2-body",
			fields: fields{CentralSpin: State{0.0, 1.0, 0.0}, Bath: []State{{0.0, 1.0, 0.0}}, PhysicsConfig: PhysicsConfig{BathDipoleMoment: 1.1e-10, AtomDipoleMoment: 1.0, Spin: 0.5}},
			args: args{
				hamiltonian: mat.NewSymDense(4, []float64{
					-1.0, 0.0, 0.0, 0.0,
					0.0, -1.0, -0.98865, 0.0,
					0.0, -0.98865, 1.0, 0.0,
					0.0, 0.0, 0.0, 1.0,
				}),
			},
		},
		{
			name: "Degenerate",
			fields: fields{
				CentralSpin: State{
					Angle:               0,
					Distance:            0,
					InteractionStrength: 0,
				},
				Bath:          []State{},
				PhysicsConfig: PhysicsConfig{},
			},
			args: args{hamiltonian: mat.NewSymDense(3, []float64{
				2.0, 0.0, 2.0,
				0.0, -2.0, 0.0,
				2.0, 0.0, -1.0,
			})},
		},
		{
			name:   "3-body",
			fields: fields{CentralSpin: State{0.0, 1.0, 0.0}, Bath: []State{{0.0, 1.0, 0.0}, {0.0, 2.0, 0.0}}, PhysicsConfig: PhysicsConfig{BathDipoleMoment: 1.1e-10, AtomDipoleMoment: 1.0, Spin: 0.5}},
			args: args{
				hamiltonian: mat.NewSymDense(8, []float64{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				CentralSpin:   tt.fields.CentralSpin,
				Bath:          tt.fields.Bath,
				PhysicsConfig: tt.fields.PhysicsConfig,
			}
			eigen := s.Diagonalize(tt.args.hamiltonian)

			for i := 0; i < eigen.EigenVectors.RawMatrix().Rows; i++ {
				left := mat.NewVecDense(len(eigen.EigenValues), nil)
				vec := eigen.EigenVectors.ColView(i)
				left.MulVec(tt.args.hamiltonian, vec)
				right := mat.NewVecDense(len(eigen.EigenValues), nil)
				right.ScaleVec(eigen.EigenValues[i], vec)

				if !mat.EqualApprox(left, right, 1e-8) {
					t.Errorf("Vector is not an eigenvector")
				}
			}
		})
	}
}

func TestSystem_diagonalize_benchmark(t *testing.T) {
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
	}{
		{
			name: "strain test",
			fields: fields{CentralSpin: State{0.0, 1.0, 0.0}, Bath: []State{
				{0.0, 1.0, 0.0}, {0.0, 2.0, 0.0}, {0.0, 1.1, 0.0}, {0.0, 1.2, 0.0}, {0.0, 1.3, 0.0},
			}, PhysicsConfig: PhysicsConfig{BathDipoleMoment: 1.1e-10, AtomDipoleMoment: 1.0, Spin: 0.5}},
			args: args{
				b0: 1.0,
				b:  3.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				CentralSpin:   tt.fields.CentralSpin,
				Bath:          tt.fields.Bath,
				PhysicsConfig: tt.fields.PhysicsConfig,
			}
			eigen := s.Diagonalize(s.Hamiltonian(tt.args.b0, tt.args.b))
			_, vecsCount := eigen.EigenVectors.Dims()

			t.Logf("num of eigvals: %v, num of eigvecs: %v", len(eigen.EigenValues), vecsCount)
		})
	}
}

func TestSystem_HamiltonianInBase(t *testing.T) {
	type fields struct {
		CentralSpin   State
		Bath          []State
		PhysicsConfig PhysicsConfig
		DownSpins     int
	}
	type args struct {
		b0      float64
		b       float64
		indices []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mat.SymDense
	}{
		{
			name: "3 -> 1",
			fields: fields{
				CentralSpin:   State{0.0, 1.0, 0.0},
				Bath:          []State{{0.0, 1.0, 0.0}, {0.0, 2.0, 0.0}},
				PhysicsConfig: PhysicsConfig{BathDipoleMoment: 1.1e-10, AtomDipoleMoment: 1.0, Spin: 0.5},
				DownSpins:     2,
			},
			args: args{
				b0:      1.0,
				b:       3.0,
				indices: []int{0, 1},
			},
			want: mat.NewSymDense(2, []float64{3.5, 0.0, 0.0, 0.5}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				CentralSpin:   tt.fields.CentralSpin,
				Bath:          tt.fields.Bath,
				PhysicsConfig: tt.fields.PhysicsConfig,
				DownSpins:     tt.fields.DownSpins,
			}
			if got := s.HamiltonianInBase(tt.args.b0, tt.args.b, tt.args.indices); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("System.HamiltonianInBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystem_HamiltonianInBasePanic(t *testing.T) {
	type fields struct {
		CentralSpin   State
		Bath          []State
		PhysicsConfig PhysicsConfig
		DownSpins     int
	}
	type args struct {
		b0      float64
		b       float64
		indices []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should panic 3 -> 0",
			fields: fields{
				CentralSpin:   State{0.0, 1.0, 0.0},
				Bath:          []State{{0.0, 1.0, 0.0}, {0.0, 2.0, 0.0}},
				PhysicsConfig: PhysicsConfig{BathDipoleMoment: 1.1e-10, AtomDipoleMoment: 1.0, Spin: 0.5},
				DownSpins:     2,
			},
			args: args{
				b0:      1.0,
				b:       3.0,
				indices: []int{0},
			},
		},
	}
	for _, tt := range tests {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		s := &System{
			CentralSpin:   tt.fields.CentralSpin,
			Bath:          tt.fields.Bath,
			PhysicsConfig: tt.fields.PhysicsConfig,
			DownSpins:     tt.fields.DownSpins,
		}
		_ = s.HamiltonianInBase(tt.args.b0, tt.args.b, tt.args.indices)
	}
}
