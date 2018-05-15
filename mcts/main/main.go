package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/RdecKa/mcts/hex"
	"github.com/RdecKa/mcts/mcts"
	"github.com/RdecKa/mcts/tree"
)

func main() {
	boardSize := flag.Int("size", 3, "Board size")
	numIterations := flag.Int("iter", 10000, "Number of iterations")
	indentJSON := flag.Bool("indent", false, "Indent JSON output")
	output := flag.String("output", ".", "Output folder")
	flag.Parse()
	fmt.Printf("Using boardSize = %d, numIterations = %d\n", *boardSize, *numIterations)

	initState := hex.NewState(byte(*boardSize))
	explorationFactor := math.Sqrt(2)
	minBeforeExpand := uint(10)
	mc := mcts.InitMCTS(*initState, explorationFactor, minBeforeExpand)
	root := mc

	toBeExpanded := make([]*tree.Node, 1, 100)             // First element is nil, will be popped in the first iteration
	probExpand := 1.0 / float64((*boardSize)*(*boardSize)) // Initial probability of continuing MCTS from a leaf node

	for ok := true; ok; ok = len(toBeExpanded) > 0 {
		if toBeExpanded[0] != nil {
			mc = mc.ContinueMCTSFromNode(toBeExpanded[0])
		} // else it is the first round, MCTS is already initialised
		toBeExpanded = toBeExpanded[1:] // Pop the queue

		fmt.Printf("Starting new MCTS from node\n%v\n", mc.GetInitialNode())
		fmt.Printf("Queue size: %d\n", len(toBeExpanded))
	for i := 0; i < *numIterations; i++ {
		if i > 0 && i%10000 == 0 {
			fmt.Printf("Finished iteration %d\n", i)
		}
		mc.RunIteration()
	}

		// Write input-output pairs for supervised machine learning, generate
		// new nodes tocontinue MCTS
	filePrefix := fmt.Sprintf("sample_%02d_%d", *boardSize, *numIterations)
		expandCandidates, err := mc.GenSamples(*output, filePrefix, 100)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

		for _, c := range expandCandidates {
			if rand.Float64() <= probExpand {
				toBeExpanded = append(toBeExpanded, c)
			}
		}

		probExpand = 1.0 / float64(len(toBeExpanded)*len(toBeExpanded)+1)
	}

	// Write JSON
	filePrefix := fmt.Sprintf("out_%02d_%d", *boardSize, *numIterations)
	err := mcts.WriteToFile(*root, *output, filePrefix, *indentJSON)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
