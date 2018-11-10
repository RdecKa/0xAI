// Package hex - in this file all possible attributes that can be used for
// evaluating a hex board are listed.
//
// To add an attribute, do the following:
// 	- Implement a type implementing game.Attribute
// 	- Initialize instance(s) of that attribute
// 	- Add this/these instance(s) to the slice GenSamAttributes (together with
// 		matching opposite attribute)
// 	- In 2-ml/learn.py, select a type of the attribute when reading a CSV
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
	AttrLastPlayer         = AttrLastPlayerTurn{true}
	AttrLastPlayerOpponent = AttrLastPlayerTurn{false}
	AttrDistanceToCenter   = AttrLastActionDistanceToCenter{}

	AttrOccRedRows  = AttrOccupiedRowsCols{Red, true}
	AttrOccRedCols  = AttrOccupiedRowsCols{Red, false}
	AttrOccBlueRows = AttrOccupiedRowsCols{Blue, true}
	AttrOccBlueCols = AttrOccupiedRowsCols{Blue, false}

	AttrPatCountRed0  = AttrPatternCount{Red, 0}
	AttrPatCountRed1  = AttrPatternCount{Red, 1}
	AttrPatCountRed2  = AttrPatternCount{Red, 2}
	AttrPatCountRed3  = AttrPatternCount{Red, 3}
	AttrPatCountRed4  = AttrPatternCount{Red, 4}
	AttrPatCountRed5  = AttrPatternCount{Red, 5}
	AttrPatCountRed6  = AttrPatternCount{Red, 6}
	AttrPatCountRed7  = AttrPatternCount{Red, 7}
	AttrPatCountRed8  = AttrPatternCount{Red, 8}
	AttrPatCountRed9  = AttrPatternCount{Red, 9}
	AttrPatCountRed10 = AttrPatternCount{Red, 10}
	AttrPatCountRed11 = AttrPatternCount{Red, 11}
	AttrPatCountRed12 = AttrPatternCount{Red, 12}
	AttrPatCountRed13 = AttrPatternCount{Red, 13}
	AttrPatCountRed14 = AttrPatternCount{Red, 14}
	AttrPatCountRed15 = AttrPatternCount{Red, 15}
	AttrPatCountRed16 = AttrPatternCount{Red, 16}
	AttrPatCountRed17 = AttrPatternCount{Red, 17}
	AttrPatCountRed18 = AttrPatternCount{Red, 18}
	AttrPatCountRed19 = AttrPatternCount{Red, 19}

	AttrPatCountBlue0  = AttrPatternCount{Blue, 0}
	AttrPatCountBlue1  = AttrPatternCount{Blue, 1}
	AttrPatCountBlue2  = AttrPatternCount{Blue, 2}
	AttrPatCountBlue3  = AttrPatternCount{Blue, 3}
	AttrPatCountBlue4  = AttrPatternCount{Blue, 4}
	AttrPatCountBlue5  = AttrPatternCount{Blue, 5}
	AttrPatCountBlue6  = AttrPatternCount{Blue, 6}
	AttrPatCountBlue7  = AttrPatternCount{Blue, 7}
	AttrPatCountBlue8  = AttrPatternCount{Blue, 8}
	AttrPatCountBlue9  = AttrPatternCount{Blue, 9}
	AttrPatCountBlue10 = AttrPatternCount{Blue, 10}
	AttrPatCountBlue11 = AttrPatternCount{Blue, 11}
	AttrPatCountBlue12 = AttrPatternCount{Blue, 12}
	AttrPatCountBlue13 = AttrPatternCount{Blue, 13}
	AttrPatCountBlue14 = AttrPatternCount{Blue, 14}
	AttrPatCountBlue15 = AttrPatternCount{Blue, 15}
	AttrPatCountBlue16 = AttrPatternCount{Blue, 16}
	AttrPatCountBlue17 = AttrPatternCount{Blue, 17}
	AttrPatCountBlue18 = AttrPatternCount{Blue, 18}
	AttrPatCountBlue19 = AttrPatternCount{Blue, 19}
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
	[2]game.Attribute{AttrLastPlayer, AttrLastPlayerOpponent},
	[2]game.Attribute{AttrDistanceToCenter, nil},

	[2]game.Attribute{AttrOccRedRows, AttrOccBlueCols},
	[2]game.Attribute{AttrOccRedCols, AttrOccBlueRows},
	[2]game.Attribute{AttrOccBlueRows, AttrOccRedCols},
	[2]game.Attribute{AttrOccBlueCols, AttrOccRedRows},

	[2]game.Attribute{AttrPatCountRed0, AttrPatCountBlue0},
	[2]game.Attribute{AttrPatCountRed1, AttrPatCountBlue3},
	[2]game.Attribute{AttrPatCountRed2, AttrPatCountBlue2},
	[2]game.Attribute{AttrPatCountRed3, AttrPatCountBlue1},
	[2]game.Attribute{AttrPatCountRed4, AttrPatCountBlue4},
	[2]game.Attribute{AttrPatCountRed5, AttrPatCountBlue7},
	[2]game.Attribute{AttrPatCountRed6, AttrPatCountBlue6},
	[2]game.Attribute{AttrPatCountRed7, AttrPatCountBlue5},
	[2]game.Attribute{AttrPatCountRed8, AttrPatCountBlue10},
	[2]game.Attribute{AttrPatCountRed9, AttrPatCountBlue9},
	[2]game.Attribute{AttrPatCountRed10, AttrPatCountBlue8},
	[2]game.Attribute{AttrPatCountRed11, AttrPatCountBlue12},
	[2]game.Attribute{AttrPatCountRed12, AttrPatCountBlue11},
	[2]game.Attribute{AttrPatCountRed13, AttrPatCountBlue13},
	[2]game.Attribute{AttrPatCountRed14, AttrPatCountBlue15},
	[2]game.Attribute{AttrPatCountRed15, AttrPatCountBlue14},
	[2]game.Attribute{AttrPatCountRed16, AttrPatCountBlue16},
	[2]game.Attribute{AttrPatCountRed17, AttrPatCountBlue17},
	[2]game.Attribute{AttrPatCountRed18, AttrPatCountBlue19},
	[2]game.Attribute{AttrPatCountRed19, AttrPatCountBlue18},

	[2]game.Attribute{AttrPatCountBlue0, AttrPatCountRed0},
	[2]game.Attribute{AttrPatCountBlue1, AttrPatCountRed3},
	[2]game.Attribute{AttrPatCountBlue2, AttrPatCountRed2},
	[2]game.Attribute{AttrPatCountBlue3, AttrPatCountRed1},
	[2]game.Attribute{AttrPatCountBlue4, AttrPatCountRed4},
	[2]game.Attribute{AttrPatCountBlue5, AttrPatCountRed7},
	[2]game.Attribute{AttrPatCountBlue6, AttrPatCountRed6},
	[2]game.Attribute{AttrPatCountBlue7, AttrPatCountRed5},
	[2]game.Attribute{AttrPatCountBlue8, AttrPatCountRed10},
	[2]game.Attribute{AttrPatCountBlue9, AttrPatCountRed9},
	[2]game.Attribute{AttrPatCountBlue10, AttrPatCountRed8},
	[2]game.Attribute{AttrPatCountBlue11, AttrPatCountRed12},
	[2]game.Attribute{AttrPatCountBlue12, AttrPatCountRed11},
	[2]game.Attribute{AttrPatCountBlue13, AttrPatCountRed13},
	[2]game.Attribute{AttrPatCountBlue14, AttrPatCountRed15},
	[2]game.Attribute{AttrPatCountBlue15, AttrPatCountRed14},
	[2]game.Attribute{AttrPatCountBlue16, AttrPatCountRed16},
	[2]game.Attribute{AttrPatCountBlue17, AttrPatCountRed17},
	[2]game.Attribute{AttrPatCountBlue18, AttrPatCountRed19},
	[2]game.Attribute{AttrPatCountBlue19, AttrPatCountRed18},
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
	actionInQuestion := (*args)[2].(*Action)
	x, y := actionInQuestion.GetCoordinates()
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

	// If the player who is evaluating the state is not the player who's turn it
	// is, negate the value
	playerEvaluating := s.lastAction.c
	dist := getDistanceBetween(centerX, centerY, x, y)
	if playerEvaluating == actionInQuestion.c {
		return dist
	}
	return -dist
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
