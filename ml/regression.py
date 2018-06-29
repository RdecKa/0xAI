import sys, getopt
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

from sklearn import tree
from sklearn.tree import DecisionTreeRegressor
from sklearn.model_selection import train_test_split


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

	# Create file for statistics
	with open(outfolder + "stats.txt", "w") as stats_file:

		# Read data from file
		print("Reading data from file:", datafile)
		df = pd.read_csv(datafile, comment = "#",
						dtype = {"final_result": np.int8, "num_stones": np.uint8,
								"occ_red_rows": np.uint8, "occ_red_cols": np.uint8,
								"occ_blue_rows": np.uint8, "occ_blue_cols": np.uint8,
								"red_p1": np.uint8, "blue_p1": np.uint8,
								"red_p2": np.uint8, "blue_p2": np.uint8})

		#print(df)
		#print("KEYS", df.keys())
		#print("COLUMNS", df.columns)
		#print("SHAPE:", df.shape)
		#print("TYPES", df.dtypes)

		y = df["value"]
		X = df.drop(columns = ["value"])

		# Split
		X_train, X_test, y_train, y_test = train_test_split(X, y, test_size = 0.1, random_state = 4224)

		# Create Regressors
		dtrs = [DecisionTreeRegressor(max_depth =  2, min_samples_leaf = 5),
				DecisionTreeRegressor(max_depth =  5, min_samples_leaf = 5),
				DecisionTreeRegressor(max_depth = 10, min_samples_leaf = 5)]

		# Create a plot
		plt.figure(figsize=(10,6))
		plt.plot(y_test.tolist(), label = "actual values", linewidth = 0.7)
		plt.ylim(-1.2, 1.2)

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
			tree.export_graphviz(dtr, out_file = outfolder + "tree" + str(dtri) + ".dot")

			# Add to plot
			plt.plot(y1, label = "predicted values (max_depth=" + str(dtr.max_depth) + ")",
					linewidth = 0.7)

	plt.xlabel("Samples")
	plt.ylabel("Value")
	plt.title("Comparison of predicted and actual values")
	plt.legend(fontsize = "x-small")
	plt.savefig(outfolder + "plot.pdf")
	plt.close()

if __name__ == "__main__":
	main(sys.argv[1:])
