package hexplayer

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
	"github.com/gorilla/websocket"
)

// HumanPlayer accepts client's (human's) moves. It uses a websocket to connect
// to the client.
type HumanPlayer struct {
	Color  hex.Color       // Player's color
	Webso  *websocket.Conn // Websocket connecting server and client
	numWin int             // Number of wins
}

// OpenConn opens a new websocket.
func OpenConn(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// CreateHumanPlayer creates a human player with given websocket and color of
// the player.
func CreateHumanPlayer(conn *websocket.Conn, color hex.Color) *HumanPlayer {
	hp := HumanPlayer{color, conn, 0}
	return &hp
}

// InitGame initializes the game. It sends board size and player's color to the
// client and waits for replay.
func (hp HumanPlayer) InitGame(boardSize int, firstPlayer hex.Color) error {
	m := []byte(fmt.Sprintf("INIT SIZE:%v COLOR:%v", boardSize, hp.Color))
	hp.Webso.WriteMessage(websocket.TextMessage, m)

	_, m, err := hp.Webso.ReadMessage()
	if err != nil {
		return err
	}
	if string(m) != "READY" && string(m) != "READY PASSIVE" {
		return errors.New("Invalid response: expected 'READY', got '" + string(m) + "'")
	}
	return nil
}

// PrevAction accepts last opponent's action
func (hp HumanPlayer) PrevAction(prevAction *hex.Action) {
	m := []byte(fmt.Sprintf("MOVE %v", prevAction))
	hp.Webso.WriteMessage(websocket.TextMessage, m)
}

// NextAction returns an action to be performed.
func (hp HumanPlayer) NextAction() (*hex.Action, error) {
	_, m, err := hp.Webso.ReadMessage()
	if err != nil {
		hp.Webso.WriteMessage(websocket.TextMessage, []byte("ERROR "+err.Error()))
		return nil, err
	}

	c := strings.Split(string(m), ",")
	if len(c) != 2 {
		e := "Exactly two coordinates are expected."
		hp.Webso.WriteMessage(websocket.TextMessage, []byte("ERROR "+e))
		return nil, errors.New(e)
	}

	coords := [2]byte{}
	for i := 0; i < 2; i++ {
		co, err := strconv.Atoi(c[i])
		if err != nil {
			hp.Webso.WriteMessage(websocket.TextMessage, []byte("ERROR "+err.Error()))
			return nil, err
		}
		coords[i] = byte(co)
	}

	a := hex.NewAction(coords[0], coords[1], hp.Color)

	return a, nil
}

// EndGame sends the following information to the client: last action made in
// the game, boolean value indicating whether the player has won or not.
func (hp *HumanPlayer) EndGame(lastAction *hex.Action, won bool) {
	r := 0
	if won {
		r = 1
		hp.numWin++
	}
	m := []byte(fmt.Sprintf("END %d %s", r, lastAction))
	hp.Webso.WriteMessage(websocket.TextMessage, m)

	// Wait for the human to stop analysing the game
	_, m, err := hp.Webso.ReadMessage()
	if err != nil {
		hp.Webso.WriteMessage(websocket.TextMessage, []byte("ERROR "+err.Error()))
		return
	}
}

// GetColor returns the color of the player
func (hp HumanPlayer) GetColor() hex.Color {
	return hp.Color
}

// GetNumberOfWins returns the number of wins for the player
func (hp HumanPlayer) GetNumberOfWins() int {
	return hp.numWin
}

// GetType returns the type of the player
func (hp HumanPlayer) GetType() PlayerType {
	return HumanType
}
