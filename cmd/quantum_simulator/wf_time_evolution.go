package main

import (
	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
)

func wf_time_evolution(conf qs.Config) {
	/*
		Given:
		- initial condition -> psi(t=0)
		- hamiltonian <- system

		Do:
		- calculate \ket{\psi(t)} = \sum_j e^(-i E_j t) \ket{E_j}\bra{E_j} \ket{\psi(0)}
		- plot ||psi(t)||
	*/
}
