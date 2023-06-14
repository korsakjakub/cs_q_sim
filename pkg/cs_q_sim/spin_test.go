package cs_q_sim

import (
	"math"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestSz(t *testing.T) {
	type args struct {
		spin float64
	}
	tests := []struct {
		name string
		args args
		want *mat.Dense
	}{
		{
			name: "spin-half pauli",
			args: args{spin: 0.5},
			want: mat.NewDense(2, 2, []float64{0.5, 0.0, 0.0, -0.5}),
		},
		{
			name: "spin-1 pauli",
			args: args{spin: 1.0},
			want: mat.NewDense(3, 3, []float64{1.0, 0.0, 0.0,
				0.0, 0.0, 0.0,
				0.0, 0.0, -1.0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sz(tt.args.spin); !mat.EqualApprox(got, tt.want, 1e-4) { // reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sz() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSp(t *testing.T) {
	type args struct {
		spin float64
	}
	tests := []struct {
		name string
		args args
		want *mat.Dense
	}{
		{
			name: "spin-half pauli",
			args: args{spin: 0.5},
			want: mat.NewDense(2, 2, []float64{0.0, 1.0, 0.0, 0.0}),
		},
		{
			name: "spin-1 pauli",
			args: args{spin: 1.0},
			want: mat.NewDense(3, 3, []float64{0.0, math.Sqrt2, 0.0,
				0.0, 0.0, math.Sqrt2,
				0.0, 0.0, 0.0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sp(tt.args.spin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSm(t *testing.T) {
	type args struct {
		spin float64
	}
	tests := []struct {
		name string
		args args
		want *mat.Dense
	}{
		{
			name: "spin-half pauli",
			args: args{spin: 0.5},
			want: mat.NewDense(2, 2, []float64{0.0, 0.0, 1.0, 0.0}),
		},
		{
			name: "spin-1 pauli",
			args: args{spin: 1.0},
			want: mat.NewDense(3, 3, []float64{0.0, 0.0, 0.0,
				math.Sqrt2, 0.0, 0.0,
				0.0, math.Sqrt2, 0.0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sm(tt.args.spin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestId(t *testing.T) {
	type args struct {
		spin float64
	}
	tests := []struct {
		name string
		args args
		want *mat.Dense
	}{
		{
			name: "spin-half identity",
			args: args{spin: 0.5},
			want: mat.NewDense(2, 2, []float64{1.0, 0.0, 0.0, 1.0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Id(tt.args.spin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Id() = %v, want %v", got, tt.want)
			}
		})
	}
}
