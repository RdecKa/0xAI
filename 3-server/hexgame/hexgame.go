package hexgame

import (
	"fmt"

	"github.com/RdecKa/bachleor-thesis/3-server/hexplayer"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

func playOneGame(players [2]hexplayer.HexPlayer, startingPlayer int) error {
	// Init game
	boardSize := 7
	for p := 0; p < 2; p++ {
		err := players[p].InitGame(boardSize, players[startingPlayer].GetColor())
		if err != nil {
			return err
		}
	}
	state := hex.NewState(byte(boardSize), players[startingPlayer].GetColor())
	turn := startingPlayer
	var prevAction *hex.Action

	// Play game
	for g, _ := state.IsGoalState(true); !g; g, _ = state.IsGoalState(true) {
		nextAction, err := players[turn].NextAction(prevAction)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if nextAction == nil {
			// Player has resigned
			fmt.Printf("Player %d resigned!\n", turn+1)
			break
		}
		s := state.GetSuccessorState(nextAction).(hex.State)
		state = &s
		prevAction = nextAction
		turn = 1 - turn
		fmt.Printf("%v", state)
	}

	// Game results
	for p := 0; p < 2; p++ {
		w := false
		if 1-turn == p {
			w = true
		}
		players[p].EndGame(prevAction, w)
	}

	fmt.Printf("%v", state)
	return nil
}

func playNGames(players [2]hexplayer.HexPlayer, numGames int) [2]int {
	startingPlayer := 0
	for g := 0; g < numGames; g++ {
		err := playOneGame(players, startingPlayer)
		if err != nil {
			fmt.Println("Game canceled: " + err.Error())
			continue
		}

		fmt.Printf("Results after %d games:\n", g+1)
		fmt.Printf("\tPlayer one: %d\n", players[0].GetNumberOfWins())
		fmt.Printf("\tPlayer two: %d\n", players[1].GetNumberOfWins())

		// Switch roles
		startingPlayer = 1 - startingPlayer
	}
	return [2]int{players[0].GetNumberOfWins(), players[1].GetNumberOfWins()}
}

// Play accepts an array of two players and number of games to be played. It
// runs numGames games of Hex between the given players.
func Play(players [2]hexplayer.HexPlayer, numGames int) {
	results := playNGames(players, numGames)
	fmt.Printf("*** Final results ***:\n")
	fmt.Printf("\tPlayer one: %d\n", results[0])
	fmt.Printf("\tPlayer two: %d\n", results[1])
}
