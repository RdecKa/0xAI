package hexplayer

import "github.com/RdecKa/bachleor-thesis/common/game/hex"

// HexPlayer represents a player of hex that can be either human or computer.
type HexPlayer interface {
	InitGame(int, hex.Color) error               // Initializes game with a grid of a given size and first player
	NextAction(*hex.Action) (*hex.Action, error) // Acepts opponent's last action, returns an action to be performed
	EndGame(*hex.Action, bool)                   // Accepts last action in the game and boolean value indicating whether the player has won or not
	SwitchColor()                                // Switches the color of the player (RED -> BLUE or BLUE -> RED)
	GetColor() hex.Color                         // Gets the color of the player
	GetNumberOfWins() int                        // Returns the number of wins
}
