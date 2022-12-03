# Central spin hamiltonian quantum simulator

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
