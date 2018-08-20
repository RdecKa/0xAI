package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/RdecKa/bachleor-thesis/1-mcts/mcts"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

func main() {
	// Read flags
	pBoardSize := flag.Int("size", 3, "Board size")
	pSecondsToRun := flag.Int("time", 10000, "Time to run")
	pWriteJSON := flag.Bool("json", false, "Output JSON file")
	pIndentJSON := flag.Bool("indent", false, "Indent JSON output")
	pOutputFolder := flag.String("output", ".", "Output folder")
	pNumWorkers := flag.Int("workers", 2, "Number of goroutines to run in parallel")
	pPatternsFile := flag.String("patterns", ".", "File with hex patterns")
	flag.Parse()
	boardSize, secondsToRun, numWorkers, patternsFile := *pBoardSize, *pSecondsToRun, *pNumWorkers, *pPatternsFile
	writeJSON, indentJSON, outputFolder := *pWriteJSON, *pIndentJSON, *pOutputFolder

	fmt.Printf("Using boardSize = %d, secondsToRun = %d, numWorkers = %d\n", boardSize, secondsToRun, numWorkers)

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
	mcts.RunMCTSinParallel(numWorkers, boardSize, time.Duration(secondsToRun)*time.Second, outputFolder, patternsFile, mc)

	if writeJSON {
		// Write JSON
		filePrefix := fmt.Sprintf("out_%02d_%d", boardSize, secondsToRun)
		err := mcts.WriteToFile(*root, outputFolder, filePrefix, indentJSON)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
