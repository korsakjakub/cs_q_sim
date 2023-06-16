import matplotlib.pyplot as plt
import sys
import yaml
import numpy as np

def plot(outdir, figdir, filenames):
    for filename in filenames:
        if "outputs" in filename:
            outdir = filename.split("/")[0]
            filename = filename.split("/")[1]
        with open(f"{outdir}/{filename}") as res_file:
            contents = yaml.safe_load(res_file)
            metadata = contents["metadata"]
            for xys in contents["xys"]:
                xs = []
                ys = []
                for xy in xys:
                    xs.append(xy["x"])
                    ys.append(xy["y"])
            plt.plot(xs, ys, label=f"Geometry: {contents['values']['system']['physicsconfig']['geometry']}\nN: {contents['values']['system']['physicsconfig']['bathcount']}")
        p = filename.split("/")[-1]
    plt.title(metadata["simulation"] + "\n" + r"$\tau(\beta) = A^{-1}(\beta) = \frac{1}{\underset{k}{\max}\,|C_k|-\underset{k}{\min}\,|C_k|}$")
    plt.xlabel(r"$\beta / \pi$")
    plt.ylabel(r"$2\pi \times \mathrm{kHz}$")
    plt.legend()
    plt.savefig(f"{figdir}/decay-couplings-plt-{p}.png")

if __name__ == '__main__':
    paths = sys.argv[1:]
    plot("outputs", "figures", paths)
