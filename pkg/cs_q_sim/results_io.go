package cs_q_sim

import (
	"fmt"
	"io"
	"os"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/plotter"
	"gopkg.in/yaml.v2"
)

type Metadata struct {
	Date           string `mapstructure:"date"`
	Simulation     string `mapstructure:"simulation"`
	SimulationId   string `mapstructure:"simulationid"`
	Cpu            string `mapstructure:"cpu"`
	Ram            string `mapstructure:"ram"`
	CompletionTime string `mapstructure:"completiontime"`
}

type ResultsIO struct {
	Filename string   `mapstructure:"filename"`
	Metadata Metadata `mapstructure:"metadata"`
	Values   struct {
		System System `mapstructure:"system"`
	} `mapstructure:"values"`
	XYs []plotter.XYs `mapstructure:"xyss"`
}

type DiagonalizationResultsIO struct {
	System       System    `mapstructure:"system"`
	EigenValues  []float64 `mapstructure:"evalues"`
	EigenVectors []float64 `mapstructure:"evectors"`
}

func (r *ResultsIO) Write(conf FilesConfig) {
	r.Filename += ".yaml"
	path := conf.OutputsDir
	file, err := os.Create(path + r.Filename)
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("couldn't close the file")
		}
	}(file)

	if b, err := yaml.Marshal(r); err != nil {
		panic(err)
	} else {
		if _, err := io.WriteString(file, string(b)); err != nil {
			panic(err)
		}
	}
	fmt.Printf("File created:\n%v%v\n", path, r.Filename)
}

func Read(conf FilesConfig, filename string) ResultsIO {
	filename += ".yaml"
	file, err := os.Open(conf.OutputsDir + filename)
	if err != nil {
		panic(err)
	}

	var results ResultsIO

	if b, err := io.ReadAll(file); err != nil {
		panic(err)
	} else {
		if err := yaml.Unmarshal(b, &results); err != nil {
			panic(err)
		}
		return results
	}
}

func LoadDiagonalizationSolutions(path string) Eigen {
	path += ".yaml"
	if _, err := os.Stat(path); err != nil {
		panic(err)
	}

	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	var eigen DiagonalizationResultsIO

	if b, err := io.ReadAll(file); err != nil {
		panic(err)
	} else {
		if err := yaml.Unmarshal(b, &eigen); err != nil {
			panic(err)
		}
		return Eigen{EigenValues: eigen.EigenValues, EigenVectors: mat.NewDense(len(eigen.EigenValues), len(eigen.EigenValues), eigen.EigenVectors)}
	}
}

func SaveDiagonalizationSolutions(eigen Eigen, s System, path string) {
	path += ".yaml"
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("couldn't close the file")
		}
	}(file)

	diagResults := DiagonalizationResultsIO{
		EigenValues:  eigen.EigenValues,
		EigenVectors: eigen.EigenVectors.RawMatrix().Data,
		System:       s,
	}

	if b, err := yaml.Marshal(diagResults); err != nil {
		panic(err)
	} else {
		if _, err := io.WriteString(file, string(b)); err != nil {
			panic(err)
		}
	}
	fmt.Printf("\nFile created:\n%v\n", path)
}
