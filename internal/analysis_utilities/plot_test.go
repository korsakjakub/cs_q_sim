package analysis_utilities

import (
	"errors"
	"os"
	"testing"

	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	"gonum.org/v1/plot/plotter"
)

func TestPlotBasic(t *testing.T) {
	type args struct {
		xys      plotter.XYs
		filename string
		conf     qs.FilesConfig
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test if image is created",
			args: args{
				xys:      plotter.XYs{plotter.XY{X: 0.0, Y: 0.0}},
				filename: "test_file.jpg",
				conf:     qs.FilesConfig{FigDir: "/tmp/"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PlotBasic(tt.args.xys, tt.args.filename, tt.args.conf)
			if _, err := os.Stat("/tmp/test_file.jpg"); errors.Is(err, os.ErrNotExist) {
				t.Errorf("file was not created")
			}
			_ = os.Remove("/tmp/test_file.jpg")
		})
	}
}

func TestPlotBasicFrom(t *testing.T) {
	type args struct {
		outputFilePath string
		figureFilePath string
		conf           qs.FilesConfig
		output         string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "testing if plot is regenerated from outputs",
			args: args{
				outputFilePath: "test_output",
				figureFilePath: "test_figure.jpg",
				conf: qs.FilesConfig{
					FigDir:     "/tmp/",
					OutputsDir: "/tmp/",
					ResultsConfig: qs.ResultsConfig{
						Cpu: "1",
						Ram: "2",
					},
				},
				output: `2022-12-19T10:29:41+01:00,test output,1,2,21.945164ms
1,1,3,0.5,100,0
99000.000000,13.465284`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Create(tt.args.conf.OutputsDir + tt.args.outputFilePath)
			if err != nil {
				t.Errorf("could not write output file")
			}
			defer f.Close()
			f.WriteString(tt.args.output)

			PlotBasicFrom(tt.args.outputFilePath, tt.args.figureFilePath, tt.args.conf)
		})
		if _, err := os.Stat(tt.args.conf.FigDir + tt.args.figureFilePath); errors.Is(err, os.ErrNotExist) {
			t.Errorf("file was not created")
		}
		_ = os.Remove(tt.args.conf.FigDir + tt.args.figureFilePath)
		_ = os.Remove(tt.args.conf.OutputsDir + tt.args.outputFilePath)
	}
}
