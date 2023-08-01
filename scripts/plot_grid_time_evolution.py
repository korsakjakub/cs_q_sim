import matplotlib.pyplot as plt
import numpy as np
import yaml
from matplotlib import rc

# Model XX
paths = [
    "outputs/outputs-08072023/ring-0-2023-07-08T12:36:17+02:00.yaml",
    "outputs/outputs-08072023/ico-0-2023-07-08T12:48:23+02:00.yaml",
    "outputs/cube/cube-0-2023-07-19T09:23:55+02:00.yaml",
    "outputs/outputs-08072023/ring-1-2023-07-08T12:34:50+02:00.yaml",
    "outputs/outputs-08072023/ico-1-2023-07-08T12:48:45+02:00.yaml",
    "outputs/cube/cube-1-2023-07-19T09:24:20+02:00.yaml",
    "outputs/outputs-08072023/ring-2-2023-07-08T12:35:16+02:00.yaml",
    "outputs/outputs-08072023/ico-2-2023-07-08T12:48:38+02:00.yaml",
    "outputs/cube/cube-2-2023-07-19T09:25:49+02:00.yaml",
    "outputs/outputs-08072023/ring-3-2023-07-08T12:35:54+02:00.yaml",
    "outputs/outputs-08072023/ico-3-2023-07-08T12:48:30+02:00.yaml",
    "outputs/cube/cube-3-2023-07-19T09:26:08+02:00.yaml"
    ]

rc('font', **{'family': 'serif', 'serif': ['Computer Modern'], 'size': 11})
rc('text', usetex=True)
# Initialize a figure with 8 subplots
fig, axes = plt.subplots(4, 3, figsize=(6.8, 6.8), sharex=True, sharey=True)
plt.rc('font', size=10) 
markings = ["(a)", "(a')", "(a'')", "(b)", "(b')", "(b'')", "(c)", "(c')", "(c'')", "(d)", "(d')", "(d'')"]

y_min, y_max = -0.5, 0.5
x_min, x_max = 0.0, 5e-5

# Iterate over the file paths and load the YAML data
for i, file_path in enumerate(paths):
    # Load YAML data
    with open(file_path, "r") as file:
        yaml_data = yaml.safe_load(file)

    # Get the subplot indices
    row = i // 3
    col = i % 3
    pc = yaml_data["values"]["system"]["physicsconfig"]

    for xys in yaml_data["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(xy["x"])
            ys.append(xy["y"])

    axes[row, col].text(0.02, 0.98, markings[i] + r" $\beta = $" + f" {pc['tiltangle']}" + r"$\pi$", transform=axes[row, col].transAxes, fontsize=11, va='top')
    axes[row, col].set_ylim(y_min, y_max)
    axes[row, col].set_xlim(x_min, x_max)
    if col == 0:
        axes[row, col].set_yticks([-0.5, 0, 0.5])
        axes[row, col].set_ylabel(r"$\langle S_Z^{(0)}\rangle(t)$")
        axes[row, col].plot(xs, ys, color="#666699", label=r"$\beta$ = " + f"{pc['tiltangle']}" + r"$\pi$")
    else:
        axes[row, col].plot(xs, ys, color="#666699", label=r"$\beta$ = " + f"{pc['tiltangle']}" + r"$\pi$")

    if row == 3:
        axes[row, col].set_xlabel(r"$t\,[\mathrm{sec}]$")
    #axes[row, col].legend(loc="upper right")

axes[0, 0].title.set_text("Ring")
axes[0, 1].title.set_text("Icosahedron")
axes[0, 2].title.set_text("Cube")
# Adjust subplot spacing and layout
fig.tight_layout()

# Show the plot
plt.savefig(f"figures/ring-vs-icosa-vs-cube.png", dpi=600)
