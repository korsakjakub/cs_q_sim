import matplotlib.pylab as plt
from matplotlib import rc
import numpy as np
import yaml

paths = [
    "outputs/spread-09072023/cube-2023-07-09T11:55:39+02:00.yaml",
    "outputs/spread-09072023/ring-2023-07-09T11:54:59+02:00.yaml",
    "outputs/spread-09072023/icosa-2023-07-09T11:55:09+02:00.yaml",
    "outputs/gaussian-01082023/spread-gaussian.yaml",
    ]

geometry = ["cube", "ring", "icosahedron", "gauss"]
markings = ["(a)", "(b)", "(c)", "(d)"]

rc('font', **{'family': 'serif', 'serif': ['Computer Modern'], 'size': 11})
rc('text', usetex=True)
fig, axes = plt.subplots(2, 2, figsize=(6.8, 4.25))

for i, file_path in enumerate(paths):
    with open(file_path, "r") as file:
        yaml_data = yaml.safe_load(file)

    row = i // 2
    col = i % 2
    pc = yaml_data["values"]["system"]["physicsconfig"]

    for xys in yaml_data["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(xy["x"])
            ys.append(xy["y"])

    axes[row, col].text(0.02, 0.95, f"{markings[i]} {geometry[i]}", transform=axes[row, col].transAxes, fontsize=11, va='top')
    axes[row, col].plot(xs, ys, color="#666699")

    # axes[1, 0].set_ylabel(r"Spread of couplings $A$ $[2\pi\,\times\,\mathrm{kHz}]$")
    # y_label = axes[1, 0].get_yaxis().get_label()
    # y_label.set_label_coords(-0.15, 0.5)  # Adjust the coordinates to move the label up

    fig.text(0.0, 0.5, r"Spread of couplings $A$ $[2\pi\,\times\,\mathrm{kHz}]$", va='center', rotation='vertical')

    if col == 1:
        axes[row, col].set_ylim(bottom=0)

    if col == 0 and row == 1:
        axes[row, col].set_ylim(top=650)

    axes[0, 0].set_xlabel(r"$\beta / \pi$")
    axes[1, 0].set_xlabel(r"$\beta / \pi$")
    axes[0, 1].set_xlabel(r"$\beta / \pi$")
    axes[1, 1].set_xlabel(r"$B$")

fig.tight_layout()

plt.savefig(f"figures/spread-09072023/spread-comparison.png", dpi=600)
