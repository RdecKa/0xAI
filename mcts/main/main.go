package main

import (
	"fmt"
	"math"

	"github.com/RdecKa/mcts/hex"
	"github.com/RdecKa/mcts/mcts"
)

func main() {
	initState := hex.NewState(2)
	explorationFactor := math.Sqrt(2)
	mc := mcts.InitMCTS(*initState, explorationFactor)
	for i := 0; i < 30000; i++ {
		mc.RunIteration()
	}
	err := mcts.WriteToFile(*mc, "./out")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mc)
}
