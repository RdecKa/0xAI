SHELL = /bin/sh

# ---> Go variables <---
GO_COMMAND = go
GO_INSTALL = $(GO_COMMAND) install
GO_CLEAN = $(GO_COMMAND) clean -i
GO_CLEAN_FILES = github.com/RdecKa/bachleor-thesis/1-mcts/mcts \
	github.com/RdecKa/bachleor-thesis/3-ab/ab \
	github.com/RdecKa/bachleor-thesis/common/astarsearch \
	github.com/RdecKa/bachleor-thesis/common/game \
	github.com/RdecKa/bachleor-thesis/common/pq \
	github.com/RdecKa/bachleor-thesis/common/tree \
	github.com/RdecKa/bachleor-thesis/common/game/hex \
	github.com/RdecKa/bachleor-thesis/server/hexgame \
	github.com/RdecKa/bachleor-thesis/server/hexplayer

# ---> Python variables <---
PYTHON_COMMAND = python3

# ---> CSS variables <---
CSS_COMPILER = sass

# ---> Common variables <---
START_TIME := $(shell date +"%Y%m%dT%H%M%S")
PATTERNS_FILE = common/game/hex/patterns.txt
OPEN_IN_BROWSER = xdg-open
SIZE = 7
OUT_DATA_DIR = data/$(SIZE)/

# ---> MCTS variables <---
TIME = 10
WORKERS = 3
TREASHOLD_N = 500
MCTS_DIR = 1-mcts/
MCTS_FILES := $(shell find $(MCTS_DIR) -type f -name "*.go")
MCTS_MAIN = $(MCTS_DIR)main/main.go
MCTS_OUT_DIR_PARENT = $(OUT_DATA_DIR)mcts/
MCTS_OUT_DIR = $(MCTS_OUT_DIR_PARENT)run-$(START_TIME)/

# ---> Visual variables <---
VISUAL_DATA_DIR = visual/mcts/
VISUAL_DATA_JSON_FILE = $(VISUAL_DATA_DIR)data.js
VISUAL_HTML_INDEX = $(VISUAL_DATA_DIR)index.html
INDENT = false
JSON = false

# ---> ML variables <---
ML_DIR = 2-ml/
ML_OUT_DIR = $(OUT_DATA_DIR)ml/ml-$(START_TIME)/
ML_MERGE_DATA_FILE = $(ML_OUT_DIR)data.in
ML_INPUT_FILES = $(shell find $(MCTS_OUT_DIR) -type f -name "*.in")
ML_MAIN = $(ML_DIR)learn.py
ML_DOT_FILES = $(shell find $(ML_OUT_DIR) -type f -name "*.dot")
ML_PS_FILES = $(ML_DOT_FILES:.dot=.ps)
ML_SELECT_TREE = 2
ML_SELECT_TREE_FILE = $(ML_OUT_DIR)tree$(ML_SELECT_TREE)code.go
ML_LINEAR_REGRESSION_FILE = $(ML_OUT_DIR)linear0code.go
ML_GEN_SAMPLE_FILE = $(ML_OUT_DIR)sample.go

# ---> AB variables <---
AB_DIR = 3-ab/ab/
AB_GEN_SAMP_FILE = $(AB_DIR)sample.go
AB_GEN_TREE_FILE = $(AB_DIR)treecode.go
AB_GEN_LINEAR_FILE = $(AB_DIR)linearcode.go

# ---> Server variables <---
SERV_DIR = server/
SERV_MAIN = $(SERV_DIR)main/server.go
SERV_BIN_NAME = server

################################################################################

all: mcts ml serv

clean:
	$(GO_CLEAN) $(GO_CLEAN_FILES)
	rm -f $(AB_GEN_TREE_FILE) $(AB_GEN_SAMP_FILE)
	rm -f $(VISUAL_DATA_DIR)style.css*
	rm -f $(SERV_DIR)static/css/style.css*

%.css: %.scss
	$(CSS_COMPILER) $< $@

# ---> MCTS targets <---
mctscomp: $(MCTS_FILES)
	# --> Compile the MCTS program <--
	$(GO_INSTALL) $(MCTS_MAIN)

mctsrun:
# Make a new folder for output files
	mkdir -p "$(MCTS_OUT_DIR)"

	# --> Run MCTS program <--
	main -output=$(MCTS_OUT_DIR) -json=$(JSON) -indent=$(INDENT) -time=$(TIME) -size=$(SIZE) -workers=$(WORKERS) -patterns=$(PATTERNS_FILE) -treasholdn=$(TREASHOLD_N)

mctsjson: DATA_FILE = "$(shell ls $(MCTS_OUT_DIR)*.json)"
mctsjson:
	# --> Create JSON <--
# File to be used for MCTS visualisation: $(VISUAL_DATA_JSON_FILE)
	echo -n "let mcst_json = " > $(VISUAL_DATA_JSON_FILE)
# Find the JSON file in output directory, copy its content: $(DATA_FILE)
	cat $(DATA_FILE) >> $(VISUAL_DATA_JSON_FILE)

mctsvisual: $(VISUAL_DATA_DIR)style.css
	# --> Open results in browser <--
	$(OPEN_IN_BROWSER) $(VISUAL_HTML_INDEX)

mcts: mctscomp mctsrun

mctsall: JSON = true
mctsall: mctscomp mctsrun mctsjson mctsvisual


# ---> ML targets <---
mlcreatedir:
	mkdir -p $(ML_OUT_DIR)

mlmerge:
	# --> Create a file to merge all learning samples: $(ML_MERGE_DATA_FILE) <--
# Copy attribute names
	head -n 1 $(word 1, $(ML_INPUT_FILES)) > $(ML_MERGE_DATA_FILE)
# Copy data
	tail -n +2 $(ML_INPUT_FILES) >> $(ML_MERGE_DATA_FILE)
# Remove redundant lines
	sed -i '/==>\|^$$/d' $(ML_MERGE_DATA_FILE)

mlrun: mlcreatedir mlmerge
	# --> Run ML program <--
	$(PYTHON_COMMAND) $(ML_MAIN) -d $(ML_MERGE_DATA_FILE) -o $(ML_OUT_DIR) -a

mltrees:
	# --> Visualize trees <--
	$(foreach FILE, $(ML_DOT_FILES), $(shell dot -Tps $(FILE) -o $(FILE:.dot=.ps)))

mlcopycode:
	# --> Copy generated Go files to AB directory <--
	cp -f "$(ML_SELECT_TREE_FILE)" "$(AB_GEN_TREE_FILE)"
	cp -f "$(ML_GEN_SAMPLE_FILE)" "$(AB_GEN_SAMP_FILE)"
	cp -f "$(ML_LINEAR_REGRESSION_FILE)" "$(AB_GEN_LINEAR_FILE)"

ml: mlrun mlcopycode

mlall: mlrun mltrees mlcopycode


# ---> Server targets <---
servcomp: $(SERV_DIR)static/css/style.css
	# --> Compile server <--
	$(GO_INSTALL) $(SERV_MAIN)

servrun:
	# --> Start server by typing '$(SERV_BIN_NAME)' <--

serv: servcomp servrun
