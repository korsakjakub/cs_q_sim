import matplotlib.pyplot as plt
import numpy as np
import yaml

paths = [
    "outputs/spread-09072023/cube-2023-07-09T11:55:39+02:00.yaml",
    "outputs/spread-09072023/dode-2023-07-09T11:55:50+02:00.yaml",
    "outputs/spread-09072023/icosa-2023-07-09T11:55:09+02:00.yaml",
    "outputs/spread-09072023/ring-2023-07-09T11:54:59+02:00.yaml",
    ]

geometry = ["cube", "dodecahedron", "icosahedron", "ring"]
markings = ["(a)", "(b)", "(c)", "(d)"]

# Initialize a figure with 8 subplots
fig, axes = plt.subplots(2, 2, figsize=(8, 5), sharex=True)
plt.rc('font', size=10) 

x_min, x_max = 0.0, 0.5

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
    axes[row, col].text(0.02, 0.95, markings[i], transform=axes[row, col].transAxes, fontsize=10, fontweight='normal', va='top')
    axes[row, col].set_xlim(x_min, x_max)
    axes[row, col].plot(xs, ys, color="xkcd:greyish blue", label=f"{geometry[i]}")
    if col == 0:
        axes[row, col].set_ylabel(r"$2\pi\,\times\,\mathrm{kHz}$")

    if row == 1:
        axes[row, col].set_xlabel(r"$\beta / \pi$")
    axes[row, col].legend(loc="lower right")

# Adjust subplot spacing and layout
fig.tight_layout()

# Show the plot
plt.savefig(f"figures/spread-09072023/spread-comparison.png", dpi=300)
