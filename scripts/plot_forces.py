import matplotlib.pyplot as plt
import sys
import yaml

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
        xs = []
        ys = []
        for b in contents["values"]["system"]["bath"]:
            xs.append(b["angle"])
            ys.append(b["force"])
        plt.scatter(xs, ys)
        plt.title(metadata["simulation"] + "\n" + r"$\Psi(0) = $" + ket_to_arrows(contents["values"]["system"]["physicsconfig"]["spinevolutionconfig"]["initialket"]))

    p = filename.split("/")[-1]
    plt.savefig(f"{figdir}/force-plt-{p}.png")

if __name__ == '__main__':
    paths = sys.argv[1:]
    plot("outputs", "figures", paths[0])