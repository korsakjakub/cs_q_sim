import matplotlib.pyplot as plt
import numpy as np
import csv
import os

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
    print(xs, ys)
    print(metadata)
    p = plt.plot(xs, ys)
    plt.savefig(f"{figdir}/plt-{filename}.png")

if __name__ == '__main__':
    plot("outputs", "figures", "2023-01-14T12:52:24+01:00")