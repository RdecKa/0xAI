import matplotlib.pyplot as plt

from sklearn import tree

import tree_to_code as ttc

from learn import Learner, Model


class DecisionTreeLearner(Learner):

    @staticmethod
    def name():
        return "DecisionTreeLearner"

    @staticmethod
    def short_name():
        return "dtl"

    @staticmethod
    def custom_output(feature_names, outfolder):
        # Create file for Go sample struct
        ttc.write_sample_file(outfolder, feature_names)


class DecisionTreeModel(Model):

    def __init__(self, model, ID):
        super().__init__(model)
        self.ID = "dtl_" + str(ID)

    def __str__(self):
        return str(self.model)

    def name(self):
        return "dt (max_depth=" + str(self.model.max_depth) + \
               ", min_leaf=" + str(self.model.min_samples_leaf) + ")"

    def feature_importances(self, feature_names, outfolder):
        plt.figure()
        plt.bar(feature_names, self.model.feature_importances_)
        plt.gcf().subplots_adjust(bottom=0.25)
        plt.xticks(rotation='vertical')
        plt.savefig(outfolder + "features_" + self.ID + ".pdf")
        plt.close()

        fi = zip(feature_names, self.model.feature_importances_)
        s = ""
        for (k, v) in fi:
            s += "\t" + k + ": " + str(v) + "\n"

        return s

    def custom_output(self, model_index, feature_names, outfolder):
        # Visualize trees
        tree.export_graphviz(self.model, feature_names=feature_names,
                             out_file=outfolder + "tree" + str(model_index) + ".dot")

        # Output Go code
        ttc.tree_to_go_code(self.model.tree_, feature_names, 0, model_index,
                            outfolder)
