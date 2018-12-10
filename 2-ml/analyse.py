import matplotlib.pyplot as plt
import seaborn as sns
import numpy as np


def analyse(outfolder, df, feature_names):
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
        stats[v] = [[], []]  # [groups, singletons]
    num_stones_index = feature_names.get_loc("num_stones")

    with open(outfolder + "same_attr_comparison.txt", "w") as comp_file:
        stds = [sd for sd in data_same_attributes["value"].std() if not np.isnan(sd)]
        comp_file.write("Number of unique groups: {}\n".format(len(stds)))
        comp_file.write("\tAverage standard deviation: {}\n".format(np.mean(stds)))
        comp_file.write("\tStandard deviation of standard deviations: {}\n".format(np.std(stds)))
        comp_file.write("\n")

        for key in data_same_attributes.groups.keys():
            ind = data_same_attributes.groups[key]
            num_stones = key[num_stones_index]
            if len(ind) > 1:
                val = df.loc[ind]
                m = val["value"].mean()
                s = val["value"].std()
                # Append tuple: (key = tuple of attribute values, mean, standard deviation):
                stats[num_stones][0].append((key, m, s))
            else:
                stats[num_stones][1].append(key)

        means, stds = [], []
        keys, keys_sizes = [], []
        sizes_groups, sizes_singletons = [], []
        for key, value in stats.items():
            groups = value[0]
            singletons = value[1]
            size_single = len(singletons)

            comp_file.write("num_stones: " + str(key) + "\n")

            if len(groups) > 0:
                std_devs = [x[2] for x in groups]
                size = len(std_devs)
                mean = np.mean(std_devs)
                std = np.std(std_devs)

                means.append(mean)
                stds.append(std)
                keys.append(key)
                sizes_groups.append(size)

                comp_file.write("\tNumber of groups with same attributes: {}\n".format(size))
                comp_file.write("\tAverage standard deviation: {}\n".format(mean))
                comp_file.write("\tStandard deviation of standard deviations: {}\n".format(std))

                if size_single == 0:
                    keys_sizes.append(key)
                    sizes_singletons.append(0)

            if size_single > 0:
                sizes_singletons.append(size_single)
                keys_sizes.append(key)

                comp_file.write("\tNumber of singleton groups: {}\n".format(size_single))

                if len(groups) == 0:
                    sizes_groups.append(0)

        fig1, ax1 = plt.subplots()
        color1 = "#0F2F8C"
        color2 = "#A60303"
        color3 = "#33BB33"

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
        plt.yscale("log")
        ax2.bar(keys_sizes, sizes_groups, color=color2, alpha=0.15)
        ax2.tick_params("y", colors=color2)
        ax2.set_ylabel("Number of groups with a given number of stones", color=color2)

        ax2.bar(keys_sizes, sizes_singletons, color=color3, alpha=0.15, bottom=sizes_groups)

        ax1.grid(True)
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
