package hexplayer

import (
	"github.com/RdecKa/bachleor-thesis/3-ab/ab"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

// AbPlayer represents a computer player that uses alpha-beta pruning for
// selecting moves
type AbPlayer struct {
	Color              hex.Color   // Player's color
	numWin             int         // Number of wins
	state              *hex.State  // Current state in a game
	lastOpponentAction *hex.Action // Opponent's last action
}

// CreateAbPlayer creates a new player
func CreateAbPlayer(c hex.Color) *AbPlayer {
	ap := AbPlayer{c, 0, nil, nil}
	return &ap
}

// InitGame initializes the game
func (ap *AbPlayer) InitGame(boardSize int, firstPlayer hex.Color) error {
	ap.state = hex.NewState(byte(boardSize), firstPlayer)
	return nil
}

// PrevAction accepts opponent's last action
func (ap *AbPlayer) PrevAction(prevAction *hex.Action) {
	if prevAction != nil {
		s := ap.state.GetSuccessorState(prevAction).(hex.State)
		ap.state = &s
		ap.lastOpponentAction = prevAction
	}
}

// NextAction returns an action to be performed
func (ap *AbPlayer) NextAction() (*hex.Action, error) {
	chosecAction := ab.AlphaBeta(ap.state)
	s := ap.state.GetSuccessorState(chosecAction).(hex.State)
	ap.state = &s
	return chosecAction, nil
}

// EndGame accepts the result of the game
func (ap *AbPlayer) EndGame(lastAction *hex.Action, won bool) {
	if won {
		ap.numWin++
	}
}

// GetColor returns the color of the player
func (ap AbPlayer) GetColor() hex.Color {
	return ap.Color
}

// GetNumberOfWins returns the number of wins for this player
func (ap AbPlayer) GetNumberOfWins() int {
	return ap.numWin
}

// GetType returns the type of the player
func (ap AbPlayer) GetType() PlayerType {
	return AbType
}
