package main

import (
	"fmt"
	"os"

	"github.com/korsakjakub/cs_q_sim/internal/quantum_simulator"
)

var conf quantum_simulator.Config

func main() {
	conf = quantum_simulator.LoadConfig([]string{os.Getenv("CONFIG_PATH")}, os.Getenv("CONFIG_NAME"), os.Getenv("CONFIG_TYPE"))
	s := quantum_simulator.System{}
	fmt.Println(s)
	fmt.Println(conf)
}
