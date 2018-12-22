import sys
import getopt
import math
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import seaborn as sns

from sklearn.model_selection import train_test_split
from sklearn.metrics.pairwise import cosine_distances
from sklearn.preprocessing import scale
from itertools import product

import decision_tree as dt
import linear_regression as lr
import analyse as ana


def write_sample_file(outfolder, feature_names):
    with open(outfolder + "sample.go", "w") as sample_file:
        sample_file.write("// Package ab (Code generated by a Python script)\n")
        sample_file.write("package ab\n\n")
        sample_file.write("type Sample struct {\n\t")
        sample_file.write(feature_names[0])
        for f in feature_names[1:]:
            sample_file.write(", " + f)
        sample_file.write(" int\n")
        sample_file.write("}")


def sort_value(x):
    minimum = x[2:]
    minimum = math.inf if minimum == "inf" else int(minimum)
    return minimum, x[0]


def generate_colors_and_patterns(n):
    patterns = ("-", "+", "x", "\\", "/", "|", "*", "o", "O", ".")
    third = math.ceil(n ** (1/3))
    intervals = [x/third for x in range(third)]
    colors = list(product(intervals, intervals, intervals))
    for i in range(len(colors)):
        colors[i] = (colors[i], patterns[i % len(patterns)])
    return colors


