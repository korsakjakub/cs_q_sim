import math
import matplotlib.pyplot as plt
import yaml
from matplotlib import rc

paths = [
    "outputs/interactions-vs-geometry/cube-0-2023-08-05T09:25:37+02:00.yaml",
    "outputs/interactions-vs-geometry/cube-0.03-2023-08-05T09:25:55+02:00.yaml",
    "outputs/interactions-vs-geometry/cube-0.05-2023-08-05T09:26:11+02:00.yaml",
    "outputs/interactions-vs-geometry/cube-0.4-2023-08-05T09:26:23+02:00.yaml",
    "outputs/interactions-vs-geometry/dode-0-2023-08-05T09:24:16+02:00.yaml",
    "outputs/interactions-vs-geometry/dode-0.03-2023-08-05T09:24:33+02:00.yaml",
    "outputs/interactions-vs-geometry/dode-0.05-2023-08-05T09:24:46+02:00.yaml",
    "outputs/interactions-vs-geometry/dode-0.4-2023-08-05T09:25:01+02:00.yaml",
    "outputs/interactions-vs-geometry/icosa-0-2023-08-05T09:22:33+02:00.yaml",
    "outputs/interactions-vs-geometry/icosa-0.03-2023-08-05T09:23:06+02:00.yaml",
    "outputs/interactions-vs-geometry/icosa-0.05-2023-08-05T09:23:19+02:00.yaml",
    "outputs/interactions-vs-geometry/icosa-0.4-2023-08-05T09:23:31+02:00.yaml",
    "outputs/interactions-vs-geometry/ring-0-2023-08-05T09:20:56+02:00.yaml",
    "outputs/interactions-vs-geometry/ring-0.03-2023-08-05T09:21:19+02:00.yaml",
    "outputs/interactions-vs-geometry/ring-0.05-2023-08-05T09:21:31+02:00.yaml",
    "outputs/interactions-vs-geometry/ring-0.4-2023-08-05T09:21:44+02:00.yaml"
]

rc('font', **{'family': 'serif', 'serif': ['Computer Modern'], 'size': 11})
rc('text', usetex=True)
geometry = ["cube", "dodecahedron", "icosahedron", "ring"]
markings = ["(a)", "(b)", "(c)", "(d)"]

# Initialize a figure with 4 subplots in a 2x2 grid
fig, axes = plt.subplots(2, 2, figsize=(6.8, 4))

# Iterate over the geometries and their respective file paths
for i, geo in enumerate(geometry):
    if geo == "dodecahedron":
        geo = "dode"
    if geo == "icosahedron":
        geo = "icosa"
    # Get the subplot indices
    row = i // 2
    col = i % 2
    ax = axes[row, col]

    # Iterate over the 4 files for each geometry
    for j in range(4):
        file_path = f"outputs/interactions-vs-geometry/{geo}-{j}-2023-08-05.yaml"

        # Load YAML data
        with open(file_path, "r") as file:
            yaml_data = yaml.safe_load(file)
        pc = yaml_data["values"]["system"]["physicsconfig"]

        for xys in yaml_data["xys"]:
            xs = [xy["x"] for xy in xys]
            ys = [xy["y"] for xy in xys]
        xs = [x + 1 for x in xs]
        ax.scatter(xs, ys, label=r"$\beta = " + f"{pc['tiltangle']}" + r"\,\pi$")
        ax.plot(xs, ys)

    #ax.set_title(f"{geometry[i].capitalize()}")
    xTicks = range(min(xs), math.ceil(max(xs))+2)
    xTicksUsed = [1] + [x for x in xTicks if x % 2 == 0]
    ax.set_xticks(xTicksUsed)
    # if col == 0:
    #     ax.set_ylabel(r"Interaction strength [$2\pi\,\times\,\mathrm{kHz}$]")
    if row == 1:
        ax.set_xlabel("Index $j$")
    lines_labels = [axes[0, 0].get_legend_handles_labels()]
    lines, labels = [sum(lol, []) for lol in zip(*lines_labels)]
    fig.legend(lines, labels, loc='upper center', ncol=4, fancybox=True)

axes[0, 0].text(0.02, 0.1, markings[0] + f" {geometry[0]}", transform=axes[0, 0].transAxes, fontsize=11, va='top')
axes[0, 1].text(0.02, 0.1, markings[1] + f" {geometry[1]}", transform=axes[0, 1].transAxes, fontsize=11, va='top')
axes[1, 0].text(0.02, 0.1, markings[2] + f" {geometry[2]}", transform=axes[1, 0].transAxes, fontsize=11, va='top')
axes[1, 1].text(0.02, 0.1, markings[3] + f" {geometry[3]}", transform=axes[1, 1].transAxes, fontsize=11, va='top')
fig.text(0.0, 0.5, r"Interaction strength $C_j$ [$2\pi\,\times\,\mathrm{kHz}$]", va='center', rotation='vertical')

# plt.tight_layout()
#plt.legend(title=r"$\lambda$ values", bbox_to_anchor=(1.05, 1), loc='upper left')

plt.savefig(f"figures/interactions-vs-geometries/interactions-vs-geometries.png", dpi=600)