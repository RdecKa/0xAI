import sys, getopt
import pandas as pd

def main(argv):
	datafile = 'sample_data/sample_06_5000_0.in'
	try:
		opts, args = getopt.getopt(argv, "d:")
	except getopt.GetoptError:
		print("Error parsing the command line arguments")
		sys.exit(1)
	for o, a in opts:
		if o == "-d":
			datafile = a

	print("Reading data from file:", datafile)

	data = pd.read_csv(datafile, comment="#")

	print(data)

if __name__ == "__main__":
	main(sys.argv[1:])
