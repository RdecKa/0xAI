package hexplayer

import (
	"fmt"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

// ------------------
// |     PlayerType     |
// ------------------

type PlayerType byte

// enum for player types
const (
	Unknown   PlayerType = 0
	HumanType PlayerType = 1
	MctsType  PlayerType = 2
	AbDtType  PlayerType = 3
	AbLrType  PlayerType = 4
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

func GetPlayerTypeFromString(t string) PlayerType {
	switch t {
	case "human":
		return HumanType
	case "mcts":
		return MctsType
	case "abDT":
		return AbDtType
	case "abLR":
		return AbLrType
	default:
		fmt.Println(fmt.Errorf("Invalid type '%s'", t))
		return Unknown
	}
}

func (t PlayerType) String() string {
	switch t {
	case HumanType:
		return "human"
	case MctsType:
		return "mcts"
	case AbDtType:
		return "abDT"
	case AbLrType:
		return "abLR"
	default:
		fmt.Println(fmt.Errorf("Invalid type '%s'", t))
		return ""
	}
}
