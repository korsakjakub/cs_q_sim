import matplotlib.pyplot as plt
import numpy as np
import yaml

paths = [
    "outputs/outputs-icosa-reduction-08072023/2023-07-08T15:42:22+02:00.yaml",
    "outputs/outputs-icosa-reduction-08072023/2023-07-19T10:02:00+02:00.yaml"
    ]

# Initialize a figure with 8 subplots
fig, axes = plt.subplots(1, 2, figsize=(8, 2), sharex=True)
plt.rc('font', size=10) 
markings = ["(a)" + r", $\beta = 0.3\pi$", "(b)"]

y_min, y_max = -0.5, 0.5
x_min, x_max = 0.0, 5e-5
prev_ys = []
prev_model = ""

# Iterate over the file paths and load the YAML data
for i, file_path in enumerate(paths):
    # Load YAML data
    with open(file_path, "r") as file:
        yaml_data = yaml.safe_load(file)

    # Get the subplot indices
    #row = i // 2
    pc = yaml_data["values"]["system"]["physicsconfig"]

    for xys in yaml_data["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(xy["x"])
            ys.append(xy["y"])

    axes[i].text(0.02, 0.95, markings[i], transform=axes[i].transAxes, fontsize=10, fontweight='normal', va='top')
    axes[i].set_ylim(y_min, y_max)
    axes[i].set_xlim(x_min, x_max)
    if i == 0:
        prev_ys = ys
        prev_bc = pc['bathcount']
        axes[i].set_yticks([-0.5, 0, 0.5])
        axes[i].set_ylabel(r"$\langle S_Z^{(0)}\rangle(t)$")
    else:
        line, = axes[i-1].plot(xs, ys, color="#666699", label=r"$N_b = $" + f"{prev_bc}")
        prev_line, = axes[i-1].plot(xs, prev_ys, color="xkcd:orange", label=r"$N_b = $" + f"{pc['bathcount']}")
        axes[i-1].legend(handles=[prev_line, line], loc="upper right")
        axes[i].set_ylabel(r"$|\Delta|(t)$")
        axes[i].set_ylim([0, 0.5])
        axes[i].plot(xs, list(map(lambda a, b: abs(a-b), ys, prev_ys)), color="xkcd:orange")
        prev_ys = []
        axes[i].set_xlabel(r"$t\,[\mathrm{sec}]$")

# Adjust subplot spacing and layout
fig.tight_layout()

# Show the plot
plt.savefig(f"figures/figures-icosa-reduction-08072023/bath-12-vs-bath-10-icosahedron.png", dpi=300)
