package quantum_simulator

import (
	"math"
	"testing"
)

func almostEqual(x []float64, y []float64, eps float64) bool {
	if len(x) != len(y) {
		return false
	}
	for i := 0; i < len(x); i += 1 {
		if math.Abs(x[i]-y[i]) > eps {
			return false
		}
	}
	return true
}

func TestSystem_forces(t *testing.T) {
	type fields struct {
		CentralSpin State
		Bath        []State
	}
	type args struct {
		pc PhysicsConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []float64
	}{
		// TODO: Add test cases.
		{
			name:   "test",
			fields: fields{CentralSpin: State{0.0, 1.0, 0.0, 0.0}, Bath: []State{{0.0, 1.0, 0.0, 0.0}}},
			args:   args{PhysicsConfig{MoleculeMass: "1.1e-10", AtomMass: "1.0", BathCount: "1"}},
			want:   []float64{-0.98865},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				CentralSpin: tt.fields.CentralSpin,
				Bath:        tt.fields.Bath,
			}
			if got := s.forces(tt.args.pc); !almostEqual(got, tt.want, 1e-4) {
				t.Errorf("System.forces() = %v, want %v", got, tt.want)
			}
		})
	}
}
