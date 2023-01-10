package quantum_simulator

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type PhysicsConfig struct {
	MoleculeMass        float64
	AtomMass            float64
	BathCount           int
	Spin                float64
	SpectrumConfig      SpectrumConfig      `mapstructure:"spectrum"`
	SpinEvolutionConfig SpinEvolutionConfig `mapstructure:"timeevolution"`
}

type SpectrumConfig struct {
	FieldRange int
}

type SpinEvolutionConfig struct {
	MagneticField float64
	TimeRange     int
}

type FilesConfig struct {
	FigDir        string        `mapstructure:"figdir"`
	OutputsDir    string        `mapstructure:"outputsdir"`
	ResultsConfig ResultsConfig `mapstructure:"results"`
}

type ResultsConfig struct {
	Cpu string `mapstructure:"cpu"`
	Ram string `mapstructure:"ram"`
}

type Config struct {
	Physics PhysicsConfig `mapstructure:"physics"`
	Files   FilesConfig   `mapstructure:"files"`
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
