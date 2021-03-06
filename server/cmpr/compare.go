package cmpr

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/RdecKa/0xAI/common/game/hex"
	"github.com/RdecKa/0xAI/server/hexgame"
	"github.com/RdecKa/0xAI/server/hexplayer"
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
	extraInfo1  interface{}
	extraInfo2  interface{}
}

// CreateMatch sets up the comparison of two players
//	bs: boardsize
//	ng: number of games to be played (actual number of games is twice that much,
// 		because roles are switched after ng games)
// 	p1, p2: player types
// 	t1, t2: time limits for both players
// 	patternFile: file with patterns in hex grid
// ei1, ei1: additional parameters for players
func CreateMatch(bs, ng int, p1, p2 hexplayer.PlayerType, t1, t2 int, patternFile string, ei1, ei2 interface{}) MatchSetup {
	return MatchSetup{
		boardSize:   bs,
		numGames:    ng,
		player1type: p1,
		player2type: p2,
		time1:       t1,
		time2:       t2,
		patternFile: patternFile,
		extraInfo1:  ei1,
		extraInfo2:  ei2,
	}
}

type result struct {
	results [2][2]int
	lengths [2][2][2]float64
}

func runParallel(ms MatchSetup, outDir string, players [2]hexplayer.HexPlayer, ch chan result) {
	results, lengths := ms.Run(outDir, players)
	ch <- result{results, lengths}
}

// Run runs a set of matches between two players
func (ms MatchSetup) Run(outDir string, players [2]hexplayer.HexPlayer) ([2][2]int, [2][2][2]float64) {
	resultChanWins := make(chan [2][2]int, 1)
	resultChanLength := make(chan [2][2][2]float64, 1)
	outDir += "games/"
	os.Mkdir(outDir, os.ModePerm)

	hexgame.Play(ms.boardSize, players, ms.numGames, nil, resultChanWins, resultChanLength, outDir)
	results := <-resultChanWins
	lengths := <-resultChanLength

	return results, lengths
}

func (ms MatchSetup) String() string {
	s := fmt.Sprintf("Player 1: %v (%ds)\nPlayer 2: %v (%ds)\n",
		ms.player1type.String(), ms.time1,
		ms.player2type.String(), ms.time2)
	s += fmt.Sprintf("Board size: %d\nNumber of games: %d (x2)\n", ms.boardSize, ms.numGames)
	return s
}

// RunAll runs all sets of matches given as argument
func RunAll(matches []MatchSetup, outDir string) {
	f, err := os.Create(outDir + "test_results.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	ch0, ch1 := make(chan result, 1), make(chan result, 1)

	f.WriteString(fmt.Sprintf("Testing started at %s.\n\n", time.Now().Format("15.04.05 (2006/01/02)")))
	for _, ms := range matches {
		f.WriteString("--------------------\n")
		f.WriteString(ms.String())
		f.WriteString("--------------------\n")

		var players [2][2]hexplayer.HexPlayer
		// player1 = Red, player2 = Blue
		players[0] = [2]hexplayer.HexPlayer{
			createPlayer(ms.player1type, hex.Red, ms.time1, ms.patternFile, ms.extraInfo1),
			createPlayer(ms.player2type, hex.Blue, ms.time2, ms.patternFile, ms.extraInfo2),
		}
		// player1 = Blue, player2 = Red
		players[1] = [2]hexplayer.HexPlayer{
			createPlayer(ms.player2type, hex.Red, ms.time2, ms.patternFile, ms.extraInfo2),
			createPlayer(ms.player1type, hex.Blue, ms.time1, ms.patternFile, ms.extraInfo1),
		}

		go runParallel(ms, outDir, players[0], ch0)
		go runParallel(ms, outDir, players[1], ch1)

		for p := 0; p <= 1; p++ {
			pls := players[p]

			var p1p1, p1p2, p2p1, p2p2 int
			var f1f1, f1f2, f2f1, f2f2 [2]float64
			if p == 0 {
				r := <-ch0
				results := r.results
				lengths := r.lengths
				p1p1, f1f1 = results[0][0], lengths[0][0]
				p1p2, f1f2 = results[0][1], lengths[0][1]
				p2p1, f2f1 = results[1][0], lengths[1][0]
				p2p2, f2f2 = results[1][1], lengths[1][1]
			} else {
				r := <-ch1
				results := r.results
				lengths := r.lengths
				p1p1, f1f1 = results[1][1], lengths[1][1]
				p1p2, f1f2 = results[1][0], lengths[1][0]
				p2p1, f2f1 = results[0][1], lengths[0][1]
				p2p2, f2f2 = results[0][0], lengths[0][0]
			}

			f.WriteString(fmt.Sprintf("\n---> Roles: Player 1 = %s, Player 2 = %s\n",
				pls[p].GetColor().String(), pls[1-p].GetColor().String()))
			f.WriteString("Number of wins:\n")
			f.WriteString("\tFirst move:  P1  P2\n")
			f.WriteString(fmt.Sprintf("\tPlayer 1:   %3d %3d\n", p1p1, p1p2))
			f.WriteString(fmt.Sprintf("\tPlayer 2:   %3d %3d\n\n", p2p1, p2p2))

			f.WriteString("Game length (avg, std):\n")
			f.WriteString("\tFirst move: P1               P2\n")
			f.WriteString(fmt.Sprintf("\tPlayer 1: (%6.2f, %6.2f) (%6.2f, %6.2f)\n",
				f1f1[0], f1f1[1], f1f2[0], f1f2[1]))
			f.WriteString(fmt.Sprintf("\tPlayer 2: (%6.2f, %6.2f) (%6.2f, %6.2f)\n",
				f2f1[0], f2f1[1], f2f2[0], f2f2[1]))
		}
		f.WriteString(fmt.Sprintf("\nGroup finished at %s.\n\n", time.Now().Format("15.04.05 (2006/01/02)")))
	}
	f.WriteString(fmt.Sprintf("\nTesting finished at %s.\n", time.Now().Format("15.04.05 (2006/01/02)")))
}

func createPlayer(t hexplayer.PlayerType, c hex.Color, tl int, patternFile string, ei interface{}) hexplayer.HexPlayer {
	switch t {
	case hexplayer.RandType:
		return hexplayer.CreateRandPlayer(c)
	case hexplayer.MctsType:
		return hexplayer.CreateMCTSplayer(c, math.Sqrt(2), time.Duration(tl)*time.Second, 10, true)
	case hexplayer.AbDtType:
		return hexplayer.CreateAbPlayer(c, nil, time.Duration(tl)*time.Second,
			true, patternFile, false, hexplayer.AbDtType)
	case hexplayer.AbLrType:
		return hexplayer.CreateAbPlayer(c, nil, time.Duration(tl)*time.Second,
			true, patternFile, false, hexplayer.AbLrType)
	case hexplayer.HybridType:
		return hexplayer.CreateHybridPlayer(c, time.Duration(tl)*time.Second,
			true, patternFile, hexplayer.AbLrType, ei.(int))
	default:
		fmt.Println(fmt.Errorf("Invalid type '%s'", t.String()))
		return nil
	}
}
