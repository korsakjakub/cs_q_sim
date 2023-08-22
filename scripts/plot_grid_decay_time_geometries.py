import math
import matplotlib.pylab as plt
from matplotlib import rc
import numpy as np
import yaml

paths = [
    "outputs/spread-09072023/ring-2023-07-09T11:54:59+02:00.yaml",
    "outputs/spread-09072023/icosa-2023-07-09T11:55:09+02:00.yaml",
    "outputs/gaussian-01082023/spread-gaussian.yaml",
    ]

geometry = ["ring", "icosahedron", "gauss"]
markings = ["(a)", "(b)", "(c)"]

rc('font', **{'family': 'serif', 'serif': ['Computer Modern'], 'size': 11})
rc('text', usetex=True)
# Initialize a figure with 8 subplots
fig, axes = plt.subplots(1, 3, figsize=(6.8, 2.4))


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
            y = math.inf if abs(xy["y"]) < 1e-10 else xy["y"]* 2 * math.pi * 1e3
            ys.append(1.0/y)

    axes[i].text(0.07, 0.95, f"{markings[i]} {geometry[i]}", transform=axes[i].transAxes, fontsize=11, va='top')
    axes[i].plot(xs, ys, color="#666699")
    if i == 0:
        axes[i].set_ylabel(r"Decay timescale $\tau$ [seconds]")

axes[2].set_ylim([1e-7, 1])
axes[1].set_ylim([2e-7, 8.5e-7])
axes[2].set_xlabel(r"$B$")
axes[1].set_xlabel(r"$\beta / \pi$")
axes[0].set_xlabel(r"$\beta / \pi$")
    
axes[1].axhline(y=3.5e-7, color='xkcd:orange', ls='dashed')
axes[0].set_yscale('log')
axes[2].set_yscale('log')

# Adjust subplot spacing and layout
fig.tight_layout()

# Show the plot
plt.savefig(f"figures/spread-09072023/decay-time-comparison.png", dpi=600)
