package hexplayer

import "github.com/RdecKa/bachleor-thesis/common/game/hex"

// HexPlayer represents a player of hex that can be either human or computer.
type HexPlayer interface {
	NextAction(hex.Action) hex.Action
}
