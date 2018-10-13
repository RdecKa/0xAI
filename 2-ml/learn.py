import sys
import getopt
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

from sklearn.tree import DecisionTreeRegressor
from sklearn.model_selection import train_test_split

import decision_tree as dt


class Learner:

    def __init__(self, models):
        self.models = models

    def get_models(self):
        return self.models


class Model:

    def __init__(self, model):
        self.model = model

    def fit(self, X, y):
        self.model.fit(X, y)

    def predict(self, X):
        return self.model.predict(X)

    def score(self, X, y):
        return self.model.score(X, y)


def main(argv):
    # Read flags
    datafile = "sample_data/data.in"
    outfolder = "sample_out/"
    try:
        opts, args = getopt.getopt(argv, "d:o:")
    except getopt.GetoptError:
        print("Error parsing the command line arguments")
        sys.exit(1)
    for o, a in opts:
        if o == "-d":
            datafile = a
        if o == "-o":
            outfolder = a

    # Read data from file
    print("Reading data from file:", datafile)
    df = pd.read_csv(datafile, comment="#",
                     dtype={"value": np.float64, "num_stones": np.uint8,
                            "occ_red_rows": np.uint8, "occ_red_cols": np.uint8,
                            "occ_blue_rows": np.uint8, "occ_blue_cols": np.uint8,
                            "red_p0": np.uint8, "blue_p0": np.uint8,
                            "red_p1": np.uint8, "blue_p1": np.uint8,
                            "red_p2": np.uint8, "blue_p2": np.uint8,
                            "red_p3": np.uint8, "blue_p3": np.uint8,
                            "red_p4": np.uint8, "blue_p4": np.uint8,
                            "lp": bool,
                            "dtc": np.uint8})

    y = df["value"]
    X = df.drop(columns=["value"])

    feature_names = X.columns

    # Split
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.1,
                                                        random_state=4224)

    dt_models = [
        DecisionTreeRegressor(max_depth=5, min_samples_leaf=5),
        DecisionTreeRegressor(max_depth=10, min_samples_leaf=5),
        DecisionTreeRegressor(max_depth=None, min_samples_leaf=1),
        DecisionTreeRegressor(max_depth=None, min_samples_leaf=10),
        DecisionTreeRegressor(max_depth=None, min_samples_leaf=20),
        DecisionTreeRegressor(max_depth=None, min_samples_leaf=50),
    ]
    dt_models = [dt.DecisionTreeModel(m, index) for index, m in enumerate(dt_models)]

    learners = [
        dt.DecisionTreeLearner(dt_models),
    ]

    for learner in learners:
        # Create a plot
        plt.figure(figsize=(10, 6))
        plt.plot(y_test.tolist(), label="actual values", linewidth=0.7)
        plt.ylim(-1.2, 1.2)

        # Create file for statistics
        with open(outfolder + "stats_" + learner.short_name() + ".txt", "w") \
                as stats_file:
            models = learner.get_models()
            for model_index in range(len(models)):
                model = models[model_index]

                # Train
                model.fit(X_train, y_train)

                # Predict
                y1 = model.predict(X_test)

                # Print statistics
                stats_file.write("##########################################\n")
                stats_file.write("Statistics for:" + str(model) + "\n")

                stats_file.write("Feature importances:\n")
                stats_file.write(model.feature_importances(X.keys(), outfolder))

                sc = model.score(X_test, y_test)
                stats_file.write("SCORE: " + str(sc) + "\n")

                # Output some custom properties for current model
                model.custom_output(model_index, feature_names, outfolder)

                # Add to plot
                plt.plot(y1, label="predicted values - " + model.name(),
                         linewidth=0.7)

        # Output some custom properties for all models
        learner.custom_output(feature_names, outfolder)

        plt.xlabel("Samples")
        plt.ylabel("Value")
        plt.title("Comparison of predicted and actual values")
        plt.legend(fontsize="x-small")
        plt.savefig(outfolder + "plot.pdf")
        plt.close()


if __name__ == "__main__":
    main(sys.argv[1:])
