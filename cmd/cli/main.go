package main

import (
	"os"

	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
	sim "github.com/korsakjakub/cs_q_sim/simulations"
)

var conf qs.Config

func main() {
	conf = qs.LoadConfig([]string{os.Getenv("CONFIG_PATH"), "../../config/", "config/"}, os.Getenv("CONFIG_NAME"), os.Getenv("CONFIG_TYPE"))
	sim.Spin_time_evolution(conf)
}
