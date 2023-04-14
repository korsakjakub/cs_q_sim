package quantum_simulator

import (
	"fmt"
	"io"
	"os"

	"gonum.org/v1/plot/plotter"
	"gopkg.in/yaml.v2"
)

type Metadata struct {
	Date           string `mapstructure:"date"`
	Simulation     string `mapstructure:"simulation"`
	Cpu            string `mapstructure:"cpu"`
	Ram            string `mapstructure:"ram"`
	CompletionTime string `mapstructure:"completiontime"`
}

type ResultsIO struct {
	Filename string   `mapstructure:"filename"`
	Metadata Metadata `mapstructure:"metadata"`
	Values   struct {
		System       System   `mapstructure:"system"`
		EigenValues  []string `mapstructure:"evalues"`
		EigenVectors []string `mapstructure:"evectors"`
	} `mapstructure:"values"`
	XYs []plotter.XYs `mapstructure:"xyss"`
}

func (r *ResultsIO) Write(conf FilesConfig) {
	path := conf.OutputsDir
	file, err := os.Create(path + r.Filename)
	if err != nil {
		parse(err)
	}

	defer file.Close()

	if b, err := yaml.Marshal(r); err != nil {
		parse(err)
	} else {
		if _, err := io.WriteString(file, string(b)); err != nil {
			parse(err)
		}
	}
	fmt.Printf("File created:\n%v%v\n", path, r.Filename)
}

func Read(conf FilesConfig, filename string) ResultsIO {
	file, err := os.Open(conf.OutputsDir + filename)
	if err != nil {
		parse(err)
	}

	var results ResultsIO

	if b, err := io.ReadAll(file); err != nil {
		parse(err)
		return ResultsIO{}
	} else {
		if err := yaml.Unmarshal(b, &results); err != nil {
			parse(err)
			return ResultsIO{}
		}
		return results
	}
}
