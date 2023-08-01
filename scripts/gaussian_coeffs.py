import numpy as np

Nb = 12
B = 1
x1 = 284225.5096624722

def k(j):
    return np.exp(-(j*B/Nb)**2)

def A(j):
    if j == 0:
        return 0.0
    return x1 * Nb * k(j) / np.sum([k(j) for j in range(1, Nb+1)])

print([A(j) for j in range(Nb+1)])