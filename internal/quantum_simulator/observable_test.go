package quantum_simulator

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestObservable_ExpectationValue(t *testing.T) {
	type fields struct {
		CDense mat.CDense
	}
	type args struct {
		state StateVec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "Ket(0) and spin 1/2 S_z",
			fields: fields{
				CDense: *mat.NewCDense(2, 2, []complex128{1.0, 0.0, 0.0, -1.0}),
			},
			args: args{
				state: StateVec{
					N:    2,
					Inc:  1,
					Data: []complex128{1.0, 0.0},
				},
			},
			want: 1.0,
		},
		{
			name: "Ket(1) and spin 1/2 S_z",
			fields: fields{
				CDense: *mat.NewCDense(2, 2, []complex128{1.0, 0.0, 0.0, -1.0}),
			},
			args: args{
				state: StateVec{
					N:    2,
					Inc:  1,
					Data: []complex128{0.0, 1.0},
				},
			},
			want: -1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Observable{
				CDense: tt.fields.CDense,
			}
			if got := o.ExpectationValue(tt.args.state); got != tt.want {
				t.Errorf("Observable.ExpectationValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
