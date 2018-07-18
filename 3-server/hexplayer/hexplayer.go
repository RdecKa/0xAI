package hexplayer

import "github.com/RdecKa/bachleor-thesis/common/game/hex"

// HexPlayer represents a player of hex that can be either human or computer.
type HexPlayer interface {
	InitGame(int) error                          // Initializes game with a grid of a given size
	NextAction(*hex.Action) (*hex.Action, error) // Acepts opponent's last action, returns an action to be performed
	EndGame(*hex.Action, bool)                   // Accepts last action in the game and boolean value indicating whether the player has won or not
	SwitchColor()                                // Switches the color of the player (RED -> BLUE or BLUE -> RED)
	GetNumberOfWins() int                        // Returns the number of wins
}
