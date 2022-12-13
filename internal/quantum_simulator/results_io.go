package quantum_simulator

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

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
	Config   PhysicsConfig
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

	meta := []string{r.Metadata.Date, r.Metadata.Simulation, r.Metadata.Cpu, r.Metadata.Ram, r.Metadata.CompletionTime}
	if err := w.Write(meta); err != nil {
		parse(err)
	}
	config := []string{r.Config.MoleculeMass, r.Config.AtomMass, r.Config.BathCount, r.Config.Spin, r.Config.FieldRange}
	if err := w.Write(config); err != nil {
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

func Read(conf FilesConfig, fileName string) ResultsIO {
	file, err := os.Open(conf.OutputsDir + fileName)
	if err != nil {
		parse(err)
	}
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		parse(err)
	}

	r := ResultsIO{
		Filename: fileName,
		Metadata: Metadata{
			Date:           records[0][0],
			Simulation:     records[0][1],
			Cpu:            records[0][2],
			Ram:            records[0][3],
			CompletionTime: records[0][4],
		},
		Config: PhysicsConfig{
			MoleculeMass: records[1][0],
			AtomMass:     records[1][1],
			BathCount:    records[1][2],
			Spin:         records[1][3],
			FieldRange:   records[1][4],
		},
	}
	for i, record := range records {
		if i < 2 {
			continue
		}
		x, _ := strconv.ParseFloat(record[0], 64)
		y, _ := strconv.ParseFloat(record[1], 64)
		r.XYs = append(r.XYs, plotter.XY{X: x, Y: y})
	}
	return r
}
