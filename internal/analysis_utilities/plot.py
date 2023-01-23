import matplotlib.pyplot as plt
import sys
import yaml

def ket_to_arrows(ket: str) -> str:
    c_spin = ket[0].replace('u', r'$\Uparrow$').replace('d', r'$\Downarrow$')
    bath = ket[1:].replace('u', r'$\uparrow$').replace('d', r'$\downarrow$')
    return "|" + c_spin + r"$\rangle$" + r"$\otimes$" + "|" + bath + r"$\rangle$"

def plot(outdir, figdir, filenames):
    for filename in filenames:
        if "outputs" in filename:
            if len(filename.split("/")) > 2:
                outdir = filename.split("/")[0] + "/" + filename.split("/")[1]
                filename = filename.split("/")[2]
            else:
                outdir = filename.split("/")[0]
                filename = filename.split("/")[1]
        with open(f"{outdir}/{filename}") as res_file:
            contents = yaml.safe_load(res_file)
            metadata = contents["metadata"]
            xs = []
            ys = []
            for xy in contents["xys"]:
                xs.append(xy["x"])
                ys.append(xy["y"])
            plt.plot(xs, ys, label=r"$\langle S_z^{(" + filename[0] + r")}\rangle(t)$")
            plt.title(metadata["simulation"] + "\n" + r"$\Psi(0) = $" + ket_to_arrows(contents["config"]["spinevolutionconfig"]["initialket"]))

    plt.legend(loc='upper right')
    plt.ylim([-0.55, 0.55])
    p = filenames[0].split("/")[-1]
    plt.savefig(f"{figdir}/plt-{p}.png")

if __name__ == '__main__':
    paths = sys.argv[1:]
    plot("outputs", "figures", paths)