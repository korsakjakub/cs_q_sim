package simulations

import (
	"fmt"
	"math"
	"strconv"
	"time"

	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/plotter"
)

func loadObservables(conf qs.PhysicsConfig) []qs.Observable {
	observables := make([]qs.Observable, len(conf.ObservablesConfig))
	ketLength := len(conf.InitialKet)
	for i, obs := range conf.ObservablesConfig {
		if obs.Slot > ketLength {
			continue
		}
		var operator *mat.Dense
		switch obs.Operator {
		case "Sz":
			operator = qs.Sz(conf.Spin)
		case "Sp":
			operator = qs.Sp(conf.Spin)
		case "Sm":
			operator = qs.Sm(conf.Spin)
		default:
			operator = qs.Id(conf.Spin)
		}
		observables[i] = qs.Observable{Dense: *qs.ManyBodyOperator(operator, obs.Slot, ketLength)}
	}
	return observables
}

func solveEigenProblem(s *qs.System) qs.Eigen {
	b := s.PhysicsConfig.BathMagneticField
	b0 := s.PhysicsConfig.CentralMagneticField
	hamiltonian := s.Hamiltonian(b0, b)

	return s.Diagonalize(hamiltonian)
}

func SpinTimeEvolution(conf qs.Config) {
	cs := qs.State{Angle: 0.0, Distance: 0.0}
	var bath []qs.State
	conf.Physics.BathCount = len(conf.Physics.InitialKet) - 1
	timeRange := conf.Physics.TimeRange
	spin := conf.Physics.Spin
	initialKet := mat.NewVecDense(int(math.Pow(2*spin+1, float64(len(conf.Physics.InitialKet)))), qs.ManyBodyVector(conf.Physics.InitialKet, int(2*spin+1)))
	observables := loadObservables(conf.Physics)
	start := time.Now()
	startTime := start.Format(time.RFC3339)

	if conf.Verbosity == "debug" {
		fmt.Println("Calculating initial states...")
	}
	for i := 0; i < conf.Physics.BathCount; i += 1 {
		bath = append(bath, qs.State{Angle: qs.PolarAngleCos(i, conf.Physics), Distance: conf.Physics.ConstantDistance})
	}
	s := &qs.System{
		CentralSpin:   cs,
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	var eigen qs.Eigen
	if diagDir := conf.Files.DiagonalizationDir; diagDir != "" {
		fmt.Printf("\nLoading diagonalization results from %v\n", diagDir)
		eigen = qs.LoadDiagonalizationSolutions(diagDir)
	} else {
		eigen = solveEigenProblem(s)
		qs.SaveDiagonalizationSolutions(eigen, conf.Files.OutputsDir+"diag-"+startTime)
	}

	if conf.Verbosity == "debug" {
		fmt.Println("Calculating time evolution...")
	}
	xyss := make([]plotter.XYs, len(observables))
	for i, observable := range observables {
		var xys plotter.XYs
		for t := 0; t < timeRange; t++ {
			evolutionTime := conf.Physics.Dt * float64(t)
			if conf.Verbosity == "debug" {
				fmt.Printf("t= %v\t(%.2f%%)\n", evolutionTime, 100.0*float64(t)/float64(timeRange))
			}
			xys = append(xys, plotter.XY{X: evolutionTime / (2.0 * math.Pi), Y: observable.ExpectationValue(qs.Evolve(initialKet, evolutionTime, eigen.EigenValues, eigen.EigenVectors))})
		}
		xyss[i] = xys
	}

	if conf.Verbosity == "debug" {
		fmt.Println("Wrapping up...")
	}
	elapsedTime := time.Since(start)
	r := qs.ResultsIO{
		Filename: startTime,
		Metadata: qs.Metadata{
			Date:           startTime,
			Simulation:     "Central spin expectation value time evolution",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsedTime.String(),
		},
		Values: struct {
			System       qs.System "mapstructure:\"system\""
			EigenValues  []string  "mapstructure:\"evalues\""
			EigenVectors []string  "mapstructure:\"evectors\""
		}{
			System:       *s,
			EigenValues:  []string{""}, // eValsToString(eigenValues),
			EigenVectors: []string{""}, // ketsToString(eigenVectors),
		},
		XYs: xyss,
	}
	r.Write(conf.Files)
}

func eValsToString(evals []float64) []string {
	output := make([]string, len(evals))
	for i, e := range evals {
		output[i] = strconv.FormatFloat(e, 'e', 8, 64)
	}
	return output
}

func ketsToString(kets *mat.Dense) []string {
	output := make([]string, kets.RawMatrix().Cols)
	for i := 0; i < kets.RawMatrix().Rows; i++ {
		ket := kets.ColView(i)
		output[i] = fmt.Sprintf("%v", ket)
	}
	return output
}
