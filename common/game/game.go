package game

// -----------------
// |     State     |
// -----------------

// State represents a state in a game
type State interface {
	String() string
	GetPossibleActions() []Action
	GetSuccessorState(Action) State
	IsGoalState(bool) (bool, interface{})
	EvaluateGoalState() float64
	Same(State) bool
	GenSample(float64, chan []uint32, chan [2][]int) string // Returns a string representing state attributes for supervised machine learning
}

// ------------------
// |     Action     |
// ------------------

// Action represents an action in a game
type Action interface {
	String() string
}
