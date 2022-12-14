# Central spin hamiltonian quantum simulator

[![codecov](https://codecov.io/gh/korsakjakub/cs_q_sim/branch/main/graph/badge.svg?token=ARL78PQTW9)](https://codecov.io/gh/korsakjakub/cs_q_sim)

The system is governed by the $H_{XX}$ hamiltonian:

$$H_{XX} = \sum_{j=1}^{N_b} H_S(0, j) + H_B$$

where

$$H_S(0, j) = C_j [S_+^{(0)} S_-^{j} + S_-^{(0)} S_+^{(j)}]$$

$$H_B = B^{(0)} S_Z^{(0)} + B \sum_{j = 1}^{N_b} S_Z^{(j)} = (B^{(0)} - B)S_Z^{(0)}$$


## Simulated system structure
- $1$ central spin Rydberg state atom (CS)
  * in the center of the system
  * interacts with all bath molecules
- $N_b$ bath polarized molecules (BPM)
  * uniformly distributed around CS with distance $|\mathbf{R}_j| = R$ and angle $\theta_j = 2 \pi j/N_b$

## Energy scales
- $C_j \sim 100$ kHz
- $B^{(0)}, B \sim 1$ GHz
- $B{(0)} - B < 100$ kHz
- lifetime of Rydberg state $\sim 100$ $\mu\text{s}$

## Dependence of $C_j$ force on angle

$$ C_j = \frac{\mu_{\text{mol}} \mu_{\text{atom}}}{4 \pi \varepsilon_0 |\mathbf{R}_j|^3} \cdot \frac{1 - 3 \cos^2 \theta_j}{2}$$

Here $\theta_j$ is an angle between $\mathbf{R}_j$ and $Z$ axis

### Units
- $[\mu] = \text{Debye}$
- $[R] = \mu m$
- $[C] = h \cdot \text{Hz}$

### notes on anticipated Go's speed:
- [spzala19/Multiprocessing-with-golang](https://github.com/spzala19/Multiprocessing-with-golang)
- [comparison for web-like tasks](https://djangostars.com/blog/my-story-with-golang/)
- [general overview](https://www.stxnext.com/blog/go-go-python-rangers-comparing-python-and-golang/)
- [tips on optimizing diagonalization speed to match numpy](https://github.com/gonum/gonum/issues/511)
- [different strategies on utilizing CPU with Go](https://liamhieuvu.com/utilize-cpu-with-golang)
