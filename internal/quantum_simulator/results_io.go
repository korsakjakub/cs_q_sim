package quantum_simulator

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
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

	metaValOf := reflect.ValueOf(r.Metadata)
	metaType := metaValOf.Type()
	metaStr := make([]string, metaType.NumField())

	for i := 0; i < metaType.NumField(); i++ {
		metaStr[i] = fmt.Sprint(metaValOf.Field(i).Interface())
	}

	if err := w.Write(metaStr); err != nil {
		parse(err)
	}

	confValOf := reflect.ValueOf(r.Config)
	confType := confValOf.Type()
	confStr := make([]string, confType.NumField())

	for i := 0; i < confType.NumField(); i++ {
		confStr[i] = fmt.Sprint(confValOf.Field(i).Interface())
	}

	if err := w.Write(confStr); err != nil {
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
	timeRange, err := strconv.Atoi(records[1][5])
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
			MoleculeMass:        moleculeMass,
			AtomMass:            atomMass,
			BathCount:           bathCount,
			Spin:                spin,
			SpectrumConfig:      SpectrumConfig{FieldRange: fieldCount},
			SpinEvolutionConfig: SpinEvolutionConfig{TimeRange: timeRange},
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
