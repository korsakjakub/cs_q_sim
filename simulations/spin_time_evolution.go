package simulations

import (
	"fmt"
	"math"
	"strconv"
	"time"

	hs "github.com/korsakjakub/cs_q_sim/internal/hilbert_space"
	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/plotter"
)

func loadObservables(conf qs.PhysicsConfig) []hs.Observable {
	observables := make([]hs.Observable, len(conf.ObservablesConfig))
	ketLength := len(conf.InitialKet)
	for i, obs := range conf.ObservablesConfig {
		if obs.Slot > ketLength {
			continue
		}
		var operator *mat.Dense
		switch obs.Operator {
		case "Sz":
			operator = hs.Sz(conf.Spin)
		case "Sp":
			operator = hs.Sp(conf.Spin)
		case "Sm":
			operator = hs.Sm(conf.Spin)
		default:
			operator = hs.Id(conf.Spin)
		}
		observables[i] = hs.Observable{Dense: mat.Dense(*hs.ManyBodyOperator(operator, obs.Slot, ketLength))}
	}
	return observables
}

func SpinTimeEvolution(conf qs.Config) {
	cs := qs.State{Angle: 0.0, Distance: 0.0}
	var bath []qs.State
	conf.Physics.BathCount = len(conf.Physics.InitialKet) - 1
	bc := conf.Physics.BathCount
	timeRange := conf.Physics.TimeRange
	spin := conf.Physics.Spin
	initialKet := hs.NewKetReal(hs.ManyBodyVector(conf.Physics.InitialKet, int(2*spin+1)))
	observables := loadObservables(conf.Physics)

	if conf.Verbosity == "debug" {
		fmt.Println("Calculating initial states...")
	}
	start := time.Now()
	for i := 0; i < bc; i += 1 {
		bath = append(bath, qs.State{Angle: qs.PolarAngleCos(i, conf.Physics), Distance: conf.Physics.ConstantDistance})
	}

	s := &qs.System{
		CentralSpin:   cs,
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	if conf.Verbosity == "debug" {
		fmt.Println("Preparing the Hamiltonian...")
	}
	b := conf.Physics.BathMagneticField
	b0 := conf.Physics.CentralMagneticField
	diagJob := qs.DiagonalizationInput{Hamiltonian: s.Hamiltonian(b0, b), B: b}
	diagOuts := make(chan qs.DiagonalizationResults)

	if conf.Verbosity == "debug" {
		fmt.Println(diagJob.Hamiltonian)
		fmt.Println("Diagonalizing...")
	}
	go s.Diagonalize(diagJob, diagOuts)

	diag := <-diagOuts
	close(diagOuts)

	start_time := start.Format(time.RFC3339)

	if conf.Verbosity == "debug" {
		fmt.Println("Calculating time evolution...")
	}
	xyss := make([]plotter.XYs, len(observables))
	for i, observable := range observables {
		var xys plotter.XYs
		for t := 0; t < timeRange; t += 1 {
			time := conf.Physics.Dt * float64(t)
			if conf.Verbosity == "debug" {
				fmt.Printf("t= %.4f\t(%.2f%%)\n", time, 100.0*float64(t)/float64(timeRange))
			}
			xys = append(xys, plotter.XY{X: time / (2.0 * math.Pi), Y: observable.ExpectationValue(initialKet.Evolve(time, diag.EigenValues, hs.KetsFromCMatrix(diag.EigenVectors)))})
		}
		xyss[i] = xys
	}

	if conf.Verbosity == "debug" {
		fmt.Println("Wrapping up...")
	}
	elapsed_time := time.Since(start)
	r := qs.ResultsIO{
		Filename: start_time,
		Metadata: qs.Metadata{
			Date:           start_time,
			Simulation:     "Central spin expectation value time evolution",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsed_time.String(),
		},
		Values: struct {
			System       qs.System "mapstructure:\"system\""
			EigenValues  []string  "mapstructure:\"evalues\""
			EigenVectors []string  "mapstructure:\"evectors\""
		}{
			System:       *s,
			EigenValues:  eValsToString(diag.EigenValues),
			EigenVectors: ketsToString(hs.KetsFromCMatrix(diag.EigenVectors)), //eValsToString(diag.EigenVectors.RawCMatrix().Data),
		},
		XYs: xyss,
	}
	r.Write(conf.Files)
}

func eValsToString(evals []complex128) []string {
	output := make([]string, len(evals))
	for i, e := range evals {
		output[i] = strconv.FormatComplex(e, 'e', 8, 128)
	}
	return output
}

func ketsToString(kets []*hs.StateVec) []string {
	output := make([]string, len(kets))
	for i, ket := range kets {
		outData := "[]complex128{"
		for _, d := range ket.Data {
			outData = outData + strconv.FormatComplex(d, 'e', 8, 128) + ","
		}
		outData += "}"
		output[i] = fmt.Sprintf("{N: %d, Inc: %d, Data: %v}", ket.N, ket.Inc, outData)
	}
	return output
}
