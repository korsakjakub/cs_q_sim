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
	config := []string{fmt.Sprint(r.Config.MoleculeMass), fmt.Sprint(r.Config.AtomMass), fmt.Sprint(r.Config.BathCount), fmt.Sprint(r.Config.Spin), fmt.Sprint(r.Config.FieldRange)}
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

func Read(conf FilesConfig, filename string) ResultsIO {
	file, err := os.Open(conf.OutputsDir + filename)
	if err != nil {
		parse(err)
	}
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		parse(err)
	}

	moleculeMass, err := strconv.ParseFloat(records[1][0], 64)
	if err != nil {
		parse(err)
	}
	atomMass, err := strconv.ParseFloat(records[1][1], 64)
	if err != nil {
		parse(err)
	}
	bathCount, err := strconv.Atoi(records[1][2])
	if err != nil {
		parse(err)
	}
	spin, err := strconv.ParseFloat(records[1][3], 64)
	if err != nil {
		parse(err)
	}
	fieldCount, err := strconv.Atoi(records[1][4])
	if err != nil {
		parse(err)
	}

	r := ResultsIO{
		Filename: filename,
		Metadata: Metadata{
			Date:           records[0][0],
			Simulation:     records[0][1],
			Cpu:            records[0][2],
			Ram:            records[0][3],
			CompletionTime: records[0][4],
		},
		Config: PhysicsConfig{
			MoleculeMass: moleculeMass,
			AtomMass:     atomMass,
			BathCount:    bathCount,
			Spin:         spin,
			FieldRange:   fieldCount,
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
