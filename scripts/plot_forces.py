import matplotlib.pyplot as plt
import sys
import yaml
import numpy as np

def ket_to_arrows(ket: str) -> str:
    c_spin = ket[0].replace('u', r'$\Uparrow$').replace('d', r'$\Downarrow$')
    bath = ket[1:].replace('u', r'$\uparrow$').replace('d', r'$\downarrow$')
    return "|" + c_spin + r"$\rangle$" + r"$\otimes$" + "|" + bath + r"$\rangle$"

def plot(outdir, figdir, filename):
    if "outputs" in filename:
        outdir = filename.split("/")[0]
        filename = filename.split("/")[1]
    with open(f"{outdir}/{filename}") as res_file:
        contents = yaml.safe_load(res_file)
        metadata = contents["metadata"]
        ys = []
        for b in contents["values"]["system"]["bath"]:
            ys.append(b["force"])
        plt.scatter(np.arange(len(ys)), ys)
        plt.title(metadata["simulation"] + "\nfor: " + contents["values"]["system"]["physicsconfig"]["geometry"])

    p = filename.split("/")[-1]
    plt.savefig(f"{figdir}/force-plt-{p}.png")

if __name__ == '__main__':
    paths = sys.argv[1:]
    plot("outputs", "figures", paths[0])