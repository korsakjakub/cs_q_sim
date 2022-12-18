package analysis_utilities

import (
	"fmt"
	"os"

	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func parse(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func PlotSpectrumMagField(xys plotter.XYs, filename string, conf qs.FilesConfig) {
	p := plot.New()
	s, err := plotter.NewScatter(xys)
	s.Radius = 1
	if err != nil {
		parse(err)
	}
	p.Add(s)
	wt, err := p.WriterTo(512, 512, "png")
	if err != nil {
		parse(err)
	}
	f, err := os.Create(conf.FigDir + filename)
	if err != nil {
		parse(err)
	}
	_, err = wt.WriteTo(f)
	if err != nil {
		parse(err)
	}
	if err = f.Close(); err != nil {
		parse(err)
	}
}

func PlotSpectrumMagFieldFrom(filename string, conf qs.FilesConfig) {
	r := qs.Read(conf, filename)
	PlotSpectrumMagField(r.XYs, filename, conf)
}
