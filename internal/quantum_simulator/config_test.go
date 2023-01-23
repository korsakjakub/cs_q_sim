package quantum_simulator

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		additionalPath []string
		args           []string
		filename       string
		lines          []string
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			name: "tmp file",
			args: args{
				additionalPath: []string{"/tmp"},
				args:           []string{"config123217318973", "yaml"},
				filename:       "/tmp/config123217318973.yaml",
				lines: []string{
					"physics:",
					"  bathdipolemoment: 1.0",
					"  atomdipolemoment: 1.0",
					"  spin: 1.0",
					"  tiltangle: 0.0",
					"  spectrum:",
					"    magneticfieldrange: 1",
					"    bathcount: 10",
					"  timeevolution:",
					"    bathmagneticfield: 1",
					"    centralmagneticfield: 1",
					"    timerange: 2",
					"    dt: 1.0",
					"    initialket: uu",
					"    observables:",
					"      - operator: Sz",
					"        slot: 0",
					"      - operator: Sz",
					"        slot: 1",
					"files:",
					"  outputsdir: test/",
					"  figdir: figtest/",
					"  results:",
					"    cpu: 2",
					"    ram: 3",
				},
			},
			want: Config{PhysicsConfig{
				BathDipoleMoment: 1.0,
				AtomDipoleMoment: 1.0,
				Spin:             1,
				TiltAngle:        0.0,
				SpectrumConfig: SpectrumConfig{
					MagneticFieldRange: 1,
					BathCount:          10,
				},
				SpinEvolutionConfig: SpinEvolutionConfig{
					BathMagneticField:    1,
					CentralMagneticField: 1,
					TimeRange:            2,
					Dt:                   1.0,
					InitialKet:           "uu",
					ObservablesConfig: []ObservableConfig{
						{Operator: "Sz", Slot: 0},
						{Operator: "Sz", Slot: 1},
					},
				},
			}, FilesConfig{
				OutputsDir: "test/",
				FigDir:     "figtest/",
				ResultsConfig: ResultsConfig{
					Cpu: "2",
					Ram: "3",
				},
			}},
		},
	}
	for _, tt := range tests {
		f, err := os.Create(tt.args.filename)
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		for _, line := range tt.args.lines {
			_, err := f.WriteString(line + "\n")
			if err != nil {
				t.Error(err)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadConfig(tt.args.additionalPath, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
		err = os.Remove(tt.args.filename)
		if err != nil {
			t.Error(err)
		}
	}
}
