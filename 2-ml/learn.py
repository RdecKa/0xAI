import sys
import getopt
import math
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import seaborn as sns

from sklearn.model_selection import train_test_split
from matplotlib.ticker import FixedLocator, FormatStrFormatter

import decision_tree as dt
import linear_regression as lr


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

    if data_analysis:
        print("Analysing data ...")

        # Create a plot showing attribute distributions
        plt.gcf().subplots_adjust(bottom=0.25)
        sns.boxplot(data=df, color="darkorchid")
        plt.xticks(rotation="vertical")
        plt.savefig(outfolder + "attrs_boxplot.pdf")
        plt.close()

        # Check how similar values do samples with the same attributes have
        data_same_attributes = df.groupby([a for a in feature_names])

        stats = {}
        for v in df.num_stones.unique():
            stats[v] = []
        num_stones_index = feature_names.get_loc("num_stones")

        with open(outfolder + "same_attr_comparison.txt", "w") as comp_file:
            stds = [sd for sd in data_same_attributes["value"].std() if not np.isnan(sd)]
            comp_file.write("Number of unique groups: {}\n".format(len(stds)))
            comp_file.write("\tAverage standard deviation: {}\n".format(np.mean(stds)))
            comp_file.write("\tStandard deviation of standard deviations: {}\n".format(np.std(stds)))
            comp_file.write("\n")

            for key in data_same_attributes.groups.keys():
                ind = data_same_attributes.groups[key]
                if len(ind) > 1:
                    val = df.loc[ind]
                    m = val["value"].mean()
                    s = val["value"].std()
                    num_stones = key[num_stones_index]
                    # Append tuple: (key = tuple of attribute values, mean, standard deviation):
                    stats[num_stones].append((key, m, s))

            means = []
            stds = []
            keys = []
            sizes = []
            for key, value in stats.items():
                if len(value) == 0:
                    continue

                std_devs = [x[2] for x in value]
                size = len(std_devs)
                mean = np.mean(std_devs)
                std = np.std(std_devs)

                comp_file.write("num_stones: " + str(key) + "\n")
                comp_file.write("\tNumber of groups with same attributes: {}\n".format(size))
                comp_file.write("\tAverage standard deviation: {}\n".format(mean))
                comp_file.write("\tStandard deviation of standard deviations: {}\n".format(std))

                means.append(mean)
                stds.append(std)
                keys.append(key)
                sizes.append(size)

            fig1, ax1 = plt.subplots()
            color1 = "#0F2F8C"
            color2 = "#A60303"

            means = np.array(means)
            stds = np.array(stds)
            lower_bounds = means - stds
            upper_bounds = means + stds

            # Plot mean of standard deviations
            ax1.plot(keys, means, color=color1)
            ax1.fill_between(keys, lower_bounds, upper_bounds, alpha=0.3)
            ax1.set_ylabel("Mean of standard deviations", color=color1)
            ax1.set_xlabel("Number of stones")
            ax1.tick_params("y", colors=color1)
            ax1.set_title("Mean and std of stds of sample values in groups, " +
                          "created by samples\nhaving the same attribute values, " +
                          "then grouped by number of stones")

            # Plot number of groups with a given number of stones
            ax2 = ax1.twinx()
            ax2.bar(keys, sizes, color=color2, alpha=0.15)
            ax2.tick_params("y", colors=color2)
            ax2.set_ylabel("Number of groups with a given number of stones", color=color2)

            # Transform the yticks on the first axis to match with the second
            y1_bot, y1_top = ax1.get_ylim()
            y2_bot, y2_top = ax2.get_ylim()

            def transform_ticks(yy):
                return y1_bot + (yy-y2_bot)/(y2_top-y2_bot) * (y1_top-y1_bot)

            new_y1_ticks = transform_ticks(ax2.get_yticks())
            new_y1_ticks = np.append(new_y1_ticks, 0)
            ax1.yaxis.set_major_locator(FixedLocator(new_y1_ticks))
            ax1.yaxis.set_major_formatter(FormatStrFormatter('%.2f'))

            ax1.grid(True)
            plt.xticks(keys)
            plt.savefig(outfolder + "std.pdf")
            plt.close()

        # Create a plot showing relations between attributes
        attrs = df.keys()
        group1 = [a for a in attrs if "red_p" in a]
        group2 = [a for a in attrs if "blue_p" in a]
        group3 = [a for a in attrs if "occ_" in a]
        group4 = [a for a in attrs if a not in group1 and a not in group2 and a not in group3]
        list_of_groups = [group1, group2, group3, group4]

        for a in range(len(list_of_groups)):
            for b in range(a, len(list_of_groups)):
                print("Creating a pairplot of:")
                print("\t" + str(list_of_groups[a]))
                print("\t" + str(list_of_groups[b]))
                sns.pairplot(df, x_vars=list_of_groups[a], y_vars=list_of_groups[b],
                             plot_kws={"alpha": 0.03, "s": 80})
                plt.savefig(outfolder + "attrs_pairplot_" + str(a) + "_" + str(b) + ".pdf")
                plt.close()

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

            # Create a plot for feature importances
            fig, ax = plt.subplots()
            positions = np.arange(len(feature_names))

            for model_index in range(len(models)):
                model = models[model_index]
                print("Training model: " + str(model) + " ...")

                # Train
                model.fit(X_train, y_train)

                # Predict
                y1 = model.predict(X_test)

                # Print statistics
                stats_file.write("##########################################\n")
                stats_file.write("Statistics for: " + str(model) + "\n")

                feature_importances = model.feature_importances(feature_names)

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

                    bar_width = (1 / (len(models) * len(feature_importances)))
                    pos = positions + (model_index * len(feature_importances) + ind) * bar_width
                    ax.bar(pos, fi, width=bar_width, label=name, tick_label=feature_names)

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

            plt.gcf().subplots_adjust(bottom=0.35, right=0.85)
            plt.xticks(positions-bar_width/2, rotation="vertical")
            plt.grid(True)
            ax.set_title("Influence of attributes")
            ax.legend(loc="center left", bbox_to_anchor=(1, 0.5), fontsize="x-small")
            ax.set_xlabel("Attribute")
            ax.set_ylabel("Influence")

            plt.yscale("linear")
            plt.savefig(outfolder + "features_" + learner.short_name() + "_lin.pdf")
            plt.yscale("symlog")
            plt.savefig(outfolder + "features_" + learner.short_name() + "_log.pdf")
            plt.close()


if __name__ == "__main__":
    main(sys.argv[1:])
