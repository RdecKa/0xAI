package cmpr

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
	"github.com/RdecKa/bachleor-thesis/server/hexgame"
	"github.com/RdecKa/bachleor-thesis/server/hexplayer"
)

// MatchSetup contains information needed for comparing two players
type MatchSetup struct {
	boardSize   int
	numGames    int
	player1type hexplayer.PlayerType
	player2type hexplayer.PlayerType
	time1       int
	time2       int
	patternFile string
}

// CreateMatch sets up the comparison of two players
//	bs: boardsize
//	ng: number of games to be played (actual number of games is twice that much,
// 		because roles are switched after ng games)
// 	p1, p2: player types
// 	t1, t2: time limits for both players
// 	patternFile: file with patterns in hex grid
func CreateMatch(bs, ng int, p1, p2 hexplayer.PlayerType, t1, t2 int, patternFile string) MatchSetup {
	return MatchSetup{
		boardSize:   bs,
		numGames:    ng,
		player1type: p1,
		player2type: p2,
		time1:       t1,
		time2:       t2,
		patternFile: patternFile,
	}
}

// Run runs a set of matches between two players
func (ms MatchSetup) Run() ([2][2][2]int, [2][2][2][2]float64) {
	resultChanWins := make(chan [2][2]int, 1)
	resultChanLength := make(chan [2][2][2]float64, 1)

	// player1 = Red, player2 = Blue
	players := [2]hexplayer.HexPlayer{
		createPlayer(ms.player1type, hex.Red, ms.time1, ms.patternFile),
		createPlayer(ms.player2type, hex.Blue, ms.time2, ms.patternFile),
	}
	hexgame.Play(ms.boardSize, players, ms.numGames, nil, resultChanWins, resultChanLength)
	results1 := <-resultChanWins
	lengths1 := <-resultChanLength

	// player1 = Blue, player2 = Red
	players = [2]hexplayer.HexPlayer{
		createPlayer(ms.player2type, hex.Red, ms.time2, ms.patternFile),
		createPlayer(ms.player1type, hex.Blue, ms.time1, ms.patternFile),
	}
	hexgame.Play(ms.boardSize, players, ms.numGames, nil, resultChanWins, resultChanLength)
	results2 := <-resultChanWins
	lengths2 := <-resultChanLength

	r := [2][2][2]int{
		results1,
		results2,
	}
	l := [2][2][2][2]float64{
		lengths1,
		lengths2,
	}
	return r, l
}

func (ms MatchSetup) String() string {
	s := fmt.Sprintf("Player 1: %v (%ds)\nPlayer 2: %v (%ds)\n",
		hexplayer.GetStringFromPlayerType(ms.player1type), ms.time1,
		hexplayer.GetStringFromPlayerType(ms.player2type), ms.time2)
	s += fmt.Sprintf("Board size: %d\nNumber of games: %d (×2)\n", ms.boardSize, ms.numGames)
	return s
}

// RunAll runs all sets of matches given as argument
func RunAll(matches []MatchSetup, resultsFileName string) {
	f, err := os.Create(resultsFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	for i, ms := range matches {
		f.WriteString("--------------------\n")
		f.WriteString(ms.String())
		f.WriteString("--------------------\n")
		results, lengths := ms.Run()
		f.WriteString(fmt.Sprintf("FINAL RESULTS for set #%d:\n", i))

		r1 := results[0]
		l1 := lengths[0]
		f.WriteString("\n---> Roles: Player 1 = Red, Player 2 = Blue\n")
		f.WriteString("Number of wins:\n")
		f.WriteString("\tFirst move:  P1  P2\n")
		f.WriteString(fmt.Sprintf("\tPlayer 1:   %3d %3d\n", r1[0][0], r1[0][1]))
		f.WriteString(fmt.Sprintf("\tPlayer 2:   %3d %3d\n\n", r1[1][0], r1[1][1]))

		f.WriteString("Game length (avg, std):\n")
		f.WriteString("\tFirst move: P1               P2\n")
		f.WriteString(fmt.Sprintf("\tPlayer 1: (%6.2f, %6.2f) (%6.2f, %6.2f)\n",
			l1[0][0][0], l1[0][0][1], l1[0][1][0], l1[0][1][1]))
		f.WriteString(fmt.Sprintf("\tPlayer 2: (%6.2f, %6.2f) (%6.2f, %6.2f)\n",
			l1[1][0][0], l1[1][0][1], l1[1][1][0], l1[1][1][1]))

		r2 := results[1]
		l2 := lengths[1]
		f.WriteString("\n---> Roles: Player 1 = Blue, Player 2 = Red\n")
		f.WriteString("First move:  P1  P2\n")
		f.WriteString(fmt.Sprintf("\tPlayer 1:   %3d %3d\n", r2[1][1], r2[1][0]))
		f.WriteString(fmt.Sprintf("\tPlayer 2:   %3d %3d\n\n", r2[0][1], r2[0][0]))

		f.WriteString("Game length (avg, std):\n")
		f.WriteString("\tFirst move: P1               P2\n")
		f.WriteString(fmt.Sprintf("\tPlayer 1: (%6.2f, %6.2f) (%6.2f, %6.2f)\n",
			l2[1][1][0], l2[1][1][1], l2[1][0][0], l2[1][0][1]))
		f.WriteString(fmt.Sprintf("\tPlayer 2: (%6.2f, %6.2f) (%6.2f, %6.2f)\n",
			l2[0][1][0], l2[0][1][1], l2[0][0][0], l2[0][0][1]))
	}
	f.WriteString("\nTesting finished.\n")
}

func createPlayer(t hexplayer.PlayerType, c hex.Color, tl int, patternFile string) hexplayer.HexPlayer {
	switch t {
	case hexplayer.MctsType:
		return hexplayer.CreateMCTSplayer(c, math.Sqrt(2), time.Duration(tl)*time.Second, 10, true)
	case hexplayer.AbDtType:
		return hexplayer.CreateAbPlayer(c, nil, time.Duration(tl)*time.Second,
			true, patternFile, false, hexplayer.AbDtType)
	case hexplayer.AbLrType:
		return hexplayer.CreateAbPlayer(c, nil, time.Duration(tl)*time.Second,
			true, patternFile, false, hexplayer.AbLrType)
	default:
		fmt.Println(fmt.Errorf("Invalid type '%s'", hexplayer.GetStringFromPlayerType(t)))
		return nil
	}
}
