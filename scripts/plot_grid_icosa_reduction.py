import matplotlib.pyplot as plt
from matplotlib import rc
import numpy as np
import yaml

paths = [
    "outputs/outputs-icosa-reduction-08072023/2023-07-08T15:42:22+02:00.yaml",
    "outputs/outputs-icosa-reduction-08072023/bath-10-new-2023-07-22T22:21:29+02:00.yaml"
    ]

rc('font', **{'family': 'serif', 'serif': ['Computer Modern'], 'size': 11})
rc('text', usetex=True)
# Initialize a figure with 8 subplots
fig, axes = plt.subplots(1, 2, figsize=(6.8, 1.7), sharex=True)
#plt.rc('font', size=10) 
markings = [r"(a) $\beta = 0.3\pi$", "(b)"]

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

    axes[i].text(0.02, 0.95, markings[i], transform=axes[i].transAxes, fontsize=11, va='top')
    axes[i].set_ylim(y_min, y_max)
    axes[i].set_xlim(x_min, x_max)
    if i == 0:
        prev_ys = ys
        prev_bc = pc['bathcount']
        axes[i].set_yticks([-0.5, 0, 0.5])
        axes[i].set_ylabel(r"$\langle S_Z^{(0)}\rangle(t)$")
    else:
        prev_line, = axes[i-1].plot(xs, prev_ys, color="xkcd:orange", label=r"$N_b = $" + f"{prev_bc}")
        line, = axes[i-1].plot(xs, ys, color="#666699", label=r"$N_b = $" + f"{pc['bathcount']}", ls='dashed')
        axes[i-1].legend(handles=[prev_line, line], loc="upper right")
        axes[i].set_ylabel(r"$|\Delta|(t)$")
        axes[i].set_ylim(-1e-4, 5e-2)
        axes[i].set_yticks([0.0, 2e-2, 5e-2])
        axes[i].plot(xs, list(map(lambda a, b: abs(a-b), ys, prev_ys)), color="xkcd:orange")
        prev_ys = []
    axes[i].set_xlabel(r"$t\,[\mathrm{seconds}]$")

fig.tight_layout()

plt.savefig(f"figures/figures-icosa-reduction-08072023/bath-12-vs-bath-10-icosahedron.png", dpi=600)
