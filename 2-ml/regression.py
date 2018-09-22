import sys, getopt
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import math

from sklearn import tree
from sklearn.tree import DecisionTreeRegressor, _tree
from sklearn.model_selection import train_test_split

def tree_to_go_code(tree, feature_names, node, tree_index, outfolder):
	with open(outfolder + "tree" + str(tree_index) + "code.go", "w") as code_file:
		def subtree_to_go_code(node, depth):
			indent = "\t" * depth
			if tree.feature[node] != _tree.TREE_UNDEFINED:
				feature = feature_names[tree.feature[node]]
				threshold = tree.threshold[node]
				code_file.write("{}if s.{} <= {} {{\n".format(indent, feature, math.floor(threshold)))
				subtree_to_go_code(tree.children_left[node], depth + 1)
				code_file.write("{}}}\n".format(indent))
				subtree_to_go_code(tree.children_right[node], depth)
			else:
				code_file.write("{}return {}\n".format(indent, tree.value[node][0][0]))

		code_file.write("// Package ab (Code generated by a Python script)\n")
		code_file.write("package ab\n\n")
		code_file.write("func (s Sample) getEstimatedValue() float64 {\n")
		subtree_to_go_code(0, 1)
		code_file.write("}\n")

def main(argv):
	# Read flags
	datafile = "sample_data/data.in"
	outfolder = "./"
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
	df = pd.read_csv(datafile, comment = "#",
					dtype = {"value": np.float64, "num_stones": np.uint8,
							"occ_red_rows": np.uint8, "occ_red_cols": np.uint8,
							"occ_blue_rows": np.uint8, "occ_blue_cols": np.uint8,
							"red_p0": np.uint8, "blue_p0": np.uint8,
							"red_p1": np.uint8, "blue_p1": np.uint8,
							"red_p2": np.uint8, "blue_p2": np.uint8,
							"lp": bool})

	y = df["value"]
	X = df.drop(columns = ["value"])

	feature_names = X.columns

	# Split
	X_train, X_test, y_train, y_test = train_test_split(X, y, test_size = 0.1, random_state = 4224)

	# Create Regressors
	dtrs = [DecisionTreeRegressor(max_depth =    2, min_samples_leaf = 5),
			DecisionTreeRegressor(max_depth =    5, min_samples_leaf = 5),
			DecisionTreeRegressor(max_depth =   10, min_samples_leaf = 5),
			DecisionTreeRegressor(max_depth = None, min_samples_leaf = 50),
			DecisionTreeRegressor(max_depth = None, min_samples_leaf = 1)]

	# Create a plot
	plt.figure(figsize=(10,6))
	plt.plot(y_test.tolist(), label = "actual values", linewidth = 0.7)
	plt.ylim(-1.2, 1.2)

	# Create file for statistics
	with open(outfolder + "stats.txt", "w") as stats_file:
		for dtri in range(len(dtrs)):
			dtr = dtrs[dtri]

			# Train
			dtr.fit(X_train, y_train)

			# Predict
			y1 = dtr.predict(X_test)

			# Print statistics
			stats_file.write("#############################################\n")
			stats_file.write("Statistics for:" + str(dtr) + "\n")
			stats_file.write("Feature importances:\n")
			fi = zip(X.keys(), dtr.feature_importances_)
			for (k, v) in fi:
				stats_file.write("\t" + k + ": " + str(v) + "\n")
			sc = dtr.score(X_test, y_test)
			stats_file.write("SCORE: " + str(sc) + "\n")

			tree.export_graphviz(dtr, out_file = outfolder + "tree" + str(dtri) + ".dot",
								feature_names = feature_names)

			# Add to plot
			plt.plot(y1, label = "predicted values (max_depth=" + str(dtr.max_depth) + ")",
					linewidth = 0.7)

			# Output Go code
			tree_to_go_code(dtr.tree_, feature_names, 0, dtri, outfolder)

	# Create file for Go sample struct
	with open(outfolder + "sample.go", "w") as sample_file:
		sample_file.write("// Package ab (Code generated by a Python script)\n")
		sample_file.write("package ab\n\n")
		sample_file.write("type Sample struct {\n\t")
		sample_file.write(feature_names[0])
		for f in feature_names[1:]:
			sample_file.write(", " + f)
		sample_file.write(" int\n")
		sample_file.write("}")

	plt.xlabel("Samples")
	plt.ylabel("Value")
	plt.title("Comparison of predicted and actual values")
	plt.legend(fontsize = "x-small")
	plt.savefig(outfolder + "plot.pdf")
	plt.close()

if __name__ == "__main__":
	main(sys.argv[1:])
