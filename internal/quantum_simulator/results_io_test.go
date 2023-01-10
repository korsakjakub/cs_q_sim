package quantum_simulator

import (
	"encoding/csv"
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
		want   [][]string
	}{
		{
			name: "tmp file",
			fields: fields{
				Filename: "test_output_file",
				Metadata: Metadata{"1", "2", "3", "4", "5"},
				Config:   PhysicsConfig{6, 7, 8, 9, SpectrumConfig{10}, SpinEvolutionConfig{MagneticField: 11.0}},
				XYs:      plotter.XYs{plotter.XY{X: 0.0, Y: 42.0}},
			},
			args: args{
				conf: FilesConfig{
					FigDir:     "/tmp/",
					OutputsDir: "/tmp/",
				},
			},
			want: [][]string{{"1", "2", "3", "4", "5"}, {"6", "7", "8", "9", "{10}", "{11 0 []}"}, {"0.000000", "42.000000"}},
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
			reader := csv.NewReader(f)
			reader.FieldsPerRecord = -1
			got, err := reader.ReadAll()
			if err != nil {
				t.Errorf("could not read the file: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
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
					MoleculeMass:        6,
					AtomMass:            7,
					BathCount:           8,
					Spin:                9,
					SpectrumConfig:      SpectrumConfig{FieldRange: 10},
					SpinEvolutionConfig: SpinEvolutionConfig{TimeRange: 11},
				},
				XYs: plotter.XYs{plotter.XY{X: 0.0, Y: 42.0}},
			},
		},
	}
	for _, tt := range tests {
		var lines = []string{
			"1,2,3,4,5",
			"6,7,8,9,10,11",
			"0.000000,42.000000",
		}
		f, err := os.Create(tt.args.conf.OutputsDir + tt.args.fileName)
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		for _, line := range lines {
			_, err := f.WriteString(line + "\n")
			if err != nil {
				t.Error(err)
			}
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