def main(argv):
    # Read flags
    datafile = "sample_data/data.in"
    outfolder = "sample_out/"
    data_analysis = True
    try:
        opts, args = getopt.getopt(argv, "ad:o:")
    except getopt.GetoptError:
        print("Error parsing the command line arguments")
        sys.exit(1)
    for o, a in opts:
        if o == "-d":
            datafile = a
        elif o == "-o":
            outfolder = a
        elif o == "-a":
            data_analysis = False

    # Read data from file
    print("Reading data from file:", datafile)
    df = pd.read_csv(datafile, comment="#",
                     dtype={"value": np.float64, "num_stones": np.uint8,
                            "lp": np.uint8, "sdtc_r": np.uint16, "sdtc_b": np.uint16,
                            "rec_r": np.uint16, "rec_b": np.uint16,
                            "occ_red_rows": np.uint8, "occ_red_cols": np.uint8,
                            "occ_blue_rows": np.uint8, "occ_blue_cols": np.uint8,
                            "red_p0": np.uint8, "blue_p0": np.uint8,
                            "red_p1": np.uint8, "blue_p1": np.uint8,
                            "red_p2": np.uint8, "blue_p2": np.uint8,
                            "red_p3": np.uint8, "blue_p3": np.uint8,
                            "red_p4": np.uint8, "blue_p4": np.uint8,
                            "red_p5": np.uint8, "blue_p5": np.uint8,
                            "red_p6": np.uint8, "blue_p6": np.uint8,
                            "red_p7": np.uint8, "blue_p7": np.uint8,
                            "red_p8": np.uint8, "blue_p8": np.uint8,
                            "red_p9": np.uint8, "blue_p9": np.uint8,
                            "red_p10": np.uint8, "blue_p10": np.uint8,
                            "red_p11": np.uint8, "blue_p11": np.uint8,
                            "red_p12": np.uint8, "blue_p12": np.uint8,
                            "red_p13": np.uint8, "blue_p13": np.uint8,
                            "red_p14": np.uint8, "blue_p14": np.uint8,
                            "red_p15": np.uint8, "blue_p15": np.uint8,
                            "red_p16": np.uint8, "blue_p16": np.uint8,
                            "red_p17": np.uint8, "blue_p17": np.uint8,
                            "red_p18": np.uint8, "blue_p18": np.uint8,
                            "red_p19": np.uint8, "blue_p19": np.uint8,
                            "red_p20": np.uint8, "blue_p20": np.uint8,
                            "red_p21": np.uint8, "blue_p21": np.uint8,
                            "red_p22": np.uint8, "blue_p22": np.uint8,
                            "red_p23": np.uint8, "blue_p23": np.uint8,
                            })

    y = df["value"]
    X = df.drop(columns=["value"])

    feature_names = X.columns
    write_sample_file(outfolder, feature_names)

    feature_colors = generate_colors_and_patterns(len(feature_names))

    if data_analysis:
        ana.analyse(outfolder, df, feature_names)

    # Split
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.1,
                                                        random_state=4224)

    # Create decision tree models
    dt_models_args = [  # [max_depth, min_samples_leaf]
        [5, 5],
        [10, 5],
        [None, 1],
        [None, 10],
        [None, 20],
        [None, 50],
    ]

    # Create linear regression models
    lr_models_args = [
        [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 15, 18, 22, 28, 36, 46, 58, 70, 85, 100],
        [x for x in range(0, 48)],
        [0, 3, 5, 7, 9, 12, 15, 19, 25, 32, 40, 50, 65, 85, 100]
    ]

    learners = [
        dt.DecisionTreeLearner(dt_models_args, feature_names),
        lr.LinearRegressionLearner(lr_models_args, feature_names),
    ]

    for learner in learners:
        print("Making learner:" + str(learner) + " ...")
        # Create file for statistics
        with open(outfolder + "stats_" + learner.short_name() + ".txt", "w") \
                as stats_file:
            models = learner.get_models()

            for model_index in range(len(models)):
                model = models[model_index]
                print("Training model: " + str(model) + " ...")

                # Train
                model.fit(X_train, y_train)

                # Print statistics
                stats_file.write("##########################################\n")
                stats_file.write("Statistics for: " + str(model) + "\n")

                # Create a plot for feature importances
                submodel_groups = model.get_num_submodels()
                plt_setup = []
                for i, num_submodels in enumerate(submodel_groups):
                    if num_submodels > 3:
                        num_cols_in_plot = 4
                        num_rows_in_plot = math.ceil(num_submodels / num_cols_in_plot)
                    else:
                        num_cols_in_plot = num_submodels
                        num_rows_in_plot = 1
                    fig_name = learner.short_name() + str(model_index) + str(i)
                    plt.figure(fig_name, figsize=(10, 20))
                    plt_setup.append((num_cols_in_plot, num_rows_in_plot, fig_name))

                feature_importances = model.feature_importances(feature_names)
                vectors = [[], []]
                model_names = [[], []]

                for ind, feat in enumerate(feature_importances):
                    if len(feat) == 0:
                        continue

                    name = feat[0]
                    fi = feat[1]

                    if isinstance(learner, dt.DecisionTreeLearner):
                        stats_file.write("Feature importances:\n")
                    else:
                        stats_file.write("Feature coefficients (split " + name + "):\n")

                    s = ""
                    for (n, v) in zip(feature_names, fi):
                        s += "\t" + n + ": " + str(v) + "\n"
                    stats_file.write(s)

                    if len(plt_setup) == 1:
                        num_cols_in_plot, num_rows_in_plot, fig_name = plt_setup[0]
                        index = ind + 1
                    else:
                        if ".r_" in name:
                            num_cols_in_plot, num_rows_in_plot, fig_name = plt_setup[0]
                            index = ind + 1
                            vectors[0].append(fi)
                            model_names[0].append(name)
                        else:  # if "._b" in name
                            num_cols_in_plot, num_rows_in_plot, fig_name = plt_setup[1]
                            index = ind - submodel_groups[0] + 1
                            vectors[1].append(fi)
                            model_names[1].append(name)

                    total = sum([abs(f) for f in fi])
                    normalised = [abs(f) / total for f in fi]

                    plt.figure(fig_name).add_subplot(num_rows_in_plot, num_cols_in_plot, index)
                    acc = 0
                    for i, f in enumerate(normalised):
                        plt.bar([name], f, bottom=acc, color=feature_colors[i][0],
                                edgecolor="black", hatch=feature_colors[i][1],
                                label=feature_names[i])
                        acc += f

                sc = model.score(X_test, y_test)
                stats_file.write("SCORE:\n")
                s = ""
                if isinstance(learner, lr.LinearRegressionLearner):
                    sorted_keys = sorted(sc, key=sort_value)
                else:
                    sorted_keys = sorted(sc)
                for key in sorted_keys:
                    v = sc[key]
                    if v is None:
                        s += "{}: No testing samples\n".format(key)
                    else:
                        s += "{}: {} ({} testing samples)\n".format(key, v[0], v[1])
                stats_file.write(s)

                # Output some custom properties for current model
                model.custom_output(model_index, outfolder)

                if len(plt_setup) == 1:
                    plt.legend(loc=7)
                    plt.savefig(outfolder + "features_" + learner.short_name() + "_" + str(model_index) + ".pdf")
                    plt.close()
                else:
                    f = plt.figure(plt_setup[0][2])
                    handles, labels = f.axes[0].get_legend_handles_labels()
                    handles = handles[:len(feature_names)]
                    labels = labels[:len(feature_names)]

                    f.legend(handles, labels, loc=7)
                    plt.savefig(outfolder + "features_" + learner.short_name() + "_" + str(model_index) + "_red.pdf")
                    plt.close(f)

                    f = plt.figure(plt_setup[1][2])
                    f.legend(handles, labels, loc=7)
                    plt.savefig(outfolder + "features_" + learner.short_name() + "_" + str(model_index) + "_blue.pdf")
                    plt.close(f)

                if isinstance(learner, lr.LinearRegressionLearner):
                    for i in range(2):
                        matrix = pd.DataFrame(vectors[i], columns=feature_names, index=model_names[i])

                        matrix_adjusted = scale(matrix)
                        cs = cosine_distances(matrix_adjusted)

                        fig, ax = plt.subplots(figsize=(14, 10))
                        ax = sns.heatmap(cs, xticklabels=model_names[i], yticklabels=model_names[i], square=True, ax=ax)
                        plt.savefig(outfolder + "heatmap_" + learner.short_name() + "_" + str(model_index) + "." + str(i) + ".pdf")
                        plt.close()

if __name__ == "__main__":
    main(sys.argv[1:])
