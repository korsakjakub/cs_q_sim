import matplotlib.pyplot as plt
import sys
import yaml
import os

def ket_to_arrows(ket: str) -> str:
    c_spin = ket[0].replace('u', r'$\Uparrow$').replace('d', r'$\Downarrow$')
    bath = ket[1:].replace('u', r'$\uparrow$').replace('d', r'$\downarrow$')
    return "|" + c_spin + r"$\rangle$" + r"$\otimes$" + "|" + bath + r"$\rangle$"

def decay(contents, figdir, filenames):
    metadata = contents["metadata"]
    pc = contents['values']['system']['physicsconfig']
    for xys in contents["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(xy["x"])
            ys.append(xy["y"])
    plt.plot(xs, ys, label=f"Geometry: {pc['geometry']}\nN: {pc['bathcount']}")
    p = filename.split("/")[-1]
    p = os.path.splitext(p)[0]
    plt.title(metadata["simulation"] + "\n" + r"$\tau(\beta) = A^{-1}(\beta) = \frac{1}{\underset{k}{\max}\,|C_k|-\underset{k}{\min}\,|C_k|}$")
    plt.xlabel(r"$\beta / \pi$")
    plt.ylabel("sec")
    plt.yscale('log')
    plt.legend()
    plt.savefig(f"{figdir}/decay-couplings-plt-{p}.png")

def spread(contents, figdir, filename):
    metadata = contents["metadata"]
    pc = contents['values']['system']['physicsconfig']
    for xys in contents["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(xy["x"])
            ys.append(xy["y"])
    plt.plot(xs, ys, label=f"Geometry: {pc['geometry']}\nN: {pc['bathcount']}")
    p = filename.split("/")[-1]
    p = os.path.splitext(p)[0]
    plt.title(metadata["simulation"] + "\n" + r"$A(\beta) = \underset{k}{\max}\,|C_k|-\underset{k}{\min}\,|C_k|$")
    plt.xlabel(r"$\beta / \pi$")
    plt.ylabel(r"$2\pi \times \mathrm{kHz}$")
    plt.legend()
    plt.savefig(f"{figdir}/spread-couplings-plt-{p}.png")

def time_evolution(contents, figdir, filename):
    metadata = contents["metadata"]
    pc = contents['values']['system']['physicsconfig']
    i = 0
    for xys in contents["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(xy["x"])
            ys.append(xy["y"])
        slot = pc["observablesconfig"][i]["slot"]
        plt.plot(xs, ys, label=r"$\langle S_z^{(" + f"{slot}" + r")}\rangle(t)$" + f"\nbeta: {pc['tiltangle']}" + r"$\pi$" + f"\nGeometry: {pc['geometry']}\nmodel: {pc['model']}")
        i += 1
    plt.title(metadata["simulation"] + "\n" + r"$\Psi(0) = $" + ket_to_arrows(pc["initialket"]))

    plt.legend(loc='upper right')
    plt.xlabel("sec")
    plt.ylabel(r"$\langle S_z \rangle$")
    plt.ylim([-0.55, 0.55])
    p = filename.split("/")[-1]
    p = os.path.splitext(p)[0]
    plt.savefig(f"{figdir}/plt-{p}.png")

def interaction_strength(contents, figdir, filename):
    metadata = contents["metadata"]
    pc = contents['values']['system']['physicsconfig']
    for xys in contents["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(int(xy["x"]))
            ys.append(xy["y"])
    plt.scatter(xs, ys, label=f"Geometry: {pc['geometry']}\nN: {pc['bathcount']}\nbeta: {pc['tiltangle']}")
    p = filename.split("/")[-1]
    p = os.path.splitext(p)[0]
    plt.title(metadata["simulation"])
    plt.xlabel("Molecule number")
    plt.ylabel(r"$2\pi \times \mathrm{kHz}$")
    plt.legend()
    plt.savefig(f"{figdir}/interaction-strength-{pc['geometry']}-plt-{p}.png")


if __name__ == '__main__':
    paths = sys.argv[1:]
    filename = paths[0]
    outdir = os.path.dirname(filename)
    filename = os.path.basename(filename)
    with open(f"{outdir}/{filename}") as res_file:
        contents = yaml.safe_load(res_file)
        s = contents["metadata"]["simulationid"]
        figdir = contents["metadata"]["figuresdir"]
        match s:
            case "spin-evolution":
                time_evolution(contents, figdir, paths[0])
            case "spin-evolution-selected-coeffs":
                time_evolution(contents, figdir, paths[0])
            case "spread-of-couplings":
                spread(contents, figdir, paths[0])
            case "decay-time":
                decay(contents, figdir, paths[0])
            case "interactions":
                interaction_strength(contents, figdir, paths[0])

