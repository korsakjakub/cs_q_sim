import matplotlib.pyplot as plt
from matplotlib import rc
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

rc('font', **{'family': 'serif', 'serif': ['Computer Modern'], 'size': 11})
rc('text', usetex=True)

fig, axes = plt.subplots(2, 2, figsize=(6.8, 3), sharex=True, sharey=True)
markings = ["(a)", "(b)", "(c)", "(d)"]

y_min, y_max = -0.5, 0.5
x_min, x_max = 0.0, 5e-5
prev_ys = []
prev_model = ""

for i, file_path in enumerate(paths):
    with open(file_path, "r") as file:
        yaml_data = yaml.safe_load(file)

    row = i//2
    col = i%2
    if row >= 2:
        row -= 2
        col += 1
    
    pc = yaml_data["values"]["system"]["physicsconfig"]

    for xys in yaml_data["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(xy["x"])
            ys.append(xy["y"])

    if i%2==0:
        prev_ys = ys
        prev_model = pc['model']
        axes[row, col].set_yticks([-0.5, 0, 0.5])
        axes[row, col].set_ylim([-0.5, 0.5])
    if i%2==1:
        line, = axes[row, col-1].plot(xs, ys, label=pc['model'], color="xkcd:orange")
        prev_line, = axes[row, col-1].plot(xs, prev_ys, label=prev_model, color="#666699")
        prev_ys = []

axes[0, 0].text(0.02, 0.98, markings[0] + r" $\beta = 0\pi$", transform=axes[0, 0].transAxes, fontsize=11, va='top')
axes[1, 0].text(0.02, 0.98, markings[1] + r" $\beta = 0.03\pi$", transform=axes[1, 0].transAxes, fontsize=11, va='top')
axes[0, 1].text(0.02, 0.98, markings[2] + r" $\beta = 0.05\pi$", transform=axes[0, 1].transAxes, fontsize=11, va='top')
axes[1, 1].text(0.02, 0.98, markings[3] + r" $\beta = 0.4\pi$", transform=axes[1, 1].transAxes, fontsize=11, va='top')
axes[0, 0].set_ylabel(r"$\langle S_Z^{(0)}\rangle(t)$")
axes[1, 0].set_ylabel(r"$\langle S_Z^{(0)}\rangle(t)$")
axes[1, 0].set_xlabel(r"$t\,[\mathrm{sec}]$")
axes[1, 1].set_xlabel(r"$t\,[\mathrm{sec}]$")

lines_labels = [axes[0, 0].get_legend_handles_labels()]
lines, labels = [sum(lol, []) for lol in zip(*lines_labels)]
fig.legend(lines, labels, loc='lower left', ncol=2)
# plt.subplots_adjust(right=1.85)
fig.tight_layout()

plt.savefig(f"figures/figures-XXX-08072023/XX-vs-XXX-ring.png", dpi=600)
