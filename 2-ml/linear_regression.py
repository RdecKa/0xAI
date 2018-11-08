import math

from sklearn.linear_model import LinearRegression
from sklearn.feature_selection import VarianceThreshold

from learner import Learner
from model import Model


class LinearRegressionLearner(Learner):

    def __init__(self, splits, feature_names):
        super().__init__(None, feature_names)
        self.models = [LinearRegressionModel(feature_names, ind, splits[ind]) for ind, m in enumerate(splits)]

    @staticmethod
    def name():
        return "LinearRegressionLearner"

    @staticmethod
    def short_name():
        return "lrl"


class LinearRegressionModel(Model):

    class Submodel:

        def __init__(self, color, maximum, model):
            self.color = color
            self.maximum = maximum
            self.model = model
            self.used_features = []

        def __str__(self):
            return "({} - {} - {})".format(self.color, self.maximum, self.model)

        def get_ID(self):
            return self.color + "_" + str(self.maximum)

    def __init__(self, feature_names, ID, splits):
        super().__init__(None, feature_names)
        self.ID = "lrl_" + str(ID)
        self.splits = splits
        self.submodels = [[], []]
        self.submodels[0] = [None] * (len(splits)+1)  # For red player
        self.submodels[1] = [None] * (len(splits)+1)  # For blue player
        for split_index in range(len(splits)):
            maximum = splits[split_index]
            self.submodels[0][split_index] = \
                self.Submodel("r", maximum, LinearRegression(normalize=True, n_jobs=-1))
            self.submodels[1][split_index] = \
                self.Submodel("b", maximum, LinearRegression(normalize=True, n_jobs=-1))
        self.submodels[0][-1] = \
            self.Submodel("r", math.inf, LinearRegression(normalize=True, n_jobs=-1))
        self.submodels[1][-1] = \
            self.Submodel("b", math.inf, LinearRegression(normalize=True, n_jobs=-1))

    def __str__(self):
        s = ""
        for c in range(2):
            for submodel in self.submodels[c]:
                s += str(submodel)
        return s

    @staticmethod
    def name():
        return "lr"

    def feature_importances(self, feature_names):
        s = []
        for c in range(2):
            for submodel in self.submodels[c]:
                if submodel is None:
                    continue
                fi = [0] * len(feature_names)
                for (ind, fn) in enumerate(feature_names):
                    if fn in submodel.used_features:
                        i = submodel.used_features.index(fn)
                        fi[ind] = submodel.model.coef_[i]
                s.append((self.ID + "." + submodel.get_ID(), fi))
        return s

    def custom_output(self, model_index, outfolder):
        self.lr_to_go_code(model_index, outfolder)

    def lr_to_go_code(self, model_index, outfolder):
        with open(outfolder + "linear" + str(model_index) + "code.go", "w") as code_file:
            def one_color_to_go_code(color):
                if color == "r":
                    submodels = self.submodels[0]
                elif color == "b":
                    submodels = self.submodels[1]
                else:
                    raise ValueError(color)

                coefficients = {}
                for ind in range(len(submodels)):
                    submodel = submodels[ind]
                    coefficients[submodel.get_ID()] = submodel.model.coef_

                code_file.write("\t\tswitch {\n")

                for ind in range(len(submodels)):
                    submodel = submodels[ind]
                    if ind == len(submodels) - 1:
                        s = "\t\tdefault"
                    else:
                        s = "\t\tcase s.num_stones <= " + str(submodel.maximum)
                    s += ":\n"
                    s += (get_one_submodel(submodel, coefficients))
                    code_file.write(s)

                code_file.write("\t\t}\n")  # end of switch

            def get_one_factor(feature_name, coefficient):
                return "({})*float64(s.{})".format(coefficient, feature_name)

            def get_one_submodel(submodel, coefficients):
                key = submodel.get_ID()
                z = [(fn, c) for (fn, c) in zip(submodel.used_features, coefficients[key]) if c != 0]
                if len(z) == 0:
                    # All coefficients are zero
                    return "\t\t\treturn 0\n"

                s = "\t\t\treturn {}".format(get_one_factor(*z[0]))
                for (n, c) in z[1:]:
                    s += " +\n\t\t\t\t{}".format(get_one_factor(n, c))
                s += "\n"
                return s

            code_file.write("// Package ab (Code generated by a Python script)\n")
            code_file.write("package ab\n\n")
            code_file.write("func getEstimatedValueLR(s Sample) float64 {\n")

            code_file.write("\tif s.lp == 0 {\n")
            one_color_to_go_code("r")
            code_file.write("\t} else {\n")
            one_color_to_go_code("b")
            code_file.write("\t}\n")

            code_file.write("}\n")

    def group_func(self, X, ind_X):
        """
        Tells to which group does a sample with index ind in the DataFrame X
        belong
        :param X: DataFrame
        :param ind: index of the investigated sample
        :return: group ID (from 0 to len(self.submodels)-1)
        """
        color = 1 if X["lp"].loc[ind_X] else 0
        for ind in range(len(self.submodels[color])):  # Take advantage of submodels being ordered
            submodel = self.submodels[color][ind]
            if X["num_stones"].loc[ind_X] <= submodel.maximum:
                return ind, color
        return len(self.submodels[color]) - 1, color

    def group_func_train(self, X, ind_X):
        color = 1 if X["lp"].loc[ind_X] else 0
        for ind in range(len(self.splits)):
            if X["num_stones"].loc[ind_X] <= self.splits[ind]:
                return ind, color
        return len(self.splits), color

    def fit(self, X, y):
        data_in_subsets = X.groupby(lambda i: self.group_func_train(X, i))
        all_features = [a for a in X]

        for c in range(2):
            for ind in range(len(self.submodels[c])):
                key = (ind, c)
                if key not in data_in_subsets.groups.keys():
                    # No samples that would fall in this group
                    self.submodels[c][ind] = None
                    continue

                inp = X.loc[data_in_subsets.groups[key]]
                out = y.loc[data_in_subsets.groups[key]]

                if len(inp) > 1:
                    # Remove redundant features (when all samples in the group
                    # have the same value of an attribute)
                    feature_selector = VarianceThreshold()
                    feature_selector.fit(inp)

                    # Get indices of attributes that are kept:
                    imp_feat_ind = feature_selector.get_support(True)
                    imp_feat = [all_features[i] for i in imp_feat_ind]
                    inp = inp[imp_feat]
                    self.submodels[c][ind].used_features = imp_feat
                else:
                    self.submodels[c][ind].used_features = all_features

                self.submodels[c][ind].model.fit(inp, out)

        # Keep only models that are not None
        models = [[], []]
        for c in range(2):
            for submodel in self.submodels[c]:
                if submodel is not None:
                    models[c].append(submodel)
        models[0][-1].maximum = math.inf  # Set the split of the last model to infinity
        models[1][-1].maximum = math.inf

        self.submodels = models

    def predict(self, X):
        y = [0] * len(X.index)
        for (ind, X_ind) in enumerate(X.index):
            i, c = self.group_func(X, X_ind)
            submodel = self.submodels[c][i]
            y[ind] = submodel.model.predict([X.loc[X_ind, submodel.used_features]])
        return y

    def score(self, X, y):
        data_in_subsets = X.groupby(lambda i: self.group_func(X, i))

        s = {}
        for c in range(2):
            for ind in range(len(self.submodels[c])):
                key = (ind, c)
                submodel = self.submodels[c][ind]
                if key not in data_in_subsets.groups.keys():
                    # No samples of this group in test data
                    s[submodel.get_ID()] = None
                    continue

                inp = X.loc[data_in_subsets.groups[key], submodel.used_features]
                out = y.loc[data_in_subsets.groups[key]]
                s[submodel.get_ID()] = (submodel.model.score(inp, out), len(inp))

        return s
