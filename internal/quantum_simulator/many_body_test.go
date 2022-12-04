package quantum_simulator

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func Test_many_body(t *testing.T) {
	type args struct {
		operator *mat.Dense
		particle int
		dim      int
	}
	tests := []struct {
		name string
		args args
		want *mat.Dense
	}{
		{
			name: "trivial example",
			args: args{
				operator: Sz(0.5),
				particle: 1,
				dim:      1,
			},
			want: Sz(0.5),
		},
		{
			name: "first conditional check",
			args: args{
				operator: Sz(0.5),
				particle: 2,
				dim:      2,
			},
			want: mat.NewDense(4, 4, []float64{0.5, 0.0, 0.0, 0.0,
				0.0, -0.5, 0.0, 0.0,
				0.0, 0.0, 0.5, 0.0,
				0.0, 0.0, 0.0, -0.5,
			}),
		},
		{
			name: "second conditional check",
			args: args{
				operator: Sz(0.5),
				particle: 2,
				dim:      3,
			},
			want: mat.NewDense(8, 8, []float64{0.5, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
				0.0, 0.5, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
				0.0, 0.0, -0.5, 0.0, 0.0, 0.0, 0.0, 0.0,
				0.0, 0.0, 0.0, -0.5, 0.0, 0.0, 0.0, 0.0,
				0.0, 0.0, 0.0, 0.0, 0.5, 0.0, 0.0, 0.0,
				0.0, 0.0, 0.0, 0.0, 0.0, 0.5, 0.0, 0.0,
				0.0, 0.0, 0.0, 0.0, 0.0, 0.0, -0.5, 0.0,
				0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, -0.5,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := many_body(tt.args.operator, tt.args.particle, tt.args.dim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("many_body() = %v, want %v", got, tt.want)
			}
		})
	}
}
