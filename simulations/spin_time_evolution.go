package simulations

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/plotter"
)

func prepareObservables(conf qs.PhysicsConfig, downSpins int) []qs.Observable {
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
		fullObservable := qs.ManyBodyOperator(operator, obs.Slot, ketLength)
		if downSpins < 2 {
			observables[i] = qs.Observable{Dense: *fullObservable}
		} else {
			indices := qs.BasisIndices(conf.BathCount+1, downSpins)
			observables[i] = qs.Observable{Dense: *qs.RestrictMatrixToSubspace(fullObservable, indices)}
		}
	}
	return observables
}

func solveEigenProblem(s *qs.System) qs.Eigen {
	b := s.PhysicsConfig.BathMagneticField
	b0 := s.PhysicsConfig.CentralMagneticField
	if s.DownSpins < 1 {
		hamiltonian := s.Hamiltonian(b0, b)
		return s.Diagonalize(hamiltonian)
	}
	indices := qs.BasisIndices(s.PhysicsConfig.BathCount+1, s.DownSpins)
	hamiltonian := s.HamiltonianInBase(b0, b, indices)

	return s.Diagonalize(hamiltonian)
}

func prepareInitialKet(s *qs.System) *mat.VecDense {
	fullKet := mat.NewVecDense(int(math.Pow(2*s.PhysicsConfig.Spin+1, float64(len(s.PhysicsConfig.InitialKet)))), qs.ManyBodyVector(s.PhysicsConfig.InitialKet, int(2*s.PhysicsConfig.Spin+1)))
	if s.DownSpins < 1 {
		return fullKet
	}
	indices := qs.BasisIndices(s.PhysicsConfig.BathCount+1, s.DownSpins)
	ketData := make([]float64, len(indices))
	for i, index := range indices {
		ketData[i] = fullKet.AtVec(index)
	}
	return mat.NewVecDense(len(indices), ketData)
}

func downSpins(ket string) int {
	downSpins := 0
	for s := range ket {
		if ket[s] == 'd' {
			downSpins++
		}
	}
	return downSpins
}

func SpinTimeEvolution(conf qs.Config) {
	cs := qs.State{Angle: 0.0, Distance: 0.0}
	var bath []qs.State
	conf.Physics.BathCount = len(conf.Physics.InitialKet) - 1
	downSpins := downSpins(conf.Physics.InitialKet)
	timeRange := conf.Physics.TimeRange
	observables := prepareObservables(conf.Physics, downSpins)
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
		DownSpins:     downSpins,
	}
	initialKet := prepareInitialKet(s)

	if conf.Verbosity == "debug" && s.DownSpins > 0 {
		fmt.Printf("Reduced the dimension: %v -> %v\n\n", math.Pow(2.0, float64(conf.Physics.BathCount+1)), len(initialKet.RawVector().Data))
	}

	var eigen qs.Eigen
	if diagDir := conf.Files.DiagonalizationDir; diagDir != "" {
		fmt.Printf("\nLoading diagonalization results from %v\n", diagDir)
		eigen = qs.LoadDiagonalizationSolutions(diagDir)
	} else {
		fmt.Println("Diagonalizing...")
		eigen = solveEigenProblem(s)
		qs.SaveDiagonalizationSolutions(eigen, *s, conf.Files.OutputsDir+"diag-"+startTime)
	}

	if conf.Verbosity == "debug" {
		fmt.Println("Calculating the inner product matrix...")
	}

	gramMatrix := qs.Grammian(initialKet, eigen.EigenVectors)

	if conf.Verbosity == "debug" {
		fmt.Println("Calculating time evolution...")
	}

	type expVal struct {
		exp   float64
		index int
	}

	xyss := make([]plotter.XYs, len(observables))
	for i, observable := range observables {
		var xys plotter.XYs
		expValChannel := make(chan expVal, timeRange)
		var wg sync.WaitGroup

		for t := 0; t < timeRange; t++ {
			evolutionTime := conf.Physics.Dt * float64(t)
			wg.Add(1)
			go func(ch chan expVal, time int) {
				ch <- expVal{exp: observable.ExpectationValue(qs.Evolve(initialKet, evolutionTime, eigen.EigenValues, eigen.EigenVectors, gramMatrix)), index: time}
				if conf.Verbosity == "debug" {
					fmt.Printf("t= %v\n", evolutionTime)
				}
				wg.Done()
			}(expValChannel, t)
		}

		wg.Wait()
		expValues := make([]expVal, 0, timeRange)
		close(expValChannel)

		for e := range expValChannel {
			expValues = append(expValues, e)
		}

		sort.Slice(expValues, func(i, j int) bool {
			return expValues[i].index < expValues[j].index
		})

		for t := 0; t < timeRange; t++ {
			xys = append(xys, plotter.XY{X: float64(expValues[t].index) * conf.Physics.Dt / (2.0 * math.Pi), Y: expValues[t].exp})
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
			System qs.System "mapstructure:\"system\""
		}{
			System: *s,
		},
		XYs: xyss,
	}
	r.Write(conf.Files)
}
