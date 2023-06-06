package quantum_simulator

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/spf13/viper"
)

type PhysicsConfig struct {
	BathDipoleMoment        float64            `mapstructure:"bathdipolemoment"`
	AtomDipoleMoment        float64            `mapstructure:"atomdipolemoment"`
	BathCount               int                `mapstructure:"bathcount"`
	Spin                    float64            `mapstructure:"spin"`
	TiltAngle               float64            `mapstructure:"tiltangle"`
	ConstantDistance        float64            `mapstructure:"constantdistance"`
	Geometry                string             `mapstructure:"geometry"`
	InteractionCoefficients []float64          `mapstructure:"interactioncoefficients"`
	BathMagneticField       float64            `mapstructure:"bathmagneticfield"`
	CentralMagneticField    float64            `mapstructure:"centralmagneticfield"`
	TimeRange               int                `mapstructure:"timerange"`
	Dt                      float64            `mapstructure:"dt"`
	InitialKet              string             `mapstructure:"initialket"`
	ObservablesConfig       []ObservableConfig `mapstructure:"observables"`
	MagneticFieldRange      int                `mapstructure:"magneticfieldrange"`
	Units                   string             `mapstructure:"units"`
}

type ObservableConfig struct {
	Operator string `mapstructure:"operator"`
	Slot     int    `mapstructure:"slot"`
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
	Simulation string        `mapstructure:"simulation"`
	Verbosity  string        `mapstructure:"verbosity"` // debug for more verbosity
	Physics    PhysicsConfig `mapstructure:"physics"`
	Files      FilesConfig   `mapstructure:"files"`
}

var vp *viper.Viper

func parse(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func LoadConfig(additionalPath []string) Config {
	vp = viper.New()
	var config Config
	for _, path := range additionalPath {
		vp.SetConfigFile(path)
		vp.AddConfigPath(filepath.Dir(path))
		err := vp.MergeInConfig()
		if err != nil {
			parse(err)
		}
	}

	err := vp.Unmarshal(&config)
	if err != nil {
		parse(err)
	}

	return config
}

func Validate(conf interface{}, fields []string) error {
	confValues := reflect.ValueOf(conf)
	for _, field := range fields {
		if !confValues.FieldByName(field).IsZero() {
			continue
		}
		return fmt.Errorf("failed to validate the struct. The field %s is required and not provided", field)
	}
	return nil
}
