package quantum_simulator

import (
	"io"
	"os"
	"reflect"
	"testing"

	"gonum.org/v1/plot/plotter"
)

func TestResultsIO_Write(t *testing.T) {
	type fields struct {
		Filename string
		Metadata Metadata
		Config   PhysicsConfig
		XYs      plotter.XYs
	}
	type args struct {
		conf FilesConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "tmp file",
			fields: fields{
				Filename: "test_output_file",
				Metadata: Metadata{"1", "2", "3", "4", "5"},
				Config: PhysicsConfig{6, 7, 8, 9.0, SpectrumConfig{10, 11}, SpinEvolutionConfig{
					MagneticField: 12.0,
					TimeRange:     13.0,
					Dt:            14.0,
					InitialKet:    "15",
					ObservablesConfig: []ObservableConfig{
						{
							Operator: "16",
							Slot:     17,
						},
					},
				}},
				XYs: plotter.XYs{plotter.XY{X: 0.0, Y: 42.0}, plotter.XY{X: 1.0, Y: 68.0}},
			},
			args: args{
				conf: FilesConfig{
					FigDir:     "/tmp/",
					OutputsDir: "/tmp/",
				},
			},
			want: "filename: test_output_file\n" +
				"metadata:\n" +
				"  date: \"1\"\n" +
				"  simulation: \"2\"\n" +
				"  cpu: \"3\"\n" +
				"  ram: \"4\"\n" +
				"  completiontime: \"5\"\n" +
				"config:\n" +
				"  moleculemass: 6\n" +
				"  atommass: 7\n" +
				"  spin: 8\n" +
				"  tiltangle: 9\n" +
				"  spectrumconfig:\n" +
				"    bathcount: 10\n" +
				"    fieldrange: 11\n" +
				"  spinevolutionconfig:\n" +
				"    magneticfield: 12\n" +
				"    timerange: 13\n" +
				"    dt: 14\n" +
				"    initialket: \"15\"\n" +
				"    observablesconfig:\n" +
				"    - operator: \"16\"\n" +
				"      slot: 17\n" +
				"xys:\n" +
				"- x: 0\n" +
				"  \"y\": 42\n" +
				"- x: 1\n" +
				"  \"y\": 68\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ResultsIO{
				Filename: tt.fields.Filename,
				Metadata: tt.fields.Metadata,
				Config:   tt.fields.Config,
				XYs:      tt.fields.XYs,
			}
			r.Write(tt.args.conf)
			f, err := os.Open(tt.args.conf.OutputsDir + tt.fields.Filename)
			if err != nil {
				t.Errorf("could not open file: %v", err)
			}
			b, err := io.ReadAll(f)
			got := string(b)
			if err != nil {
				t.Errorf("could not read the file: %v", err)
			}
			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
			err = os.Remove(tt.args.conf.OutputsDir + tt.fields.Filename)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestRead(t *testing.T) {
	type args struct {
		conf     FilesConfig
		fileName string
		lines    string
	}
	tests := []struct {
		name string
		args args
		want ResultsIO
	}{
		{
			name: "read",
			args: args{
				conf: FilesConfig{
					FigDir:     "/tmp/",
					OutputsDir: "/tmp/",
				},
				fileName: "test_read",
				lines: "filename: test_read\n" +
					"metadata:\n" +
					"  date: \"1\"\n" +
					"  simulation: \"2\"\n" +
					"  cpu: \"3\"\n" +
					"  ram: \"4\"\n" +
					"  completiontime: \"5\"\n" +
					"config:\n" +
					"  moleculemass: 6\n" +
					"  atommass: 7\n" +
					"  spin: 8\n" +
					"  tiltangle: 9\n" +
					"  spectrumconfig:\n" +
					"    bathcount: 10\n" +
					"    fieldrange: 11\n" +
					"  spinevolutionconfig:\n" +
					"    magneticfield: 12\n" +
					"    timerange: 13\n" +
					"    dt: 14\n" +
					"    initialket: \"15\"\n" +
					"    observablesconfig:\n" +
					"    - operator: \"16\"\n" +
					"      slot: 17\n" +
					"xys:\n" +
					"- x: 0\n" +
					"  \"y\": 42\n" +
					"- x: 1\n" +
					"  \"y\": 68\n",
			},
			want: ResultsIO{
				Filename: "test_read",
				Metadata: Metadata{
					Date:           "1",
					Simulation:     "2",
					Cpu:            "3",
					Ram:            "4",
					CompletionTime: "5",
				},
				Config: PhysicsConfig{
					MoleculeMass:   6,
					AtomMass:       7,
					Spin:           8,
					TiltAngle:      9,
					SpectrumConfig: SpectrumConfig{BathCount: 10, FieldRange: 11},
					SpinEvolutionConfig: SpinEvolutionConfig{
						MagneticField: 12,
						TimeRange:     13,
						Dt:            14,
						InitialKet:    "15",
						ObservablesConfig: []ObservableConfig{
							{
								Operator: "16",
								Slot:     17,
							},
						},
					},
				},
				XYs: plotter.XYs{
					plotter.XY{X: 0.0, Y: 42.0},
					plotter.XY{X: 1.0, Y: 68.0},
				},
			},
		},
	}
	for _, tt := range tests {
		var lines = tt.args.lines
		f, err := os.Create(tt.args.conf.OutputsDir + tt.args.fileName)
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		if _, err := f.WriteString(lines); err != nil {
			t.Error(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := Read(tt.args.conf, tt.args.fileName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
			err = os.Remove(tt.args.conf.OutputsDir + tt.args.fileName)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
