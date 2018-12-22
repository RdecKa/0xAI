package hexplayer

import (
	"math"
	"time"

	"github.com/RdecKa/0xAI/common/game/hex"
)

// HybridPlayer represents a computer player that combines strategies of MCTS and AB
type HybridPlayer struct {
	Color              hex.Color
	state              *hex.State
	numWin             int
	lastOpponentAction *hex.Action
	subPlayers         [2]HexPlayer
	activeSubplayer    int // Stores index of the subplayer that is active (0 or 1)
	changeTypeAt       int // When this number of stones is reached, AB changes to MCTS
	numStonesplaced    int
}

// CreateHybridPlayer creates a new player
func CreateHybridPlayer(c hex.Color, t time.Duration, allowResignation bool,
	patFileName string, ABsubtype PlayerType, changeTypeAt int) *HybridPlayer {
	ABsubPlayer := CreateAbPlayer(c, nil, t, allowResignation, patFileName, false, ABsubtype)
	MCTSsubPlayer := CreateMCTSplayer(c, math.Sqrt(2), t, 10, allowResignation)
	hp := HybridPlayer{c, nil, 0, nil, [2]HexPlayer{ABsubPlayer, MCTSsubPlayer}, 0, changeTypeAt, 0}
	return &hp
}

// InitGame initializes the game
func (hp *HybridPlayer) InitGame(boardSize int, firstPlayer hex.Color) error {
	initState := hex.NewState(byte(boardSize), firstPlayer)
	hp.state = initState
	hp.lastOpponentAction = nil
	hp.activeSubplayer = 0
	hp.numStonesplaced = 0
	hp.subPlayers[0].InitGame(boardSize, firstPlayer)
	return nil
}

// PrevAction accepts opponent's last action
func (hp *HybridPlayer) PrevAction(prevAction *hex.Action) {
	if prevAction != nil {
		hp.subPlayers[hp.activeSubplayer].PrevAction(prevAction)
		hp.lastOpponentAction = prevAction
		hp.updatePlayerState(prevAction)
	}
}

// NextAction returns an action to be performed
func (hp *HybridPlayer) NextAction() (*hex.Action, error) {
	selected, err := hp.subPlayers[hp.activeSubplayer].NextAction()
	if err != nil {
		return nil, err
	}
	hp.updatePlayerState(selected)
	return selected, err
}

// updatePlayerState updates the game state of the player
func (hp *HybridPlayer) updatePlayerState(a *hex.Action) {
	if a == nil {
		return
	}

	s := hp.state.GetSuccessorState(a).(hex.State)
	hp.state = &s
	hp.numStonesplaced++
	if hp.numStonesplaced == hp.changeTypeAt {
		hp.activeSubplayer = 1
		s := hp.state.Clone().(hex.State)
		(hp.subPlayers[1]).(*MCTSplayer).initGameFromState(&s, hp.lastOpponentAction)
	}
}

// EndGame accepts the result of the game
func (hp *HybridPlayer) EndGame(lastAction *hex.Action, won bool) {
	if won {
		hp.numWin++
		hp.subPlayers[0].EndGame(lastAction, won)
		hp.subPlayers[1].EndGame(lastAction, won)
	}
}

// GetColor returns the color of the player
func (hp HybridPlayer) GetColor() hex.Color {
	return hp.Color
}

// GetNumberOfWins returns the number of wins for this player
func (hp HybridPlayer) GetNumberOfWins() int {
	return hp.numWin
}

// GetType returns the type of the player
func (hp HybridPlayer) GetType() PlayerType {
	return HybridType
}
