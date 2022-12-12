package quantum_simulator

import (
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func Plot_spectrum_mag_field(xys plotter.XYs, filename string, conf FilesConfig) {
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
