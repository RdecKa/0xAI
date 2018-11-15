package hexgame

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
	"github.com/RdecKa/bachleor-thesis/server/hexplayer"
	"github.com/gorilla/websocket"
)

// playOneGame returns 0 if the first player (of players) won and 1 if the
// second player won
func playOneGame(boardSize int, players [2]hexplayer.HexPlayer, passiveClient hexplayer.HexPlayer,
	startingPlayer int, outFile *os.File) (int, int, error) {

	fmt.Println("---------------------------------------------")

	// Init game
	for p := 0; p < 3; p++ {
		var err error
		if p < 2 {
			err = players[p].InitGame(boardSize, players[startingPlayer].GetColor())
		} else if p == 2 && passiveClient != nil {
			err = passiveClient.InitGame(boardSize, players[startingPlayer].GetColor())
		}

		if err != nil {
			return -1, -1, err
		}
	}
	state := hex.NewState(byte(boardSize), players[startingPlayer].GetColor())
	turn := startingPlayer
	var prevAction *hex.Action
	gameLength := 0

	// Play game
	for g, _ := state.IsGoalState(true); !g; g, _ = state.IsGoalState(true) {
		players[turn].PrevAction(prevAction)
		nextAction, err := players[turn].NextAction()
		if err != nil {
			fmt.Println(err)
			return -1, -1, err
		}
		if nextAction == nil {
			// Player has resigned
			outFile.WriteString(fmt.Sprintf("Player %v resigned!\n", players[turn].GetColor()))
			break
		}
		if passiveClient != nil {
			passiveClient.PrevAction(nextAction)
		}
		s := state.GetSuccessorState(nextAction).(hex.State)
		state = &s
		prevAction = nextAction
		turn = 1 - turn
		gameLength++
		outFile.WriteString(fmt.Sprintf("%v", state))

		// Call garbage collector
		runtime.GC()
	}

	outFile.WriteString(fmt.Sprintf("Game length: %d\n", gameLength))

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
	return 1 - turn, gameLength, nil
}

func playNGames(boardSize int, players [2]hexplayer.HexPlayer, passiveClient hexplayer.HexPlayer,
	numGames int, outFile *os.File) ([2][2]int, [2][2][]int) {
	startingPlayer := 0
	gameLengthList := [2][2][]int{}
	gameLengthList[0][0] = make([]int, 0)
	gameLengthList[0][1] = make([]int, 0)
	gameLengthList[1][0] = make([]int, 0)
	gameLengthList[1][1] = make([]int, 0)
	results := [2][2]int{}
	// results[0][0]: players[0] won, player[0] started a game
	// results[0][1]: players[0] won, player[1] started a game
	// results[1][0]: players[1] won, player[0] started a game
	// results[1][1]: players[1] won, player[1] started a game
	for g := 0; g < numGames; g++ {
		winPlayer, gameLength, err := playOneGame(boardSize, players, passiveClient, startingPlayer, outFile)
		if err != nil {
			fmt.Println("Game canceled: " + err.Error())
			continue
		}

		results[winPlayer][startingPlayer]++
		gameLengthList[winPlayer][startingPlayer] = append(gameLengthList[winPlayer][startingPlayer], gameLength)
		fmt.Printf("Results after %d games:\n", g+1)
		fmt.Printf("\tPlayer %v: %d\n", players[0].GetColor(), players[0].GetNumberOfWins())
		fmt.Printf("\tPlayer %v: %d\n", players[1].GetColor(), players[1].GetNumberOfWins())

		// Switch roles
		startingPlayer = 1 - startingPlayer
	}
	return results, gameLengthList
}

// Play accepts an array of two players and number of games to be played. It
// runs numGames games of Hex between the given players.
func Play(boardSize int, players [2]hexplayer.HexPlayer, numGames int,
	conn *websocket.Conn, resultChanWins chan [2][2]int,
	resultChanLengths chan [2][2][2]float64, outDir string) {

	if conn != nil {
		defer conn.Close()
	}

	var passiveClient hexplayer.HexPlayer
	if conn != nil && players[0].GetType() != hexplayer.HumanType && players[1].GetType() != hexplayer.HumanType {
		// Create a passive player to show the game in browser
		passiveClient = hexplayer.CreateHumanPlayer(conn, hex.None)
	}

	outFile, err := os.Create(fmt.Sprintf("%sgames_%s_%s_%s_%d.txt",
		outDir, time.Now().Format("150405"),
		players[0].GetType().String(),
		players[1].GetType().String(), numGames))
	if err != nil {
		fmt.Println(err)
	}

	outFile.WriteString(fmt.Sprintf("%s: %s\n", players[0].GetColor().String(),
		players[0].GetType().String()))
	outFile.WriteString(fmt.Sprintf("%s: %s\n", players[1].GetColor().String(),
		players[1].GetType().String()))

	results, gameLengthList := playNGames(boardSize, players, passiveClient, numGames, outFile)
	lengths := [2][2][2]float64{}
	for wp := range gameLengthList {
		for sp := range gameLengthList[wp] {
			avgLen := avg(gameLengthList[wp][sp])
			stdLen := std(gameLengthList[wp][sp], avgLen)
			lengths[wp][sp] = [2]float64{avgLen, stdLen}
		}
	}

	outFile.WriteString("\n*** Final results ***:\n")
	outFile.WriteString(fmt.Sprintf("Player %s: %d\n", players[0].GetColor().String(), results[0]))
	outFile.WriteString(fmt.Sprintf("Player %s: %d\n", players[1].GetColor().String(), results[1]))

	if resultChanWins != nil {
		resultChanWins <- results
	}
	if resultChanLengths != nil {
		resultChanLengths <- lengths
	}
}

func avg(lst []int) float64 {
	if len(lst) == 0 {
		return 0
	}
	s := 0
	for _, el := range lst {
		s += el
	}
	return float64(s) / float64(len(lst))
}

func std(lst []int, avg float64) float64 {
	if len(lst) == 0 {
		return 0.0
	}

	std := 0.0
	for _, el := range lst {
		std += math.Pow(float64(el)-avg, 2)
	}
	std /= float64(len(lst)) // This is variance
	return math.Sqrt(std)    // Square root to get the standard deviation
}
