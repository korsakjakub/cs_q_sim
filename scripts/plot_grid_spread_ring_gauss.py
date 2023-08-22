import matplotlib.pylab as plt
from matplotlib import rc
import numpy as np
import yaml

paths = [
    "outputs/spread-09072023/ring-2023-07-09T11:54:59+02:00.yaml",
    "outputs/gaussian-01082023/spread-gaussian.yaml",
    ]

geometry = ["ring", "gauss"]
markings = ["(a)", "(b)"]

rc('font', **{'family': 'serif', 'serif': ['Computer Modern'], 'size': 11})
rc('text', usetex=True)
# Initialize a figure with 8 subplots
fig, axes = plt.subplots(1, 2, figsize=(6.8, 2.2), sharey=True)

# Iterate over the file paths and load the YAML data
for i, file_path in enumerate(paths):
    # Load YAML data
    with open(file_path, "r") as file:
        yaml_data = yaml.safe_load(file)

    # Get the subplot indices
    row = i // 2
    col = i % 2
    pc = yaml_data["values"]["system"]["physicsconfig"]

    for xys in yaml_data["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(xy["x"])
            ys.append(xy["y"])

    # axes[row, col].set_ylim(y_min, y_max)
    axes[col].text(0.02, 0.95, f"{markings[i]} {geometry[i]}", transform=axes[col].transAxes, fontsize=11, va='top')
    axes[0].set_xlim(0, 0.5)
    
    axes[col].plot(xs, ys, color="#666699")
    if col == 0:
        axes[col].set_ylabel(r"$2\pi\,\times\,\mathrm{kHz}$")
    else:
        axes[col].set_ylim(bottom=0)

    if col == 0 and row == 1:
        axes[col].set_ylim([0, 650])

    axes[0].set_xlabel(r"$\beta / \pi$")
    axes[1].set_xlabel(r"$B$")
    #axes[row, col].legend(loc="lower right")

# Adjust subplot spacing and layout
fig.tight_layout()

# Show the plot
plt.savefig(f"figures/gaussian-01082023/spread-ring-gauss.png", dpi=600)
