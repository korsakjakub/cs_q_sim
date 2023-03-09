package quantum_simulator

import (
	"math"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	hs "github.com/korsakjakub/cs_q_sim/internal/hilbert_space"
	"gonum.org/v1/gonum/mat"
)

func TestOrtho_NewOrtho(t *testing.T) {
	type args struct {
		eigenvalues  []float64
		eigenvectors []*hs.StateVec
	}
	tests := []struct {
		name string
		args args
		want *Ortho
	}{
		{
			name: "single",
			args: args{
				eigenvalues: []float64{1.0},
				eigenvectors: []*hs.StateVec{
					{N: 1, Inc: 1, Data: []complex128{1.0}},
				},
			},
			want: &Ortho{Eigen{Eigenvalue: 1.0, EigenVectors: []*hs.StateVec{{N: 1, Inc: 1, Data: []complex128{1.0}}}}},
		},
		{
			name: "double",
			args: args{
				eigenvalues: []float64{1.0, 1.0},
				eigenvectors: []*hs.StateVec{
					{N: 1, Inc: 1, Data: []complex128{1.0}},
					{N: 1, Inc: 1, Data: []complex128{0.0}},
				},
			},
			want: &Ortho{
				Eigen{Eigenvalue: 1.0, EigenVectors: []*hs.StateVec{{N: 1, Inc: 1, Data: []complex128{1.0}}, {N: 1, Inc: 1, Data: []complex128{0.0}}}},
			},
		},
		{
			name: "two degenerate energies, four vectors",
			args: args{
				eigenvalues: []float64{1.0, 1.0, 2.0, 2.0},
				eigenvectors: []*hs.StateVec{
					{N: 1, Inc: 1, Data: []complex128{1.0}}, {N: 1, Inc: 1, Data: []complex128{0.0}},
					{N: 1, Inc: 1, Data: []complex128{1.0}}, {N: 1, Inc: 1, Data: []complex128{0.0}},
				},
			},
			want: &Ortho{
				Eigen{Eigenvalue: 1.0, EigenVectors: []*hs.StateVec{{N: 1, Inc: 1, Data: []complex128{1.0}}, {N: 1, Inc: 1, Data: []complex128{0.0}}}},
				Eigen{Eigenvalue: 2.0, EigenVectors: []*hs.StateVec{{N: 1, Inc: 1, Data: []complex128{1.0}}, {N: 1, Inc: 1, Data: []complex128{0.0}}}},
			},
		},
		{
			name: "real 3-body with degenerate values",
			args: args{
				eigenvalues: []float64{1.05031803e+02, -4.49498371e+01, 1.09127643e+03, 9.08723566e+02, -1.09127643e+03, 6.00966503e+01, 4.49498371e+01, -1.00100000e+03, -1.05031803e+02, 1.00100000e+03, -2.00100000e+03, -9.08723566e+02, -1.00100000e+03, -6.00966503e+01, 1.00100000e+03, 2.00100000e+03},
				eigenvectors: []*hs.StateVec{
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (1.29384942e-15 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (-3.84473224e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (4.80591530e-01 + 0.00000000e+00i), (-7.88170109e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (-7.03222678e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (-3.50057265e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (4.37571581e-01 + 0.00000000e+00i), (4.37571581e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (6.03185425e-18 + 0.00000000e+00i), (3.34862697e-18 + 0.00000000e+00i), (-2.19562234e-15 + 0.00000000e+00i), (-6.45923611e-18 + 0.00000000e+00i), (-4.94406963e-01 + 0.00000000e+00i), (4.94406963e-01 + 0.00000000e+00i), (3.68887921e-35 + 0.00000000e+00i), (-1.49831281e-17 + 0.00000000e+00i), (-5.05531162e-01 + 0.00000000e+00i), (5.05531162e-01 + 0.00000000e+00i), (5.55963396e-18 + 0.00000000e+00i), (2.35866771e-15 + 0.00000000e+00i), (4.44770716e-18 + 0.00000000e+00i), (2.16742656e-36 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (8.91796606e-17 + 0.00000000e+00i), (1.77993679e-17 + 0.00000000e+00i), (-7.95394975e-34 + 0.00000000e+00i), (-1.48762225e-16 + 0.00000000e+00i), (2.97451218e-17 + 0.00000000e+00i), (-2.97451218e-17 + 0.00000000e+00i), (8.51422590e-16 + 0.00000000e+00i), (4.99501028e-17 + 0.00000000e+00i), (-1.30549383e-17 + 0.00000000e+00i), (9.13259512e-18 + 0.00000000e+00i), (-8.40381124e-01 + 0.00000000e+00i), (7.71255162e-33 + 0.00000000e+00i), (-5.20221225e-01 + 0.00000000e+00i), (-1.52083674e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (4.37571581e-01 + 0.00000000e+00i), (4.37571581e-01 + 0.00000000e+00i), (1.62176985e-18 + 0.00000000e+00i), (-3.50057265e-01 + 0.00000000e+00i), (-7.31050110e-17 + 0.00000000e+00i), (7.30980871e-17 + 0.00000000e+00i), (-1.45449336e-18 + 0.00000000e+00i), (7.03222678e-01 + 0.00000000e+00i), (2.12530789e-16 + 0.00000000e+00i), (-1.74932109e-16 + 0.00000000e+00i), (-2.62218418e-17 + 0.00000000e+00i), (-3.56080329e-22 + 0.00000000e+00i), (-6.51417487e-18 + 0.00000000e+00i), (-3.90502047e-20 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (-7.10969665e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (3.46242912e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (-4.32803640e-01 + 0.00000000e+00i), (-4.32803640e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (1.98658372e-32 + 0.00000000e+00i), (1.10286614e-32 + 0.00000000e+00i), (5.59239856e-01 + 0.00000000e+00i), (-5.25121444e-32 + 0.00000000e+00i), (2.99107716e-01 + 0.00000000e+00i), (2.99107716e-01 + 0.00000000e+00i), (1.26933787e-49 + 0.00000000e+00i), (-6.32247707e-32 + 0.00000000e+00i), (-3.04126935e-01 + 0.00000000e+00i), (-3.04126935e-01 + 0.00000000e+00i), (1.83105855e-32 + 0.00000000e+00i), (-5.68624258e-01 + 0.00000000e+00i), (1.46484684e-32 + 0.00000000e+00i), (7.57814409e-51 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (5.90421858e-18 + 0.00000000e+00i), (3.27776911e-18 + 0.00000000e+00i), (1.79019052e-15 + 0.00000000e+00i), (-6.85356212e-18 + 0.00000000e+00i), (5.05531162e-01 + 0.00000000e+00i), (-5.05531162e-01 + 0.00000000e+00i), (2.80425141e-35 + 0.00000000e+00i), (9.81440552e-18 + 0.00000000e+00i), (-4.94406963e-01 + 0.00000000e+00i), (4.94406963e-01 + 0.00000000e+00i), (5.44199059e-18 + 0.00000000e+00i), (1.67205691e-15 + 0.00000000e+00i), (4.35359247e-18 + 0.00000000e+00i), (1.50604076e-36 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (-2.64390620e-33 + 0.00000000e+00i), (-1.46778341e-33 + 0.00000000e+00i), (4.24530569e-01 + 0.00000000e+00i), (2.64449363e-32 + 0.00000000e+00i), (-3.96871096e-01 + 0.00000000e+00i), (-3.96871096e-01 + 0.00000000e+00i), (-1.92853291e-50 + 0.00000000e+00i), (1.38938479e-32 + 0.00000000e+00i), (-4.00667836e-01 + 0.00000000e+00i), (-4.00667836e-01 + 0.00000000e+00i), (-2.43692073e-33 + 0.00000000e+00i), (4.28591919e-01 + 0.00000000e+00i), (-1.94953659e-33 + 0.00000000e+00i), (-1.20910800e-51 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (-4.75650635e-01 + 0.00000000e+00i), (7.88145605e-01 + 0.00000000e+00i), (-1.97504872e-18 + 0.00000000e+00i), (3.90618713e-01 + 0.00000000e+00i), (-7.10111551e-18 + 0.00000000e+00i), (7.11115050e-18 + 0.00000000e+00i), (-2.68418964e-18 + 0.00000000e+00i), (9.60190464e-16 + 0.00000000e+00i), (3.55777097e-17 + 0.00000000e+00i), (-3.68066484e-17 + 0.00000000e+00i), (5.98727391e-17 + 0.00000000e+00i), (5.62574410e-22 + 0.00000000e+00i), (1.03915347e-17 + 0.00000000e+00i), (-7.53147014e-20 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (2.61885246e-33 + 0.00000000e+00i), (1.45387465e-33 + 0.00000000e+00i), (4.28591919e-01 + 0.00000000e+00i), (-2.67446947e-32 + 0.00000000e+00i), (-4.00667836e-01 + 0.00000000e+00i), (-4.00667836e-01 + 0.00000000e+00i), (1.05641723e-50 + 0.00000000e+00i), (1.16101989e-32 + 0.00000000e+00i), (3.96871096e-01 + 0.00000000e+00i), (3.96871096e-01 + 0.00000000e+00i), (2.41382840e-33 + 0.00000000e+00i), (-4.24530569e-01 + 0.00000000e+00i), (1.93106272e-33 + 0.00000000e+00i), (5.36535544e-52 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (4.32803640e-01 + 0.00000000e+00i), (4.32803640e-01 + 0.00000000e+00i), (1.93126797e-18 + 0.00000000e+00i), (-3.46242912e-01 + 0.00000000e+00i), (-5.67208287e-17 + 0.00000000e+00i), (5.67089008e-17 + 0.00000000e+00i), (-1.51056865e-18 + 0.00000000e+00i), (-7.10969665e-01 + 0.00000000e+00i), (-1.59522481e-16 + 0.00000000e+00i), (2.02844650e-16 + 0.00000000e+00i), (-1.18914284e-17 + 0.00000000e+00i), (-7.36521294e-22 + 0.00000000e+00i), (-3.50807519e-17 + 0.00000000e+00i), (-4.44324842e-20 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (6.28112659e-01 + 0.00000000e+00i), (-5.63102927e-03 + 0.00000000e+00i), (-3.93424427e-18 + 0.00000000e+00i), (7.78102037e-01 + 0.00000000e+00i), (-1.79750798e-17 + 0.00000000e+00i), (1.79950692e-17 + 0.00000000e+00i), (1.91776118e-20 + 0.00000000e+00i), (-2.77045492e-15 + 0.00000000e+00i), (4.13834700e-17 + 0.00000000e+00i), (-1.52848502e-16 + 0.00000000e+00i), (5.13318850e-17 + 0.00000000e+00i), (1.12063319e-21 + 0.00000000e+00i), (4.12061782e-17 + 0.00000000e+00i), (5.38097638e-22 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(1.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (1.00000000e+00 + 0.00000000e+00i)}},
					{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (1.95379774e-32 + 0.00000000e+00i), (1.08466477e-32 + 0.00000000e+00i), (-5.68624258e-01 + 0.00000000e+00i), (-5.39947978e-32 + 0.00000000e+00i), (-3.04126935e-01 + 0.00000000e+00i), (-3.04126935e-01 + 0.00000000e+00i), (8.90210756e-50 + 0.00000000e+00i), (4.61263797e-32 + 0.00000000e+00i), (-2.99107716e-01 + 0.00000000e+00i), (-2.99107716e-01 + 0.00000000e+00i), (1.80083932e-32 + 0.00000000e+00i), (-5.59239856e-01 + 0.00000000e+00i), (1.44067145e-32 + 0.00000000e+00i), (4.71268627e-51 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
				},
			},
			want: &Ortho{
				Eigen{
					Eigenvalue:   -2.00100000e+03,
					EigenVectors: []*hs.StateVec{{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (2.61885246e-33 + 0.00000000e+00i), (1.45387465e-33 + 0.00000000e+00i), (4.28591919e-01 + 0.00000000e+00i), (-2.67446947e-32 + 0.00000000e+00i), (-4.00667836e-01 + 0.00000000e+00i), (-4.00667836e-01 + 0.00000000e+00i), (1.05641723e-50 + 0.00000000e+00i), (1.16101989e-32 + 0.00000000e+00i), (3.96871096e-01 + 0.00000000e+00i), (3.96871096e-01 + 0.00000000e+00i), (2.41382840e-33 + 0.00000000e+00i), (-4.24530569e-01 + 0.00000000e+00i), (1.93106272e-33 + 0.00000000e+00i), (5.36535544e-52 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}}}},
				Eigen{
					Eigenvalue: -1.09127643e+03,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (4.37571581e-01 + 0.00000000e+00i), (4.37571581e-01 + 0.00000000e+00i), (1.62176985e-18 + 0.00000000e+00i), (-3.50057265e-01 + 0.00000000e+00i), (-7.31050110e-17 + 0.00000000e+00i), (7.30980871e-17 + 0.00000000e+00i), (-1.45449336e-18 + 0.00000000e+00i), (7.03222678e-01 + 0.00000000e+00i), (2.12530789e-16 + 0.00000000e+00i), (-1.74932109e-16 + 0.00000000e+00i), (-2.62218418e-17 + 0.00000000e+00i), (-3.56080329e-22 + 0.00000000e+00i), (-6.51417487e-18 + 0.00000000e+00i), (-3.90502047e-20 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					},
				},
				Eigen{
					Eigenvalue: -1.00100000e+03,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (6.28112659e-01 + 0.00000000e+00i), (-5.63102927e-03 + 0.00000000e+00i), (-3.93424427e-18 + 0.00000000e+00i), (7.78102037e-01 + 0.00000000e+00i), (-1.79750798e-17 + 0.00000000e+00i), (1.79950692e-17 + 0.00000000e+00i), (1.91776118e-20 + 0.00000000e+00i), (-2.77045492e-15 + 0.00000000e+00i), (4.13834700e-17 + 0.00000000e+00i), (-1.52848502e-16 + 0.00000000e+00i), (5.13318850e-17 + 0.00000000e+00i), (1.12063319e-21 + 0.00000000e+00i), (4.12061782e-17 + 0.00000000e+00i), (5.38097638e-22 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (5.90421858e-18 + 0.00000000e+00i), (3.27776911e-18 + 0.00000000e+00i), (1.79019052e-15 + 0.00000000e+00i), (-6.85356212e-18 + 0.00000000e+00i), (5.05531162e-01 + 0.00000000e+00i), (-5.05531162e-01 + 0.00000000e+00i), (2.80425141e-35 + 0.00000000e+00i), (9.81440552e-18 + 0.00000000e+00i), (-4.94406963e-01 + 0.00000000e+00i), (4.94406963e-01 + 0.00000000e+00i), (5.44199059e-18 + 0.00000000e+00i), (1.67205691e-15 + 0.00000000e+00i), (4.35359247e-18 + 0.00000000e+00i), (1.50604076e-36 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
					},
				},
				Eigen{
					Eigenvalue: -9.08723566e+02,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (4.32803640e-01 + 0.00000000e+00i), (4.32803640e-01 + 0.00000000e+00i), (1.93126797e-18 + 0.00000000e+00i), (-3.46242912e-01 + 0.00000000e+00i), (-5.67208287e-17 + 0.00000000e+00i), (5.67089008e-17 + 0.00000000e+00i), (-1.51056865e-18 + 0.00000000e+00i), (-7.10969665e-01 + 0.00000000e+00i), (-1.59522481e-16 + 0.00000000e+00i), (2.02844650e-16 + 0.00000000e+00i), (-1.18914284e-17 + 0.00000000e+00i), (-7.36521294e-22 + 0.00000000e+00i), (-3.50807519e-17 + 0.00000000e+00i), (-4.44324842e-20 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}}}},
				Eigen{
					Eigenvalue: -1.05031803e+02,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (-2.64390620e-33 + 0.00000000e+00i), (-1.46778341e-33 + 0.00000000e+00i), (4.24530569e-01 + 0.00000000e+00i), (2.64449363e-32 + 0.00000000e+00i), (-3.96871096e-01 + 0.00000000e+00i), (-3.96871096e-01 + 0.00000000e+00i), (-1.92853291e-50 + 0.00000000e+00i), (1.38938479e-32 + 0.00000000e+00i), (-4.00667836e-01 + 0.00000000e+00i), (-4.00667836e-01 + 0.00000000e+00i), (-2.43692073e-33 + 0.00000000e+00i), (4.28591919e-01 + 0.00000000e+00i), (-1.94953659e-33 + 0.00000000e+00i), (-1.20910800e-51 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}}}},
				Eigen{
					Eigenvalue: -6.00966503e+01,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(1.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}}}},
				Eigen{
					Eigenvalue: -4.49498371e+01,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (-7.03222678e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (-3.50057265e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (4.37571581e-01 + 0.00000000e+00i), (4.37571581e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}}}},
				Eigen{
					Eigenvalue: 4.49498371e+01,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (1.98658372e-32 + 0.00000000e+00i), (1.10286614e-32 + 0.00000000e+00i), (5.59239856e-01 + 0.00000000e+00i), (-5.25121444e-32 + 0.00000000e+00i), (2.99107716e-01 + 0.00000000e+00i), (2.99107716e-01 + 0.00000000e+00i), (1.26933787e-49 + 0.00000000e+00i), (-6.32247707e-32 + 0.00000000e+00i), (-3.04126935e-01 + 0.00000000e+00i), (-3.04126935e-01 + 0.00000000e+00i), (1.83105855e-32 + 0.00000000e+00i), (-5.68624258e-01 + 0.00000000e+00i), (1.46484684e-32 + 0.00000000e+00i), (7.57814409e-51 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}}}},
				Eigen{
					Eigenvalue: 6.00966503e+01,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (-7.10969665e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (3.46242912e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (-4.32803640e-01 + 0.00000000e+00i), (-4.32803640e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}}}},
				Eigen{
					Eigenvalue: 1.05031803e+02,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (1.29384942e-15 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (-3.84473224e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (4.80591530e-01 + 0.00000000e+00i), (-7.88170109e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}}}},
				Eigen{
					Eigenvalue: 9.08723566e+02,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (8.91796606e-17 + 0.00000000e+00i), (1.77993679e-17 + 0.00000000e+00i), (-7.95394975e-34 + 0.00000000e+00i), (-1.48762225e-16 + 0.00000000e+00i), (2.97451218e-17 + 0.00000000e+00i), (-2.97451218e-17 + 0.00000000e+00i), (8.51422590e-16 + 0.00000000e+00i), (4.99501028e-17 + 0.00000000e+00i), (-1.30549383e-17 + 0.00000000e+00i), (9.13259512e-18 + 0.00000000e+00i), (-8.40381124e-01 + 0.00000000e+00i), (7.71255162e-33 + 0.00000000e+00i), (-5.20221225e-01 + 0.00000000e+00i), (-1.52083674e-01 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}}}},
				Eigen{
					Eigenvalue: 1.00100000e+03,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (-4.75650635e-01 + 0.00000000e+00i), (7.88145605e-01 + 0.00000000e+00i), (-1.97504872e-18 + 0.00000000e+00i), (3.90618713e-01 + 0.00000000e+00i), (-7.10111551e-18 + 0.00000000e+00i), (7.11115050e-18 + 0.00000000e+00i), (-2.68418964e-18 + 0.00000000e+00i), (9.60190464e-16 + 0.00000000e+00i), (3.55777097e-17 + 0.00000000e+00i), (-3.68066484e-17 + 0.00000000e+00i), (5.98727391e-17 + 0.00000000e+00i), (5.62574410e-22 + 0.00000000e+00i), (1.03915347e-17 + 0.00000000e+00i), (-7.53147014e-20 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}},
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i), (1.00000000e+00 + 0.00000000e+00i)}},
					},
				},
				Eigen{
					Eigenvalue: 1.09127643e+03,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (6.03185425e-18 + 0.00000000e+00i), (3.34862697e-18 + 0.00000000e+00i), (-2.19562234e-15 + 0.00000000e+00i), (-6.45923611e-18 + 0.00000000e+00i), (-4.94406963e-01 + 0.00000000e+00i), (4.94406963e-01 + 0.00000000e+00i), (3.68887921e-35 + 0.00000000e+00i), (-1.49831281e-17 + 0.00000000e+00i), (-5.05531162e-01 + 0.00000000e+00i), (5.05531162e-01 + 0.00000000e+00i), (5.55963396e-18 + 0.00000000e+00i), (2.35866771e-15 + 0.00000000e+00i), (4.44770716e-18 + 0.00000000e+00i), (2.16742656e-36 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}}}},
				Eigen{
					Eigenvalue: 2.00100000e+03,
					EigenVectors: []*hs.StateVec{
						{N: 16, Inc: 1, Data: []complex128{(0.00000000e+00 + 0.00000000e+00i), (1.95379774e-32 + 0.00000000e+00i), (1.08466477e-32 + 0.00000000e+00i), (-5.68624258e-01 + 0.00000000e+00i), (-5.39947978e-32 + 0.00000000e+00i), (-3.04126935e-01 + 0.00000000e+00i), (-3.04126935e-01 + 0.00000000e+00i), (8.90210756e-50 + 0.00000000e+00i), (4.61263797e-32 + 0.00000000e+00i), (-2.99107716e-01 + 0.00000000e+00i), (-2.99107716e-01 + 0.00000000e+00i), (1.80083932e-32 + 0.00000000e+00i), (-5.59239856e-01 + 0.00000000e+00i), (1.44067145e-32 + 0.00000000e+00i), (4.71268627e-51 + 0.00000000e+00i), (0.00000000e+00 + 0.00000000e+00i)}}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrtho(tt.args.eigenvalues, tt.args.eigenvectors); !reflect.DeepEqual(got, *tt.want) {
				for i := range got {
					t.Errorf("\n\nElement: %d\n\n", i)
					spew.Dump(got[i])
					spew.Dump((*tt.want)[i])
				}
				//t.Errorf("NewOrtho() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrtho_Orthonormalize(t *testing.T) {
	tests := []struct {
		name string
		o    Ortho
		want []*hs.StateVec
	}{
		{
			name: "two",
			o: NewOrtho([]float64{1.0, 1.0}, []*hs.StateVec{
				{N: 2.0, Inc: 1.0, Data: []complex128{1.0, 0.0}},
				{N: 2.0, Inc: 1.0, Data: []complex128{0.6, 0.8}}},
			),
			want: []*hs.StateVec{
				{N: 2.0, Inc: 1.0, Data: []complex128{1.0, 0.0}},
				{N: 2.0, Inc: 1.0, Data: []complex128{0.0, 1.0}},
			},
		},
		{
			name: "3d",
			o: NewOrtho([]float64{1.0, 1.0, 1.0}, []*hs.StateVec{
				{N: 3.0, Inc: 1.0, Data: []complex128{1.0, -1.0, 1.0}},
				{N: 3.0, Inc: 1.0, Data: []complex128{1.0, 0.0, 1.0}},
				{N: 3.0, Inc: 1.0, Data: []complex128{1.0, 1.0, 2.0}},
			}),
			want: []*hs.StateVec{
				{N: 3.0, Inc: 1.0, Data: []complex128{complex(-1.0/math.Sqrt(3.0), 0.0), complex(-1.0/math.Sqrt(3.0), 0.0), complex(-1.0/math.Sqrt(3.0), 0.0)}},
				{N: 3.0, Inc: 1.0, Data: []complex128{complex(1.0/math.Sqrt(6.0), 0.0), complex(2.0/math.Sqrt(6.0), 0.0), complex(1.0/math.Sqrt(6.0), 0.0)}},
				{N: 3.0, Inc: 1.0, Data: []complex128{complex(-1.0/math.Sqrt(2.0), 0.0), 0.0, complex(1.0/math.Sqrt(2.0), 0.0)}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.Orthonormalize()
			for _, vecs := range tt.o {
				got := mat.NewDense(len(vecs.EigenVectors[0].Data), len(vecs.EigenVectors), nil)
				got.Mul(RealPart(hs.CMatrixFromKets(vecs.EigenVectors)).T(), RealPart(hs.CMatrixFromKets(vecs.EigenVectors)))
				d, _ := got.Dims()
				if !mat.EqualApprox(got, hs.Id(0.5*(float64(d)-1.0)), 1e-14) {
					t.Errorf("Orthonormalize() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestOrtho_OrthoToEigen(t *testing.T) {
	tests := []struct {
		name  string
		o     Ortho
		want  []complex128
		want1 *mat.CDense
	}{
		{
			name:  "Simple case",
			o:     NewOrtho([]float64{1.0, 2.0}, []*hs.StateVec{hs.NewKet([]complex128{1.0, 2.0}), hs.NewKet([]complex128{3.0, 4.0})}),
			want:  []complex128{1.0, 2.0},
			want1: mat.NewCDense(2, 2, []complex128{1.0, 2.0, 3.0, 4.0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.o.OrthoToEigen()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ortho.OrthoToEigen() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Ortho.OrthoToEigen() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}