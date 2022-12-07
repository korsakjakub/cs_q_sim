# plan podboju kosmosu

- system.go
  * recipe for $H_{XX}$ using `spin.go` and `many_body.go` - done
  * diagonalization using `mat.Eigen`
    * usage like: `eig, eigv := sys.diagonalize(params)`
    * procedure shall be concurrent - channel
- `*_plot.go`
  * functions for plotting interesting quantities
- `cmd/quantum_simulator/*.go`
  * in different files sweeps over different parameters so that calculations can be put in parallel 
