import matplotlib.pyplot as plt
import numpy as np
from numpy import fft
import yaml
from matplotlib import rc

def spectrogram(xs, ys):
    window_size = 256
    overlap = 128
    fs = 1/(xs[1]-xs[0])
    window = np.hanning(window_size)
    windows = [ys[i:i+window_size] * window for i in range(0, len(ys)-window_size, window_size-overlap)]
    spectrogram = [np.abs(np.fft.rfft(win))**2 for win in windows]
    spectrogram = np.array(spectrogram).T

    frequencies = np.fft.rfftfreq(window_size, d=1.0/fs)
    time = np.arange(len(spectrogram[0])) * (window_size - overlap) / fs
    return time, frequencies, 10 * np.log10(spectrogram), fs

paths = [
    "outputs/icosahedron/2023-08-06T12:27:52+02:00.yaml",
    "outputs/icosahedron/2023-08-06T13:00:51+02:00.yaml",
    "outputs/icosahedron/2023-08-06T12:20:34+02:00.yaml",
    "outputs/icosahedron/2023-08-06T13:01:17+02:00.yaml",
    ]

rc('font', **{'family': 'serif', 'serif': ['Computer Modern'], 'size': 11})
rc('text', usetex=True)
# fig, ax = plt.subplots(1, 1, figsize=(6.8, 3.4), sharex=True, sharey=True)
fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(6.8, 3.4))

y_min, y_max = -0.5, 0.5
x_min, x_max = 0.0, 5e-6

for i, file_path in enumerate(paths):
    with open(file_path, "r") as file:
        yaml_data = yaml.safe_load(file)

    pc = yaml_data["values"]["system"]["physicsconfig"]

    for xys in yaml_data["xys"]:
        xs = []
        ys = []
        for xy in xys:
            xs.append(xy["x"])
            ys.append(xy["y"])

    ax1.set_ylim(y_min, y_max)
    ax1.set_xlim(x_min, x_max)
    ax1.set_yticks([-0.5, 0, 0.5])
    ax1.set_ylabel(r"$\langle S_Z^{(0)}\rangle(t)$")
    ax1.plot(xs, ys, label=r"$\beta$ = " + f"{pc['tiltangle']}" + r"$\pi$")
    ax1.set_xlabel(r"$t\,[\mathrm{sec}]$")


    fft_vals = fft.fft(ys)
    freqs = fft.fftfreq(len(xs), xs[1] - xs[0])

    top_indices = np.argsort(np.abs(fft_vals))[-5:]
    top_freqs = freqs[top_indices]
    top_amplitudes = np.abs(fft_vals[top_indices])
    times = [1/f for f in top_freqs if f > 0 ]
    print(times)
    #[ax2.axvline(x=t, color='xkcd:orange', ls='dotted') for t in times] # , label=r'$t = ' + f'{round(1e6*t, 2)}' + r' \,\mu\mathrm{s}$') for t in times]
    time, frequencies, color, fs = spectrogram(xs, ys)
    ax2.pcolormesh(time, frequencies, color)
    ax2.set_xlabel("Time (s)")
    ax2.set_ylabel("Frequency (Hz)")
    #ax2.colorbar(label="Power spectral density (dB/Hz)")
    ax2.set_ylim([0, fs/2.])

#ax.axvline(x=0.35e-6, ls='dashed', label=r'$t = 0.35\,\mu\mathrm{s}$')


fig.legend(bbox_to_anchor=(0.97, 0.95), ncol=2, loc='upper right')
fig.tight_layout()
plt.savefig(f"figures/icosa-decay-time-fft.png", dpi=600)
