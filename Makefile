SHELL = /bin/sh

# ---> Go variables <---
GO_COMMAND = go
GO_INSTALL = $(GO_COMMAND) install

# ---> Python variables <---
PYTHON_COMMAND = python3

# ---> Common variables <---
START_TIME := $(shell date +"%Y%m%dT%H%M%S")
PATTERNS_FILE = common/game/hex/patterns.txt

# ---> MCTS variables <---
MCTS_DIR = 1-mcts/
MCTS_FILES := $(shell find $(MCTS_DIR) -type f -name "*.go")
MCTS_MAIN = $(MCTS_DIR)main/main.go
MCTS_OUT_DIR_PARENT = data/mcts/
MCTS_OUT_DIR = $(MCTS_OUT_DIR_PARENT)run-$(START_TIME)/
TIME = 5
SIZE = 3
WORKERS = 3

# ---> Visual variables <---
VISUAL_DATA_DIR = visual/mcts/
VISUAL_DATA_JSON_FILE = $(VISUAL_DATA_DIR)data.js
VISUAL_HTML_INDEX = $(VISUAL_DATA_DIR)index.html
INDENT = false
JSON = false

# ---> ML variables <---
ML_DIR = 2-ml/
ML_OUT_DIR = data/ml/ml-$(START_TIME)/
ML_MERGE_DATA_FILE = $(ML_OUT_DIR)data.in
ML_INPUT_FILES = $(shell find $(MCTS_OUT_DIR) -type f -name "*.in")
ML_MAIN = $(ML_DIR)regression.py
ML_DOT_FILES = $(shell find $(ML_OUT_DIR) -type f -name "*.dot")
ML_PS_FILES = $(ML_DOT_FILES:.dot=.ps)
ML_SELECT_TREE = 2
ML_SELECT_TREE_FILE = $(ML_OUT_DIR)tree$(ML_SELECT_TREE)code.go
ML_GEN_SAMPLE_FILE = $(ML_OUT_DIR)sample.go

# ---> AB variables <---
AB_DIR = 3-ab/ab/
AB_GEN_TREE_FILE = $(AB_DIR)treecode.go

# ---> MCTS targets <---
mctscomp: $(MCTS_FILES)
	# Compile the MCTS program
	$(GO_INSTALL) $(MCTS_MAIN)

mctsrun:
	# Make a new folder for output files: $(MCTS_OUT_DIR)
	mkdir "$(MCTS_OUT_DIR)"

	# Run the program
	main -output=$(MCTS_OUT_DIR) -json=$(JSON) -indent=$(INDENT) -time=$(TIME) -size=$(SIZE) -workers=$(WORKERS) -patterns=$(PATTERNS_FILE)

mctsjson: DATA_FILE = "$(shell ls $(MCTS_OUT_DIR)*.json)"
mctsjson:
	# File to be used for MCTS visualisation: $(VISUAL_DATA_JSON_FILE)
	echo -n "let mcst_json = " > $(VISUAL_DATA_JSON_FILE)
	# Find the JSON file in output directory, copy its content: $(DATA_FILE)
	cat $(DATA_FILE) >> $(VISUAL_DATA_JSON_FILE)

mctsvisual:
	# Open results in browser
	xdg-open $(VISUAL_HTML_INDEX)

mcts: mctscomp mctsrun

mctsall: JSON = true
mctsall: mctscomp mctsrun mctsjson mctsvisual


# ---> ML targets <---
mlcreatedir:
# Create a folder for ML output
	mkdir $(ML_OUT_DIR)

$(ML_MERGE_DATA_FILE): $(ML_INPUT_FILES)
	# Create a file to merge all learning samples: $(ML_MERGE_DATA_FILE)
	# Copy attribute names
	head -n 1 $< > $@
	# Copy data
	tail -n +2 $(ML_INPUT_FILES) >> $@
	# Remove redundant lines
	sed -i '/==>\|^$$/d' $(ML_MERGE_DATA_FILE)

mlrun: mlcreatedir $(ML_MERGE_DATA_FILE)
	$(PYTHON_COMMAND) $(ML_MAIN) -d $(ML_MERGE_DATA_FILE) -o $(ML_OUT_DIR)

mltrees:
	# Visualize trees
	$(foreach FILE, $(ML_DOT_FILES), $(shell dot -Tps $(FILE) -o $(FILE:.dot=.ps)))

mlcopycode:
	cp -f "$(ML_SELECT_TREE_FILE)" "$(AB_GEN_TREE_FILE)"
	cp -f "$(ML_GEN_SAMPLE_FILE)" "$(AB_DIR)"

ml: mlrun mlcopycode

mlall: mlrun mltrees mlcopycode
