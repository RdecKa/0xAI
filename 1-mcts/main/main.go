package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/RdecKa/bachleor-thesis/1-mcts/mcts"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

func main() {
	// Read flags
	pBoardSize := flag.Int("size", 3, "Board size")
	pNumIterations := flag.Int("iter", 10000, "Number of iterations")
	pIndentJSON := flag.Bool("indent", false, "Indent JSON output")
	pOutputFolder := flag.String("output", ".", "Output folder")
	pNumWorkers := flag.Int("workers", 2, "Number of goroutines to run in parallel")
	flag.Parse()
	boardSize, numIterations, indentJSON, outputFolder, numWorkers := *pBoardSize, *pNumIterations, *pIndentJSON, *pOutputFolder, *pNumWorkers
	fmt.Printf("Using boardSize = %d, numIterations = %d, numWorkers = %d\n", boardSize, numIterations, numWorkers)

	// Init the algorithm
	initState := hex.NewState(byte(boardSize))
	explorationFactor := math.Sqrt(2)
	minBeforeExpand := uint(10)
	mc := mcts.InitMCTS(*initState, explorationFactor, minBeforeExpand)
	root := mc

	// Run the algorithm
	mcts.RunMCTSinParallel(numWorkers, boardSize, numIterations, outputFolder, mc)

	// Write JSON
	filePrefix := fmt.Sprintf("out_%02d_%d", boardSize, numIterations)
	err := mcts.WriteToFile(*root, outputFolder, filePrefix, indentJSON)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
