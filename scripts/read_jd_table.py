import csv
import sys
import yaml
import matplotlib.pyplot as plt

def read(jd, jk):

    x = []
    # read in JD tsv table
    with open(jd) as file_jd:
        tsv = csv.reader(file_jd, delimiter="\t")
        i = 0
        jd_y = []
        for line in tsv:
            if i > 32:
                jd_y.append(float(line[-1]))
            i += 1

    # read in JK yaml file
    with open(jk) as file_jk:
        yml = yaml.safe_load(file_jk)
        jk_y = []
        for xy in yml["xys"][0]:
            jk_y.append(xy["y"])
            x.append(xy["x"])
    # print([jk_y[i] - jd_y[i] for i, _ in enumerate(jk_y)])
    plt.plot(x, jk_y, label="JK")
    plt.plot(x, jd_y, label="JD")
    plt.legend()
    plt.xlabel(r"$t$ $[s]$")
    plt.ylabel(r"$\langle S^{(0)}_z \rangle$")
    plt.savefig("jd_jk_comparison_07032023.png")


if __name__ == '__main__':
    paths = sys.argv[1:]
    read(paths[0], paths[1]) # 0 -> tsv, 1 -> yaml