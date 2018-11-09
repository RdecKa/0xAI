package hexgame

import (
	"fmt"
	"runtime"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
	"github.com/RdecKa/bachleor-thesis/server/hexplayer"
	"github.com/gorilla/websocket"
)

// playOneGame returns 0 if the first player (of players) won and 1 if the
// second player won
func playOneGame(boardSize int, players [2]hexplayer.HexPlayer, passiveClient hexplayer.HexPlayer, startingPlayer int) (int, error) {
	// Init game
	for p := 0; p < 3; p++ {
		var err error
		if p < 2 {
			err = players[p].InitGame(boardSize, players[startingPlayer].GetColor())
		} else if p == 2 && passiveClient != nil {
			err = passiveClient.InitGame(boardSize, players[startingPlayer].GetColor())
		}

		if err != nil {
			return -1, err
		}
	}
	state := hex.NewState(byte(boardSize), players[startingPlayer].GetColor())
	turn := startingPlayer
	var prevAction *hex.Action

	// Play game
	for g, _ := state.IsGoalState(true); !g; g, _ = state.IsGoalState(true) {
		players[turn].PrevAction(prevAction)
		nextAction, err := players[turn].NextAction()
		if err != nil {
			fmt.Println(err)
			return -1, err
		}
		if nextAction == nil {
			// Player has resigned
			fmt.Printf("Player %v resigned!\n", players[turn].GetColor())
			break
		}
		if passiveClient != nil {
			passiveClient.PrevAction(nextAction)
		}
		s := state.GetSuccessorState(nextAction).(hex.State)
		state = &s
		prevAction = nextAction
		turn = 1 - turn
		fmt.Printf("%v", state)

		// Call garbage collector
		runtime.GC()
	}

	// Game results
	for p := 0; p < 2; p++ {
		w := false
		if 1-turn == p {
			w = true
		}
		players[p].EndGame(prevAction, w)
	}
	if passiveClient != nil {
		passiveClient.EndGame(prevAction, false)
	}

	fmt.Printf("%v", state)
	return 1 - turn, nil
}

func playNGames(boardSize int, players [2]hexplayer.HexPlayer, passiveClient hexplayer.HexPlayer, numGames int) [2][2]int {
	startingPlayer := 0
	results := [2][2]int{}
	// results[0][0]: players[0] won, player[0] started a game
	// results[0][1]: players[0] won, player[1] started a game
	// results[1][0]: players[1] won, player[0] started a game
	// results[1][1]: players[1] won, player[1] started a game
	for g := 0; g < numGames; g++ {
		winPlayer, err := playOneGame(boardSize, players, passiveClient, startingPlayer)
		if err != nil {
			fmt.Println("Game canceled: " + err.Error())
			continue
		}

		results[winPlayer][startingPlayer]++
		fmt.Printf("Results after %d games:\n", g+1)
		fmt.Printf("\tPlayer %v: %d\n", players[0].GetColor(), players[0].GetNumberOfWins())
		fmt.Printf("\tPlayer %v: %d\n", players[1].GetColor(), players[1].GetNumberOfWins())

		// Switch roles
		startingPlayer = 1 - startingPlayer
	}
	return results
}

// Play accepts an array of two players and number of games to be played. It
// runs numGames games of Hex between the given players.
func Play(boardSize int, players [2]hexplayer.HexPlayer, numGames int, conn *websocket.Conn, resultChan chan [2][2]int) {
	if conn != nil {
		defer conn.Close()
	}

	var passiveClient hexplayer.HexPlayer
	if conn != nil && players[0].GetType() != hexplayer.HumanType && players[1].GetType() != hexplayer.HumanType {
		// Create a passive player to show the game in browser
		passiveClient = hexplayer.CreateHumanPlayer(conn, hex.None)
	}

	results := playNGames(boardSize, players, passiveClient, numGames)

	if resultChan == nil {
		fmt.Printf("*** Final results ***:\n")
		fmt.Printf("\tPlayer one: %d\n", results[0])
		fmt.Printf("\tPlayer two: %d\n", results[1])
	} else {
		resultChan <- results
	}
}
