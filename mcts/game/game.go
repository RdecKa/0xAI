package game

// -----------------
// |     State     |
// -----------------

// State represents a state in a game
type State interface {
	String() string
	GetPossibleActions() []Action
	GetSuccessorState(Action) State
	IsGoalState() bool
	EvaluateGoalState() float64
}

// ------------------
// |     Action     |
// ------------------

// Action represents an action in a game
type Action interface {
	String() string
}
