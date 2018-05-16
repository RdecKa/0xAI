#!/bin/bash

output_folder="data/mcts/"
visual_data_folder="visual/mcts/"
visual_data_json_file="${visual_data_folder}data.js"
visual_html_index="${visual_data_folder}index.html"

indent=false
iter=10000
size=3
browser=false
workers=2

now=$(date +"%Y%m%dT%I%M%S")
output_folder="${output_folder}run-${now}/"
# Make a new folder for output files
mkdir "$output_folder"

while getopts 'bin:o:s:w:' flag; do
	case "${flag}" in
		b) browser=true ;;
		i) indent='true' ;;
		n) iter="${OPTARG}" ;;
		o) output_folder="${OPTARG}" ;;
		s) size="${OPTARG}" ;;
		w) workers="${OPTARG}" ;;
		*) echo "Unexpected option ${flag}" ;;
	esac
done

# Run the program
go run mcts/main/main.go -output="$output_folder" -indent="$indent" -iter="$iter" -size="$size" -workers="$workers"
status=$?
if [ "$status" -ne 0 ]; then
	echo "Error occured."
	exit 1
fi

# Get the newest file in output directory
data_file_name=$(ls -t $output_folder | head -1)

echo -n "let mcst_json = " > $visual_data_json_file
cat "${output_folder}${data_file_name}" >> $visual_data_json_file

# Open results in browser
if [ "$browser" = true ]; then
	xdg-open $visual_html_index
fi
