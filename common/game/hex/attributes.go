// Package hex - in this file all possible attributes that can be used for
// evaluating a hex board are listed.
//
// To add an attribute, do the following:
// 	- Implement a type implementing game.Attribute
// 	- Initialize instance(s) of that attribute
// 	- Add this/these instance(s) to the slice GenSamAttributes (together with
// 		matching opposite attribute)
// 	- In 2-ml/regression.py, select a type of the attribute when reading a CSV
// 		file (for now only integer values are supported)
// 	- In 3-ab/ab/ab.go, add a line to initialization of Sample sample for each
// 		instance of the attribute
//
// To remove an attribute, simply delete it from the GenSamAttributes. To
// completely remove it, undo the steps listed in instructions for adding an
// attribute.
package hex

import (
	"fmt"

	"github.com/RdecKa/bachleor-thesis/common/game"
)

// Initialization of all available attributes
var (
	AttrNumStones          = AttrNumberStones{}
	AttrOccRedRows         = AttrOccupiedRowsCols{Red, true}
	AttrOccRedCols         = AttrOccupiedRowsCols{Red, false}
	AttrOccBlueRows        = AttrOccupiedRowsCols{Blue, true}
	AttrOccBlueCols        = AttrOccupiedRowsCols{Blue, false}
	AttrPatCountRed0       = AttrPatternCount{Red, 0}
	AttrPatCountRed1       = AttrPatternCount{Red, 1}
	AttrPatCountRed2       = AttrPatternCount{Red, 2}
	AttrPatCountRed3       = AttrPatternCount{Red, 3}
	AttrPatCountRed4       = AttrPatternCount{Red, 4}
	AttrPatCountBlue0      = AttrPatternCount{Blue, 0}
	AttrPatCountBlue1      = AttrPatternCount{Blue, 1}
	AttrPatCountBlue2      = AttrPatternCount{Blue, 2}
	AttrPatCountBlue3      = AttrPatternCount{Blue, 3}
	AttrPatCountBlue4      = AttrPatternCount{Blue, 4}
	AttrLastPlayer         = AttrLastPlayerTurn{true}
	AttrLastPlayerOpponent = AttrLastPlayerTurn{false}
	AttrDistanceToCenter   = AttrLastActionDistanceToCenter{}
)

// GenSamAttributes contains the attributes that are included in the sample
// generation. Each sub-slice represents a pair of attributes that are oppposite
// to each other. This information is used in generation of learning samples
// when two samples are generated for each state - one as it is and one with
// switched roles of red and blue player.
// If the second element of a pair is nil, the attribute is the same for both
// players.
var GenSamAttributes = [][2]game.Attribute{
	[2]game.Attribute{AttrNumStones, nil},
	[2]game.Attribute{AttrOccRedRows, AttrOccBlueCols},
	[2]game.Attribute{AttrOccRedCols, AttrOccBlueRows},
	[2]game.Attribute{AttrOccBlueRows, AttrOccRedCols},
	[2]game.Attribute{AttrOccBlueCols, AttrOccRedRows},
	[2]game.Attribute{AttrPatCountRed0, AttrPatCountBlue0},
	[2]game.Attribute{AttrPatCountRed1, AttrPatCountBlue1},
	[2]game.Attribute{AttrPatCountRed2, AttrPatCountBlue2},
	[2]game.Attribute{AttrPatCountRed3, AttrPatCountBlue3},
	[2]game.Attribute{AttrPatCountRed4, AttrPatCountBlue4},
	[2]game.Attribute{AttrPatCountBlue0, AttrPatCountRed0},
	[2]game.Attribute{AttrPatCountBlue1, AttrPatCountRed1},
	[2]game.Attribute{AttrPatCountBlue2, AttrPatCountRed2},
	[2]game.Attribute{AttrPatCountBlue3, AttrPatCountRed3},
	[2]game.Attribute{AttrPatCountBlue4, AttrPatCountRed4},
	[2]game.Attribute{AttrLastPlayer, AttrLastPlayerOpponent},
	[2]game.Attribute{AttrDistanceToCenter, nil},
}

// ----------------------------
// |     AttrNumberStones     |
// ----------------------------

// AttrNumberStones takes care of getting the number of stones on the board
type AttrNumberStones struct{}

// GetAttributeName returns the name of an attribute
func (a AttrNumberStones) GetAttributeName() string {
	return "num_stones"
}

// GetAttributeValue returns the value of an attribute
func (a AttrNumberStones) GetAttributeValue(args *[]interface{}) int {
	patCount := (*args)[1].([2][]int)
	return patCount[0][0] + patCount[1][0] // red_p0 + blue_p0
}

// --------------------------------
// |     AttrOccupiedRowsCols     |
// --------------------------------

