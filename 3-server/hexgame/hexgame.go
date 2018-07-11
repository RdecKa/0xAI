package hexgame

import (
	"fmt"

	"github.com/RdecKa/bachleor-thesis/3-server/hexplayer"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

func playOneGame(players [2]hexplayer.HexPlayer) (int, error) {
	// Init game
	boardSize := 5
	for p := 0; p < 2; p++ {
		err := players[p].InitGame(boardSize)
		if err != nil {
			return -1, err
		}
	}
	state := hex.NewState(byte(boardSize))
	turn := 0
	var prevAction *hex.Action

	// Play game
	for !state.IsGoalState() {
		nextAction, err := players[turn].NextAction(prevAction)
		if err != nil {
			fmt.Println(err)
			return -1, err
		}
		s := state.GetSuccessorState(nextAction).(hex.State)
		state = &s
		prevAction = nextAction
		turn = 1 - turn
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
	return 1 - turn, nil
}

func playNGames(players [2]hexplayer.HexPlayer, numGames int) [2]int {
	wins := [2]int{0, 0}
	for g := 0; g < numGames; g++ {
		winner, err := playOneGame(players)
		if err != nil {
			fmt.Println("Game canceled: " + err.Error())
			continue
		}
		wins[winner]++
	}
	return wins
}

// Play accepts an array of two players and number of games to be played. It
// runs numGames games of Hex between the given players.
func Play(players [2]hexplayer.HexPlayer, numGames int) {
	results := playNGames(players, numGames)
	fmt.Println("Results:")
	fmt.Printf("\tPlayer one: %d\n", results[0])
	fmt.Printf("\tPlayer two: %d\n", results[1])
}
