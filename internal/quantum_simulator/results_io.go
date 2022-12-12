package quantum_simulator

import (
	"encoding/csv"
	"fmt"
	"os"

	"gonum.org/v1/plot/plotter"
)

type Metadata struct {
	Date           string
	Simulation     string
	Cpu            string
	Ram            string
	CompletionTime string
}

type ResultsIO struct {
	Filename string
	Metadata Metadata
	XYs      plotter.XYs
}

func (r *ResultsIO) Write(conf FilesConfig) {
	path := conf.OutputsDir
	file, err := os.Create(path + r.Filename)
	if err != nil {
		parse(err)
	}

	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	meta := []string{r.Metadata.Date, r.Metadata.Cpu, r.Metadata.Ram, r.Metadata.Simulation, r.Metadata.CompletionTime}
	if err := w.Write(meta); err != nil {
		parse(err)
	}
	var data [][]string
	for _, record := range r.XYs {
		row := []string{fmt.Sprintf("%f", record.X),
			fmt.Sprintf("%f", record.Y)}
		data = append(data, row)
	}
	w.WriteAll(data)
}
