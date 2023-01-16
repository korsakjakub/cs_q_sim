import matplotlib.pyplot as plt
import numpy as np
import csv
import os
import sys

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
    initial_ket = "|" + metadata[-1][-1].split()[-1][:-1].replace('u', '↑').replace('d', '↓') + ">"
    p = plt.plot(xs, ys, label=f"Ψ(0) = {initial_ket}")
    plt.legend(loc='upper right')
    plt.savefig(f"{figdir}/plt-{filename}.png")

if __name__ == '__main__':
    filename = sys.argv[1]
    plot("outputs", "figures", filename)