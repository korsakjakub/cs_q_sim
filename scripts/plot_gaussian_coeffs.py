import numpy as np
import matplotlib.pyplot as plt
from matplotlib import rc

path = 'figures/interactions-vs-geometries/interactions-gaussian.png'
Nb = 12
x1 = 284225.5096624722
rc('font', **{'family': 'serif', 'serif': ['Computer Modern'], 'size': 11})
rc('text', usetex=True)

def k(j, B):
    return np.exp(-(j*B/Nb)**2)

def A(j, B):
    if j == 0:
        return 0.0
    return 1e-3* x1 * Nb * k(j, B) / np.sum([k(j, B) for j in range(1, Nb+1)])

if __name__ == "__main__":
    fig, ax = plt.subplots(1, 1, figsize=(4.5,3))
    for B in [0.0, 0.165, 0.273, 1.337]:
        ys = [A(j, B) for j in range(1, Nb+1)]
        xs = [int(i) for i in np.arange(1, 13, 1)]
        ax.scatter(xs, ys, label=r"$B = $" + f" {B}")
        ax.plot(xs, ys)
    ax.set_ylabel(r"Interaction strength $C_j$ [$2\pi\,\times\,\mathrm{kHz}$]")
    ax.set_xlabel(r"Index $j$")
    ax.set_xticks(range(1, 13))
    fig.tight_layout()
    fig.legend(bbox_to_anchor=(0.97, 0.95), loc='upper right')
    plt.savefig(path, dpi=600)
