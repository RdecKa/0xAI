package hexplayer

import "github.com/RdecKa/bachleor-thesis/common/game/hex"

// ------------------
// |     PlayerType     |
// ------------------

type PlayerType byte

// enum for player types
const (
	HumanType PlayerType = 0
	MctsType  PlayerType = 1
	AbType    PlayerType = 2
)

// HexPlayer represents a player of hex that can be either human or computer.
type HexPlayer interface {
	InitGame(int, hex.Color) error    // Initializes game with a grid of a given size and first player
	PrevAction(*hex.Action)           // Acepts opponent's last action
	NextAction() (*hex.Action, error) // Returns an action to be performed
	EndGame(*hex.Action, bool)        // Accepts last action in the game and boolean value indicating whether the player has won or not
	GetColor() hex.Color              // Gets the color of the player
	GetNumberOfWins() int             // Returns the number of wins
	GetType() PlayerType              // Returns the type of the player
}
