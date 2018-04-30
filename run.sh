#!/bin/bash

output_folder="data/mcts/"
visual_data_folder="visual/mcts/"
visual_data_json_file="${visual_data_folder}data.js"
visual_html_index="${visual_data_folder}index.html"

go run mcts/main/main.go -output=$output_folder -indent=false -iter=10000 -size=3

# Get the newest file in output directory
data_file_name=$(ls -t $output_folder | head -1)

echo -n "let mcst_json = " > $visual_data_json_file
cat "${output_folder}${data_file_name}" >> $visual_data_json_file

# Open in browser
xdg-open $visual_html_index
