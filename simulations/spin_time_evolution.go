package simulations

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	cs "github.com/korsakjakub/cs_q_sim/pkg/cs_q_sim"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/plotter"
)

func prepareObservables(conf cs.PhysicsConfig, downSpins int) []cs.Observable {
	observables := make([]cs.Observable, len(conf.ObservablesConfig))
	ketLength := len(conf.InitialKet)
	for i, obs := range conf.ObservablesConfig {
		if obs.Slot > ketLength {
			continue
		}
		var operator *mat.Dense
		switch obs.Operator {
		case "Sz":
			operator = cs.Sz(conf.Spin)
		case "Sp":
			operator = cs.Sp(conf.Spin)
		case "Sm":
			operator = cs.Sm(conf.Spin)
		default:
			operator = cs.Id(conf.Spin)
		}
		fullObservable := cs.ManyBodyOperator(operator, obs.Slot, ketLength)
		if downSpins < 2 {
			observables[i] = cs.Observable{Dense: *fullObservable}
		} else {
			indices := cs.BasisIndices(conf.BathCount+1, downSpins)
			observables[i] = cs.Observable{Dense: *cs.RestrictMatrixToSubspace(fullObservable, indices)}
		}
	}
	return observables
}

func solveEigenProblem(s *cs.System) cs.Eigen {
	b := s.PhysicsConfig.BathMagneticField
	b0 := s.PhysicsConfig.CentralMagneticField
	if s.DownSpins < 1 {
		hamiltonian := s.Hamiltonian(b0, b)
		return s.Diagonalize(hamiltonian)
	}
	indices := cs.BasisIndices(s.PhysicsConfig.BathCount+1, s.DownSpins)
	hamiltonian := s.HamiltonianInBase(b0, b, indices)

	return s.Diagonalize(hamiltonian)
}

func prepareInitialKet(s *cs.System) *mat.VecDense {
	fullKet := mat.NewVecDense(int(math.Pow(2*s.PhysicsConfig.Spin+1, float64(len(s.PhysicsConfig.InitialKet)))), cs.ManyBodyVector(s.PhysicsConfig.InitialKet, int(2*s.PhysicsConfig.Spin+1)))
	if s.DownSpins < 1 {
		return fullKet
	}
	indices := cs.BasisIndices(s.PhysicsConfig.BathCount+1, s.DownSpins)
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

func SpinTimeEvolution(conf cs.Config) {
	var bath []cs.State
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
		bath = append(bath, cs.State{Angle: cs.PolarAngleCos(i, conf.Physics), Distance: conf.Physics.ConstantDistance})
	}
	s := &cs.System{
		CentralSpin:   cs.State{Angle: 0.0, Distance: 0.0},
		Bath:          bath,
		PhysicsConfig: conf.Physics,
		DownSpins:     downSpins,
	}
	initialKet := prepareInitialKet(s)

	if conf.Verbosity == "debug" && s.DownSpins > 0 {
		fmt.Printf("Reduced the dimension: %v -> %v\n\n", math.Pow(2.0, float64(conf.Physics.BathCount+1)), len(initialKet.RawVector().Data))
	}

	var eigen cs.Eigen
	if diagDir := conf.Files.DiagonalizationDir; diagDir != "" {
		fmt.Printf("\nLoading diagonalization results from %v\n", diagDir)
		eigen = cs.LoadDiagonalizationSolutions(diagDir)
	} else {
		fmt.Println("Diagonalizing...")
		eigen = solveEigenProblem(s)
		cs.SaveDiagonalizationSolutions(eigen, *s, conf.Files.OutputsDir+"diag-"+startTime)
	}

	if conf.Verbosity == "debug" {
		fmt.Println("Calculating the inner product matrix...")
	}

	gramMatrix := cs.Grammian(initialKet, eigen.EigenVectors)

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
				ch <- expVal{exp: observable.ExpectationValue(cs.Evolve(initialKet, evolutionTime, eigen.EigenValues, eigen.EigenVectors, gramMatrix)), index: time}
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
	r := cs.ResultsIO{
		Filename: startTime,
		Metadata: cs.Metadata{
			Date:           startTime,
			Simulation:     "Central spin expectation value time evolution",
			Cpu:            conf.Files.ResultsConfig.Cpu,
			Ram:            conf.Files.ResultsConfig.Ram,
			CompletionTime: elapsedTime.String(),
		},
		Values: struct {
			System cs.System "mapstructure:\"system\""
		}{
			System: *s,
		},
		XYs: xyss,
	}
	r.Write(conf.Files)
}
