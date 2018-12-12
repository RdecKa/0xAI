package hexplayer

import (
	"math/rand"

	"github.com/RdecKa/0xAI/common/game/hex"
)

// RandPlayer represents a computer player that randomly selects actions
type RandPlayer struct {
	Color  hex.Color
	state  *hex.State
	numWin int
}

// CreateRandPlayer creates a new player
func CreateRandPlayer(c hex.Color) *RandPlayer {
	rp := RandPlayer{c, nil, 0}
	return &rp
}

// InitGame initializes the game
func (rp *RandPlayer) InitGame(boardSize int, firstPlayer hex.Color) error {
	initState := hex.NewState(byte(boardSize), firstPlayer)
	rp.state = initState
	return nil
}

// PrevAction accepts opponent's last action
func (rp *RandPlayer) PrevAction(prevAction *hex.Action) {
	if prevAction != nil {
		rp.updatePlayerState(prevAction)
	}
}

// NextAction returns a randomly chosen action to be performed
func (rp *RandPlayer) NextAction() (*hex.Action, error) {
	actions := rp.state.GetPossibleActions()
	selected := actions[rand.Intn(len(actions))].(*hex.Action)
	rp.updatePlayerState(selected)
	return selected, nil
}

// updatePlayerState updates the game state of the player
func (rp *RandPlayer) updatePlayerState(a *hex.Action) {
	s := rp.state.GetSuccessorState(a).(hex.State)
	rp.state = &s
}

// EndGame accepts the result of the game
func (rp *RandPlayer) EndGame(lastAction *hex.Action, won bool) {
	if won {
		rp.numWin++
	}
}

// GetColor returns the color of the player
func (rp RandPlayer) GetColor() hex.Color {
	return rp.Color
}

// GetNumberOfWins returns the number of wins for this player
func (rp RandPlayer) GetNumberOfWins() int {
	return rp.numWin
}

// GetType returns the type of the player
func (rp RandPlayer) GetType() PlayerType {
	return RandType
}
