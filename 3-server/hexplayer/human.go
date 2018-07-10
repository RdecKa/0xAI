package hexplayer

import (
	"fmt"
	"net/http"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
	"github.com/gorilla/websocket"
)

// HumanPlayer accepts client's (human's) moves.
type HumanPlayer struct {
	color hex.Color // Player's color
	state hex.State // Current state in a game
}

func CreateHumanPlayer(w http.ResponseWriter, r *http.Request) {
	// TODO
	c, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer c.Close()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Println("Received:", msg)
	}
}

// NextAction accepts the previous action of the opponent and returns an action
// to be performed now.
func (hp HumanPlayer) NextAction(hex.Action) hex.Action {
	// TODO
	return hex.Action{}
}
