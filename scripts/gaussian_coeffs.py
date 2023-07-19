import numpy as np

Nb = 11
B = 0.4
x1 = 2.0

def k(j):
    return np.exp(-(j*B/Nb)**2)

def A(j):
    if j == 0:
        return 0.0
    return x1 * Nb * k(j) / np.sum([k(j) for j in range(Nb)])

print([A(j) for j in range(Nb+1)])