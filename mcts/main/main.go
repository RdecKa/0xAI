package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/RdecKa/mcts/hex"
	"github.com/RdecKa/mcts/mcts"
)

func main() {
	boardSize := flag.Int("size", 3, "Board size")
	numIterations := flag.Int("iter", 10000, "Number of iterations")
	indentJSON := flag.Bool("indent", false, "Indent JSON output")
	output := flag.String("output", ".", "Output file")
	flag.Parse()
	fmt.Printf("Using boardSize = %d, numIterations = %d\n", *boardSize, *numIterations)

	initState := hex.NewState(byte(*boardSize))
	explorationFactor := math.Sqrt(2)
	minBeforeExpand := uint(10)
	mc := mcts.InitMCTS(*initState, explorationFactor, minBeforeExpand)

	for i := 0; i < *numIterations; i++ {
		if i > 0 && i%10000 == 0 {
			fmt.Printf("Finished iteration %d\n", i)
		}
		mc.RunIteration()
	}

	filePrefix := fmt.Sprintf("out_%02d_%d", *boardSize, *numIterations)
	err := mcts.WriteToFile(*mc, *output, filePrefix, *indentJSON)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
