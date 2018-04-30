package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/RdecKa/mcts/hex"
	"github.com/RdecKa/mcts/mcts"
)

func main() {
	boardSize := flag.Int("size", 3, "Board size")
	numIterations := flag.Int("iter", 10000, "Number of iterations")
	indentJSON := flag.Bool("indent", false, "Indent JSON output")
	flag.Parse()
	fmt.Printf("Using boardSize = %d, numIterations = %d\n", *boardSize, *numIterations)

	initState := hex.NewState(byte(*boardSize))
	explorationFactor := math.Sqrt(2)
	mc := mcts.InitMCTS(*initState, explorationFactor)

	for i := 0; i < *numIterations; i++ {
		mc.RunIteration()
	}

	err := mcts.WriteToFile(*mc, "./out", *indentJSON)
	if err != nil {
		fmt.Println(err)
	}
}
