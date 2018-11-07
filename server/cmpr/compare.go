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
	boardSize  int
	numGames   int
	redPlayer  hexplayer.HexPlayer
	bluePlayer hexplayer.HexPlayer
	timeRed    int
	timeBlue   int
}

// CreateMatch sets up the comparison of two players
func CreateMatch(bs, ng int, p1, p2 hexplayer.PlayerType, t1, t2 int, patternFile string) MatchSetup {
	rp := createPlayer(p1, hex.Red, t1, patternFile)
	bp := createPlayer(p2, hex.Blue, t2, patternFile)
	return MatchSetup{
		boardSize:  bs,
		numGames:   ng,
		redPlayer:  rp,
		bluePlayer: bp,
		timeRed:    t1,
		timeBlue:   t2,
	}
}

// Run runs a set of matches between two players
func (ms MatchSetup) Run() [2]int {
	players := [2]hexplayer.HexPlayer{ms.redPlayer, ms.bluePlayer}
	resultChan := make(chan [2]int, 1)
	hexgame.Play(ms.boardSize, players, ms.numGames, nil, resultChan)
	results := <-resultChan
	return results
}

func (ms MatchSetup) String() string {
	s := fmt.Sprintf("Red player: %v (%ds)\nBlue player: %v (%ds)\n",
		hexplayer.GetStringFromPlayerType(ms.redPlayer.GetType()), ms.timeRed,
		hexplayer.GetStringFromPlayerType(ms.bluePlayer.GetType()), ms.timeBlue)
	s += fmt.Sprintf("Board size: %d\nNumber of games: %d\n", ms.boardSize, ms.numGames)
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
		results := ms.Run()
		f.WriteString(fmt.Sprintf("Final results for set #%d:\n", i))
		f.WriteString(fmt.Sprintf("Red: %d\n", results[0]))
		f.WriteString(fmt.Sprintf("Blue: %d\n\n", results[1]))
	}
	f.WriteString("Testing finished.\n")
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
