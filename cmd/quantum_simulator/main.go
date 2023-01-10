package main

import (
	"os"

	qs "github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
)

var conf qs.Config

func main() {
	conf = qs.LoadConfig([]string{os.Getenv("CONFIG_PATH")}, os.Getenv("CONFIG_NAME"), os.Getenv("CONFIG_TYPE"))
	// conf = qs.LoadConfig([]string{"../../config/"})
	spectrum(conf)
	spin_time_evolution(conf)
}