// AttrOccupiedRowsCols takes care of getting the number of occupied rows or
// columns for a player
type AttrOccupiedRowsCols struct {
	color Color // For which player rows/cols are counted
	rows  bool  // true: counting occupied rows; false: counting occupied columns
}

// GetAttributeName returns the name of an attribute
func (a AttrOccupiedRowsCols) GetAttributeName() string {
	n := "occ_"

	switch a.color {
	case Red:
		n += "red"
	case Blue:
		n += "blue"
	default:
		panic(fmt.Errorf("Invalid color %v", a.color))
	}

	n += "_"

	if a.rows {
		n += "rows"
	} else {
		n += "cols"
	}

	return n
}

// GetAttributeValue returns the value of an attribute
func (a AttrOccupiedRowsCols) GetAttributeValue(args *[]interface{}) int {
	patCount := (*args)[1].([2][]int)
	i := -1
	switch a.color {
	case Red:
		i = 0
	case Blue:
		i = 1
	default:
		panic(fmt.Errorf("Invalid color %v", a.color))
	}

	var r int
	if a.rows {
		r = patCount[i][len(patCount[i])-2]
	} else {
		r = patCount[i][len(patCount[i])-1]
	}

	return r
}

// ----------------------------
// |     AttrPatternCount     |
// ----------------------------

// AttrPatternCount takes care of the number of a single pattern on the board
// for one player
type AttrPatternCount struct {
	color        Color // For which player patterns are counted
	patternIndex int   // Index of the pattern
}

// GetAttributeName returns the name of an attribute
func (a AttrPatternCount) GetAttributeName() string {
	var n string

	switch a.color {
	case Red:
		n = "red_p"
	case Blue:
		n = "blue_p"
	default:
		panic(fmt.Errorf("Invalid color %v", a.color))
	}

	return fmt.Sprintf("%s%d", n, a.patternIndex)
}

// GetAttributeValue returns the value of an attribute
func (a AttrPatternCount) GetAttributeValue(args *[]interface{}) int {
	patCount := (*args)[1].([2][]int)
	i := -1
	switch a.color {
	case Red:
		i = 0
	case Blue:
		i = 1
	default:
		panic(fmt.Errorf("Invalid color %v", a.color))
	}

	return patCount[i][a.patternIndex]
}

// ------------------------------
// |     AttrLastPlayerTurn     |
// ------------------------------

// AttrLastPlayerTurn stores information about the last player
type AttrLastPlayerTurn struct {
	isLastPlayer bool
}

// GetAttributeName returns the name of an attribute
func (a AttrLastPlayerTurn) GetAttributeName() string {
	return "lp"
}

// GetAttributeValue returns 0 if the Red player had the last turn and 1 otherwise
func (a AttrLastPlayerTurn) GetAttributeValue(args *[]interface{}) int {
	s := (*args)[0].(State)
	lp := s.GetLastPlayer()
	switch {
	// Actual state
	case a.isLastPlayer && lp == Red:
		return 0
	case a.isLastPlayer && lp == Blue:
		return 1
	// State with reversed roles of players
	case !a.isLastPlayer && lp == Red:
		return 1
	case !a.isLastPlayer && lp == Blue:
		return 0
	}
	panic(fmt.Errorf("Invalid color %v", lp))
}

// ------------------------------------------
// |     AttrLastActionDistanceToCenter     |
// ------------------------------------------

// AttrLastActionDistanceToCenter tells how far from the center was the last
// stone placed
type AttrLastActionDistanceToCenter struct{}

// GetAttributeName returns the name of an attribute
func (a AttrLastActionDistanceToCenter) GetAttributeName() string {
	return "dtc"
}

// GetAttributeValue returns the value of an attribute
func (a AttrLastActionDistanceToCenter) GetAttributeValue(args *[]interface{}) int {
	s := (*args)[0].(State)
	lastAction := s.GetLastAction()
	x, y := lastAction.GetCoordinates()
	size := s.GetSize()
	var centerX, centerY int
	if size%2 == 1 {
		// Board has only one central position
		centerX = size / 2
		centerY = centerX
	} else {
		// Board has four central positions
		centerSmall, centerBig := size/2-1, size/2
		if x <= centerSmall {
			centerX = centerSmall
		} else {
			centerX = centerBig
		}
		if y <= centerSmall {
			centerY = centerSmall
		} else {
			centerY = centerBig
		}
		// (centerX, centerY) is the central position that is closest to (x, y)
	}
	return getDistanceBetween(centerX, centerY, x, y)
}

// getDistanceBetween returns the distance between points (x1, y1) and (x2, y2)
// in a hexagonal grid
func getDistanceBetween(x1, y1, x2, y2 int) int {
	return (abs(x1-x2) + abs(x1+y1-x2-y2) + abs(y1-y2)) / 2
}

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}
