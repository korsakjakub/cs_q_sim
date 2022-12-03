package quantum_simulator

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Physics struct {
		MoleculeMass string `mapstructure:"moleculemass"`
		AtomMass     string `mapstructure:"atommass"`
	} `mapstructure:"physics"`
	Files struct {
		OutputDir string `mapstructure:"outputdir"`
	} `mapstructure:"files"`
}

var vp *viper.Viper

func parse(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func LoadConfig(additionalPath []string, args ...string) Config {
	vp = viper.New()
	var config Config
	if len(args) > 0 {
		vp.SetConfigName(args[0])
		vp.SetConfigType(args[1])
	} else {
		vp.SetConfigName("config")
		vp.SetConfigType("yaml")
	}
	vp.AddConfigPath("./")
	for _, path := range additionalPath {
		vp.AddConfigPath(path)
	}

	err := vp.ReadInConfig()
	if err != nil {
		parse(err)
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		parse(err)
	}

	return config
}
