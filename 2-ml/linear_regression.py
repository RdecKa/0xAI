import math

from sklearn.linear_model import LinearRegression

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

    def __init__(self, feature_names, ID, splits):
        super().__init__(None, feature_names)
        self.ID = "lrl_" + str(ID)
        self.submodels = [(split, LinearRegression(normalize=True, n_jobs=-1)) for split in splits]
        self.submodels.append((math.inf, LinearRegression(normalize=True, n_jobs=-1)))

    def __str__(self):
        return str(self.submodels)

    @staticmethod
    def name():
        return "lr"

    def feature_importances(self):
        s = []
        for submodel in self.submodels:
            s.append(submodel[1].coef_)
        return s

    def custom_output(self, model_index, outfolder):
        self.lr_to_go_code(model_index, outfolder)

    def lr_to_go_code(self, model_index, outfolder):
        coefficients = {}
        for submodel in self.submodels:
            coefficients[submodel[0]] = submodel[1].coef_

        with open(outfolder + "linear" + str(model_index) + "code.go", "w") as code_file:
            def get_one_factor(feature_name, coefficient):
                return "({})*float64(s.{})".format(coefficient, feature_name)

            def get_one_submodel(key):
                s = "\t\treturn {}".format(get_one_factor(self.feature_names[0], coefficients[key][0]))
                for (n, c) in zip(self.feature_names[1:], coefficients[key][1:]):
                    s += " +\n\t\t\t{}".format(get_one_factor(n, c))
                s += "\n"
                return s
            code_file.write("// Package ab (Code generated by a Python script)\n")
            code_file.write("package ab\n\n")
            code_file.write("func (s Sample) getEstimatedValue() float64 {\n")
            code_file.write("\tswitch {\n")

            for ind in range(len(self.submodels)):
                submodel = self.submodels[ind]
                if ind == len(self.submodels) - 1:
                    s = "\tdefault"
                else:
                    s = "\tcase s.num_stones <= " + str(submodel[0])
                s += ":\n"
                s += (get_one_submodel(submodel[0]))
                code_file.write(s)
            code_file.write("\t}")  # end of switch
            code_file.write("\n}\n")

    def group_func(self, X, ind_X):
        """
        Tells to which group does a sample with index ind in the DataFrame X
        belong
        :param X: DataFrame
        :param ind: index of the investigated sample
        :return: group ID (from 0 to len(self.submodels)-1)
        """
        for ind in range(len(self.submodels)):
            submodel = self.submodels[ind]
            if X["num_stones"].loc[ind_X] <= submodel[0]:
                return ind
        return len(self.submodels) - 1

    def fit(self, X, y):
        data_in_subsets = X.groupby(lambda i: self.group_func(X, i))

        for ind in range(len(self.submodels)):
            if ind not in data_in_subsets.groups.keys():
                # No samples that would fall in this group
                self.submodels[ind] = None
                continue

            inp = X.loc[data_in_subsets.groups[ind]]
            out = y.loc[data_in_subsets.groups[ind]]
            self.submodels[ind][1].fit(inp, out)

        # Keep only models that are not None
        models = []
        for submodel in self.submodels:
            if submodel is not None:
                models.append(submodel)
        models[-1] = (math.inf, models[-1][1])  # Set the split of the last model to infinity

        self.submodels = models

    def predict(self, X):
        y = [0] * len(X)
        for (ind, X_ind) in enumerate(X.index):
            group_id = self.group_func(X, X_ind)
            y[ind] = self.submodels[group_id][1].predict([X.loc[X_ind]])

        return y

    def score(self, X, y):
        data_in_subsets = X.groupby(lambda i: self.group_func(X, i))

        s = {}
        for ind in range(len(self.submodels)):
            if ind not in data_in_subsets.groups.keys():
                # No samples of this group in test data
                s[ind] = None
                continue

            submodel = self.submodels[ind]
            inp = X.loc[data_in_subsets.groups[ind]]
            out = y.loc[data_in_subsets.groups[ind]]
            s[submodel[0]] = (submodel[1].score(inp, out), len(inp))
        return s
