package quantum_simulator

import (
	"reflect"
	"testing"
)

func TestSystem_forces(t *testing.T) {
	type fields struct {
		CentralSpin State
		Bath        []State
	}
	type args struct {
		cnf []Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				CentralSpin: tt.fields.CentralSpin,
				Bath:        tt.fields.Bath,
			}
			if got := s.forces(tt.args.cnf...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("System.forces() = %v, want %v", got, tt.want)
			}
		})
	}
}
