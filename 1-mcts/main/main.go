package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/RdecKa/0xAI/1-mcts/mcts"
	"github.com/RdecKa/0xAI/common/game/hex"
)

func main() {
	// Read flags
	pBoardSize := flag.Int("size", 3, "Board size")
	pSecondsToRun := flag.Int("time", 5, "Seconds to run")
	pThresholdN := flag.Uint("thresholdn", 100, "Number of visits of a node required to generate a sample")
	pWriteJSON := flag.Bool("json", false, "Output JSON file")
	pIndentJSON := flag.Bool("indent", false, "Indent JSON output")
	pOutputFolder := flag.String("output", "./", "Output folder")
	pNumWorkers := flag.Int("workers", 3, "Number of goroutines to run in parallel")
	pPatternsFile := flag.String("patterns", "patterns.txt", "File with hex patterns")
	flag.Parse()
	boardSize, secondsToRun, thresholdN, numWorkers, patternsFile := *pBoardSize, *pSecondsToRun, *pThresholdN, *pNumWorkers, *pPatternsFile
	writeJSON, indentJSON, outputFolder := *pWriteJSON, *pIndentJSON, *pOutputFolder

	fmt.Printf("Using boardSize = %d, secondsToRun = %d, numWorkers = %d, patternsFile = %s, writeJSON = %t, indentJSON = %t, outputFolder = %s, thresholdN = %d\n",
		boardSize, secondsToRun, numWorkers, patternsFile, writeJSON, indentJSON, outputFolder, thresholdN)

	// Init the algorithm
	initState := hex.NewState(byte(boardSize), hex.Red)
	explorationFactor := math.Sqrt(2)
	minBeforeExpand := uint(10)
	mc := mcts.InitMCTS(*initState, explorationFactor, minBeforeExpand)
	var root *mcts.MCTS
	if writeJSON {
		root = mc
	}

	// Run the algorithm
	mcts.RunMCTSinParallel(numWorkers, boardSize, thresholdN, time.Duration(secondsToRun)*time.Second,
		outputFolder, patternsFile, mc, false)

	if writeJSON {
		// Write JSON
		filePrefix := fmt.Sprintf("out_%02d_%d", boardSize, secondsToRun)
		err := mcts.WriteToFile(*root, outputFolder, filePrefix, indentJSON)
		if err != nil {
			panic(err)
		}
	}
}
