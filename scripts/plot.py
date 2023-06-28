import matplotlib.pyplot as plt
import sys
import yaml

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
    plt.title(metadata["simulation"] + "\n" + r"$\tau(\beta) = A^{-1}(\beta) = \frac{1}{\underset{k}{\max}\,|C_k|-\underset{k}{\min}\,|C_k|}$")
    plt.xlabel(r"$\beta / \pi$")
    plt.ylabel(r"$2\pi \times \mathrm{kHz}$")
    plt.legend()
    plt.savefig(f"{figdir}/decay-couplings-plt-{p}.png")

def spread(contents, figdir, filenames):
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
        plt.plot(xs, ys, label=r"$\langle S_z^{(" + f"{slot}" + r")}\rangle(t)$")
        i += 1
    plt.title(metadata["simulation"] + "\n" + r"$\Psi(0) = $" + ket_to_arrows(pc["initialket"]))

    plt.legend(loc='upper right')
    plt.ylim([-0.55, 0.55])
    p = filename.split("/")[-1]
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
    plt.title(metadata["simulation"])
    plt.xlabel("Molecule number")
    plt.ylabel(r"$2\pi \times \mathrm{kHz}$")
    plt.legend()
    plt.savefig(f"{figdir}/interaction-strength-{pc['geometry']}-plt-{p}.png")


if __name__ == '__main__':
    paths = sys.argv[1:]
    filename = paths[0]
    if "outputs" in filename:
        outdir = filename.split("/")[0]
        filename = filename.split("/")[1]
    with open(f"{outdir}/{filename}") as res_file:
        contents = yaml.safe_load(res_file)
        s = contents["metadata"]["simulationid"]
        match s:
            case "spin-evolution":
                time_evolution(contents, "figures", paths[0])
            case "spread-of-couplings":
                spread(contents, "figures", paths[0])
            case "decay-time":
                decay(contents, "figures", paths[0])
            case "interactions":
                interaction_strength(contents, "figures", paths[0])