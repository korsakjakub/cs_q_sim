package hilbert_space

/*
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
*/
