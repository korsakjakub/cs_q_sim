# plan podboju kosmosu

## 7.12.2022
- system.go
  * recipe for $H_{XX}$ using `spin.go` and `many_body.go` - **done**
  * diagonalization using `mat.Eigen`
    * usage like: `eig, eigv := sys.diagonalize(params)` - **done**
    * procedure shall be concurrent - channel - **done**
- `*_plot.go`
  * functions for plotting interesting quantities - **1 plot scheme so far**
- `cmd/quantum_simulator/*.go`
  * in different files sweeps over different parameters so that calculations can be put in parallel 
    - **1 simulation so far** (could be more efficient for sure)
## 13.12.2022
- `results_io.go`
  * save calculations outputs to csv files 
  * read from csv into ResultsIO struct for easy plotting
