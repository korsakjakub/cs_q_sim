import numpy as np
import yaml

path = 'outputs/gaussian-01082023/spread-gaussian.yaml'
Nb = 12
x1 = 284225.5096624722

def k(j, B):
    return np.exp(-(j*B/Nb)**2)

def A(j, B):
    if j == 0:
        return 0.0
    return x1 * Nb * k(j, B) / np.sum([k(j, B) for j in range(1, Nb+1)])

if __name__ == "__main__":
    xy = []
    for B in np.arange(0.0, 1.5, 1e-3):
        coeffs = [A(j, B) for j in range(1, Nb+1)]
        aMax = max([abs(c) for c in coeffs])
        aMin = min([abs(c) for c in coeffs])
        xy.append({'x': float(B), 'y': float(aMax-aMin)*1e-3})
    res = {
        'filename': 'spread-gaussian.yaml',
        'metadata': {
            'simulationid': 'spread-of-couplings',
            'simulation': 'spread-of-couplings',
            'figuresdir': 'figures/gaussian-01082023'
        },
        'values': {
            'system' : {
                'physicsconfig' : {
                    'bathcount': Nb,
                    'geometry': 'gauss',
                }
            }
        },
        'xys': [xy]
    }
    with open(path, 'w') as file:
        yaml.dump(res, file)