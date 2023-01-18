import matplotlib.pyplot as plt
import numpy as np
import csv
import os
import sys

def ket_to_arrows(ket: str) -> str:
    ket_string = ket[-1][-1].split()[-1][:-1]
    c_spin = ket_string[0].replace('u', r'$\Uparrow$').replace('d', r'$\Downarrow$')
    bath = ket_string[1:].replace('u', r'$\uparrow$').replace('d', r'$\downarrow$')
    return "|" + c_spin + r"$\rangle$" + r"$\otimes$" + "|" + bath + r"$\rangle$"

def plot(outdir, figdir, filename):
    with open(f"{outdir}/{filename}") as csv_file:
        metadata = []
        xs = []
        ys = []
        csv_read=csv.reader(csv_file, delimiter=',')
        i = 0
        for row in csv_read:
            if i > 1:
                xs.append(float(row[0]))
                ys.append(float(row[1]))
            else:
                metadata.append(row)
            i += 1
    initial_ket = ket_to_arrows(metadata)
    print(metadata)
    plt.plot(xs, ys, label=r"$\langle S_z^{(0)}\rangle(t)$")
    plt.title(metadata[0][1] + "\n" + r"$\Psi(0) = $" + ket_to_arrows(metadata))
    plt.legend(loc='upper right')
    plt.ylim([-0.55, 0.55])
    plt.savefig(f"{figdir}/plt-{filename}.png")

if __name__ == '__main__':
    filename = sys.argv[1]
    plot("outputs", "figures", filename)