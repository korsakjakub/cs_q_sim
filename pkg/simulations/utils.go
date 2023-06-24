package simulations

import (
	"math"
	
	cs "github.com/korsakjakub/cs_q_sim/pkg/cs_q_sim"
	"gonum.org/v1/gonum/mat"
)

type spectrumInput struct {
	hamiltonian   *mat.SymDense
	magneticField float64
}

type spectrumOutput struct {
	eigenValues   []float64
	eigenVectors  *mat.Dense
	magneticField float64
}

func spread(states []cs.State) float64 {
	max := math.Abs(states[0].InteractionStrength)
	min := math.Abs(states[0].InteractionStrength)
	for i := range states {
		strength := math.Abs(states[i].InteractionStrength)
		if strength > max {
			max = strength
		} else if strength < min {
			min = strength
		}
	}
	return max - min
}

func prepareStates(conf cs.Config) []cs.State {
	var bath []cs.State
	bc := conf.Physics.BathCount

	for i := 0; i < bc; i += 1 {
		bath = append(bath, cs.State{Angle: cs.PolarAngleCos(i, conf.Physics), Distance: conf.Physics.ConstantDistance})
	}

	s := &cs.System{
		CentralSpin:   cs.State{Angle: 0.0, Distance: 0.0},
		Bath:          bath,
		PhysicsConfig: conf.Physics,
	}

	for j := 0; j <= bc; j += 1 {
		s.InteractionAt(j)
	}
	return s.Bath
}

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
