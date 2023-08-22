import matplotlib.pyplot as plt
import numpy as np
import yaml
from matplotlib import rc

# Model XX
paths = [
    # "outputs/icosahedron/2023-08-06T12:27:52+02:00.yaml",
    # "outputs/icosahedron/2023-08-06T13:00:51+02:00.yaml",
    "outputs/icosahedron/2023-08-06T12:20:34+02:00.yaml",
    # "outputs/icosahedron/2023-08-06T13:01:17+02:00.yaml",
    "outputs/time-evolution-20230818/ring-0.4365-2023-08-18T23:37:11+02:00.yaml",
    "outputs/time-evolution-20230818/gauss-1.469-2023-08-19T00:03:18+02:00.yaml",
    ]

rc('font', **{'family': 'serif', 'serif': ['Computer Modern'], 'size': 11})
rc('text', usetex=True)
# Initialize a figure with 8 subplots
fig, ax = plt.subplots(1, 1, figsize=(4, 1.5), sharex=True, sharey=True)
markings = [r"$\beta = 0.18\,\pi$, Icosahedron", r"$\beta = 0.44\,\pi$, Ring", r"$B = 1.469$, Gauss"]

y_min, y_max = -0.5, 0.5
x_min, x_max = 0.0, 5e-6

# Iterate over the file paths and load the YAML data
for i, file_path in enumerate(paths):
    # Load YAML data
    with open(file_path, "r") as file:
        yaml_data = yaml.safe_load(file)

    # Get the subplot indices
    pc = yaml_data["values"]["system"]["physicsconfig"]

    for xys in yaml_data["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(xy["x"])
            ys.append(xy["y"])

    ax.set_ylim(y_min, y_max)
    ax.set_xlim(x_min, x_max)
    ax.set_yticks([-0.5, 0, 0.5])
    ax.set_ylabel(r"$\langle S_Z^{(0)}\rangle(t)$")
    ax.plot(xs, ys, label=markings[i])
    ax.set_xlabel(r"$t\,[\mathrm{seconds}]$")
    if i == 2:
        ax.axvline(x=0.35e-6, ls='dashed', label=r'$\tau = 0.35\,\mu\mathrm{s}$')
    #axes[row, col].legend(loc="upper right")

# axes[0, 0].title.set_text("Ring")
# axes[0, 1].title.set_text("Icosahedron")
# axes[0, 2].title.set_text("Cube")
# Adjust subplot spacing and layout
fig.legend(bbox_to_anchor=(1.5, 0.93), ncol=1, loc='upper right')
# fig.tight_layout()
# Show the plot
plt.savefig(f"figures/icosa-gauss-ring-decay-time-comparison.png", bbox_inches="tight", dpi=600)
