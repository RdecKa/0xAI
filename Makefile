SHELL = /bin/sh

# ---> Go <---
GO_COMMAND = go
GO_INSTALL = $(GO_COMMAND) install

# ---> Common <---
NOW := $(shell date +"%Y%m%dT%H%M%S")

# ---> MCTS <---
MCTS_SRC = 1-mcts/
MCTS_FILES := $(shell find $(MCTS_SRC) -type f -name "*.go")
MCTS_MAIN = $(MCTS_SRC)main/main.go
MCTS_OUT_FOLD_PARENT = data/mcts/
MCTS_OUT_FOLD = $(MCTS_OUT_FOLD_PARENT)run-$(NOW)/
TIME = 5
SIZE = 3
WORKERS = 3

# ---> Visual <---
VISUAL_DATA_FOLD = visual/mcts/
VISUAL_DATA_JSON_FILE = $(VISUAL_DATA_FOLD)data.js
VISUAL_HTML_INDEX = $(VISUAL_DATA_FOLD)index.html
INDENT = false
JSON = false

# ---> AB <---
PATTERNS_FILE = common/game/hex/patterns.txt
AB_FOLD = 3-ab/ab/

mctscomp: $(MCTS_FILES)
	# Compile the MCTS program
	$(GO_INSTALL) $(MCTS_MAIN)

mctsrun:
	# Make a new folder for output files: $(MCTS_OUT_FOLD)
	$(shell mkdir "$(MCTS_OUT_FOLD)")

	# Run the program
	main -output=$(MCTS_OUT_FOLD) -json=$(JSON) -indent=$(INDENT) -time=$(TIME) -size=$(SIZE) -workers=$(WORKERS) -patterns=$(PATTERNS_FILE)

mctsjson: DATA_FILE = "$(shell ls $(MCTS_OUT_FOLD)*.json)"
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
