#!/bin/bash

output_folder_mcts="data/mcts/"
visual_data_folder="visual/mcts/"
visual_data_json_file="${visual_data_folder}data.js"
visual_html_index="${visual_data_folder}index.html"

indent=false
iter=10000
size=3
browser=false
workers=2

now=$(date +"%Y%m%dT%H%M%S")
output_folder_mcts="${output_folder_mcts}run-${now}/"
# Make a new folder for output files
mkdir "$output_folder_mcts"

while getopts 'bin:o:s:w:' flag; do
	case "${flag}" in
		b) browser=true ;;
		i) indent='true' ;;
		n) iter="${OPTARG}" ;;
		o) output_folder_mcts="${OPTARG}" ;;
		s) size="${OPTARG}" ;;
		w) workers="${OPTARG}" ;;
		*) echo "Unexpected option ${flag}" ;;
	esac
done

# Run the program
go run 1-mcts/main/main.go -output="$output_folder_mcts" -indent="$indent" -iter="$iter" -size="$size" -workers="$workers"
status=$?
if [ "$status" -ne 0 ]; then
	echo "Error occured."
	exit 1
fi

# Get the newest file in output directory
data_file_name=$(ls -t $output_folder_mcts | head -1)

echo -n "let mcst_json = " > $visual_data_json_file
cat "${output_folder_mcts}${data_file_name}" >> $visual_data_json_file

# Open results in browser
if [ "$browser" = true ]; then
	xdg-open $visual_html_index
fi

echo "MCTS phase completed."

# Ask user whether to continue with the next phase
answer=""
while [[ "$answer" != "y" ]] && [ "$answer" != "n" ]; do
	read -p "Do you want to continue with ML? [y/n] " answer
done

if [ "$answer" = "n" ]; then
	echo "Bye!"
	exit 0
fi

echo "OK, let's continue!"

# Create a new folder for ML outputs
output_folder_ml="data/ml/ml-${now}/"
mkdir "$output_folder_ml"

# Create a file to merge all learning samples
data_file="${output_folder_ml}data.in"
touch "$data_file"

# Merge all samples
head -n 1 ${output_folder_mcts}*_0.in >> $data_file # Copy attribute names
for filename in ${output_folder_mcts}*.in; do
	tail -n +2 "$filename" >> "$data_file" # Copy everything except attribute names
done

python3 2-ml/regression.py -d "$data_file" -o "$output_folder_ml"

# Visualize trees
for filename in ${output_folder_ml}*.dot; do
	dot -Tps "$filename" -o "${filename%.*}.ps"
done
