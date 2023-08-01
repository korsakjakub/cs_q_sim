import matplotlib.pyplot as plt
import numpy as np
import yaml

#   Model XX vs. XXX for ring
paths = [
    "outputs/outputs-08072023/ring-0-2023-07-08T12:36:17+02:00.yaml",
    "outputs/outputs-XXX-08072023/ring-0-2023-07-08T14:37:27+02:00.yaml",
    "outputs/outputs-08072023/ring-1-2023-07-08T12:34:50+02:00.yaml",
    "outputs/outputs-XXX-08072023/ring-1-2023-07-08T14:37:35+02:00.yaml",
    "outputs/outputs-08072023/ring-2-2023-07-08T12:35:16+02:00.yaml",
    "outputs/outputs-XXX-08072023/ring-2-2023-07-08T14:37:47+02:00.yaml",
    "outputs/outputs-08072023/ring-3-2023-07-08T12:35:54+02:00.yaml",
    "outputs/outputs-XXX-08072023/ring-3-2023-07-08T14:38:01+02:00.yaml"
]

# Initialize a figure with 8 subplots
fig, axes = plt.subplots(4, 2, figsize=(8, 6), sharex=True)
plt.rc('font', size=10) 
markings = ["(a)", "(a')", "(b)", "(b')", "(c)", "(c')", "(d)", "(d')"]

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
    row = i//2
    col = i%2
    pc = yaml_data["values"]["system"]["physicsconfig"]

    for xys in yaml_data["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(xy["x"])
            ys.append(xy["y"])

    axes[row, col].text(0.02, 0.95, markings[i] + r", $\beta = $" + f"{pc['tiltangle']}" + r"$\pi$", transform=axes[row, col].transAxes, fontsize=10, fontweight='normal', va='top')
    if i%2==0:
        prev_ys = ys
        prev_model = pc['model']
        axes[row, col].set_yticks([-0.5, 0, 0.5])
        axes[row, col].set_ylim([-0.5, 0.5])
        axes[row, col].set_ylabel(r"$\langle S_Z^{(0)}\rangle(t)$")
        #axes[row, col].plot(xs, ys, color="xkcd:orange", label=pc['model'])
    if i%2==1:
        line, = axes[row, col-1].plot(xs, ys, label=pc['model'], color="xkcd:orange")
        prev_line, = axes[row, col-1].plot(xs, prev_ys, label=prev_model, color="#666699")
        axes[row, col-1].legend(handles=[prev_line, line], loc="upper right")
        axes[row, col-1].set_label(["1", "2"])
        axes[row, col].set_ylabel(r"$|\Delta|(t)$")
        axes[row, col].set_ylim([0, 0.5])
        axes[row, col].plot(xs, list(map(lambda a, b: abs(a-b), ys, prev_ys)), color="xkcd:orange")
        prev_ys = []
    if row == 3:
        axes[row, col].set_xlabel(r"$t\,[\mathrm{sec}]$")
    #axes[row, col].legend(loc="upper right")

# Adjust subplot spacing and layout
fig.tight_layout()

# Show the plot
plt.savefig(f"figures/figures-XXX-08072023/XX-vs-XXX-ring.png", dpi=300)
