package hexplayer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/RdecKa/0xAI/3-ab/ab"
	"github.com/RdecKa/0xAI/common/game/hex"
	"github.com/gorilla/websocket"
)

// AbPlayer represents a computer player that uses alpha-beta pruning for
// selecting moves
type AbPlayer struct {
	Color              hex.Color                  // Player's color
	subtype            PlayerType                 // Player's subtype (DT/LR)
	Webso              *websocket.Conn            // Websocket connecting server and client
	timeToRun          time.Duration              // Time given to select an action
	numWin             int                        // Number of wins
	state              *hex.State                 // Current state in a game
	safeWinCells       [][2]cell                  // List of cells under the bridges on a winning path
	lastOpponentAction *hex.Action                // Opponent's last action
	allowResignation   bool                       // Allow the player to resign if the game is lost
	createTree         bool                       // If true, create a search tree for debugging purposes
	gridChan           chan []uint32              // Used for pattern checking
	stopChan           chan struct{}              // -||-
	patChan            chan []int                 // -||-
	resultChan         chan [2][]int              // -||-
	getEstimatedValue  func(s *ab.Sample) float64 // Function used for evaluating states
}

// CreateAbPlayer creates a new player
func CreateAbPlayer(c hex.Color, webso *websocket.Conn, t time.Duration,
	allowResignation bool, patFileName string, createTree bool, subtype PlayerType) *AbPlayer {

	gridChan, patChan, stopChan, resultChan := hex.CreatePatChecker(patFileName)
	ap := AbPlayer{
		Color:             c,
		subtype:           subtype,
		Webso:             webso,
		timeToRun:         t,
		allowResignation:  allowResignation,
		createTree:        createTree,
		gridChan:          gridChan,
		patChan:           patChan,
		stopChan:          stopChan,
		resultChan:        resultChan,
		getEstimatedValue: ab.GetEstimateFunction(subtype.String())}
	return &ap
}

// InitGame initializes the game
func (ap *AbPlayer) InitGame(boardSize int, firstPlayer hex.Color) error {
	ap.state = hex.NewState(byte(boardSize), firstPlayer)
	ap.safeWinCells = nil
	ap.lastOpponentAction = nil
	return nil
}

// PrevAction accepts opponent's last action
func (ap *AbPlayer) PrevAction(prevAction *hex.Action) {
	if prevAction != nil {
		ap.updatePlayerState(prevAction)
		ap.lastOpponentAction = prevAction
	}
}

// NextAction returns an action to be performed. It returns nil when the player
// decides to resign.
func (ap *AbPlayer) NextAction() (*hex.Action, error) {
	// Check if the player has already won (has a virtual connection)
	if a, swc, ok := getActionIfWinningPathExists(ap.lastOpponentAction, ap.safeWinCells, ap.Color); ok {
		ap.updatePlayerState(a)
		ap.safeWinCells = swc
		return a, nil
	}

	// Run Minimax with alpha-beta pruning
	chosenAction, searchedTree := ab.AlphaBeta(ap.state, ap.timeToRun, ap.createTree,
		ap.gridChan, ap.patChan, ap.resultChan, ap.getEstimatedValue, ap.subtype.String())

	if chosenAction == nil {
		if !ap.allowResignation {
			a, err := doNotLoseHope(ap.state, ap.Color)
			if err != nil {
				return nil, err
			}

			// Update state
			ap.updatePlayerState(a)

			return a, nil
		}
		return nil, nil
	}

	// Send JSON to the client for debugging purposes
	if ap.createTree {
		jsonText, err := json.Marshal(searchedTree)
		if err != nil {
			fmt.Print(fmt.Errorf("Error creating JSON of searchedTree: %s", err))
		} else {
			message := fmt.Sprintf("ABJSON %s", jsonText)
			ap.Webso.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}

	// Update state
	ap.updatePlayerState(chosenAction)

	// Check if player has a virtual connection
	if exists, solution := ap.state.IsGoalState(false); exists {
		fmt.Println(ap.subtype.String() + " player has a virtual connection!")
		winPath := solution.([][2]int)
		ap.safeWinCells = findSafeCells(winPath, ap.state.GetSize(), ap.Color)
	}

	return chosenAction, nil
}

// updatePlayerState updates the game state of the player
func (ap *AbPlayer) updatePlayerState(a *hex.Action) {
	s := ap.state.GetSuccessorState(a).(hex.State)
	ap.state = &s
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
	return ap.subtype
}
