package quantum_simulator

type System struct {
	CentralSpin State
	Bath        []State
}

type State struct {
	Angle            float64
	Distance         float64
	MagneticField    float64
	InteractionForce float64
}

func (s *System) forces(cnf ...Config) []float64 {
	return []float64{0.0}
}
