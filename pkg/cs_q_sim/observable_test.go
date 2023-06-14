package cs_q_sim

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestObservable_ExpectationValue(t *testing.T) {
	type fields struct {
		Dense mat.Dense
	}
	type args struct {
		state []complex128
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
				Dense: *mat.NewDense(2, 2, []float64{1.0, 0.0, 0.0, -1.0}),
			},
			args: args{
				state: []complex128{1.0, 0.0},
			},
			want: 1.0,
		},
		{
			name: "Ket(1) and spin 1/2 S_z",
			fields: fields{
				Dense: *mat.NewDense(2, 2, []float64{1.0, 0.0, 0.0, -1.0}),
			},
			args: args{
				state: []complex128{0.0, 1.0},
			},
			want: -1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Observable{
				Dense: tt.fields.Dense,
			}
			if got := o.ExpectationValue(tt.args.state); math.Abs(got-tt.want) > 1e-14 {
				t.Errorf("Observable.ExpectationValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
