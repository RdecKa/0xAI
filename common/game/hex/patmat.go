package hex

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// This file provides functions for pattern matching in hex grids. Grids must be
// represented as lists of integers, where one integer represents one row, and
// each two bits in an integer represent one column. Patterns must be
// represented as 2D slices.

// -------------------
// |     pattern     |
// -------------------

type pattern struct {
	w, h     int          // width and heigth of the pattern
	pat      [][]cellType // pattern
	excluded bool         // true if rows and columns where this pattern is found do not count as occupied, false otherwise
}

func (p *pattern) String() string {
	s := ""
	for ri, r := range p.pat {
		for i := 0; i < ri; i++ {
			s += " "
		}
		for _, c := range r {
			s += c.String() + " "
		}
		s += "\n"
	}
	return fmt.Sprintf("%ssize: (%d, %d), excluded: %v\n",
		s, p.w, p.h, p.excluded)
}

// --------------------
// |     cellType     |
// --------------------

type cellType byte

// enum for player types
const (
	cellEmpty      cellType = 0
	cellPlayer     cellType = 1
	cellOpponent   cellType = 2
	cellIndefinite cellType = 3
)

func (ct cellType) String() string {
	switch ct {
	case cellEmpty:
		return "."
	case cellPlayer:
		return "*"
	case cellOpponent:
		return "/"
	case cellIndefinite:
		return "?"
	default:
		return "?"
	}
}

func getCellTypeFromString(s string) cellType {
	switch s {
	case ".": // Empty cell
		return cellEmpty
	case "/": // Opponent's color
		return cellOpponent
	case "?": // Cell state not important
		return cellIndefinite
	case "*": // Player's color
		return cellPlayer
	default:
		fmt.Println(fmt.Errorf("Invalid character '%s' in pattern", s))
		return cellIndefinite
	}
}

// ----------------------
// |     PatChecker     |
// ----------------------

// CreatePatChecker creates a go routine that will serach for patterns in grids.
// It returns channels for communicatin with this goroutine.
func CreatePatChecker(fileName string) (chan []uint32, chan []int, chan struct{}, chan [2][]int) {
	gridChan := make(chan []uint32, 1)
	patChan := make(chan []int, 1)
	stopChan := make(chan struct{}, 1)
	resultChan := make(chan [2][]int, 1)

	go patChecker(fileName, gridChan, patChan, stopChan, resultChan)

	return gridChan, patChan, stopChan, resultChan
}

// readPatternsFromFile reads all patterns from a specified file and returns a
// 2D slice of patterns. The first dimension is a pattern, the second dimension
// are all rotations of that pattern.
func readPatternsFromFile(fileName string) ([][]*pattern, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	patterns := make([][]*pattern, 0, 25)
	patC, rotC, lineC := -1, -1, 0 // Counters of patterns, rotations for each pattern, and lines in each pattern
	exclude := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lineSplit := strings.Fields(line)
		if lineSplit[0] == "###" {
			exclude = false
			patC++
			rotC = -1
			patterns = append(patterns, make([]*pattern, 0, 1))
		} else if lineSplit[0] == "---" {
			rotC++
			lineC = 0
			patterns[patC] = append(patterns[patC], &pattern{0, 0, make([][]cellType, 0, 3), exclude})
		} else if lineSplit[0] == "exclude" {
			exclude = true
		} else {
			val, w := lineToNumber(lineSplit)
			patterns[patC][rotC].pat = append(patterns[patC][rotC].pat, val)
			patterns[patC][rotC].w = w // Necessary only once, but easier that way
			patterns[patC][rotC].h++
			lineC++
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}

// lineToNumber converts ASCII characters that represent a line in a format
// needed.
func lineToNumber(lineSplit []string) ([]cellType, int) {
	num := make([]cellType, len(lineSplit))
	width := len(num)

	for i, ls := range lineSplit {
		num[i] = getCellTypeFromString(ls)
	}

	return num, width
}

// countPatternsInGrid counts how many occurences the given pattern (with given
// rotation) has in the grid. It also counts how many rows and columns each
// player has occupied.
func countPatternsInGrid(patterns [][]*pattern, grid []uint32, usedPat []int) [2][]int {
	var results [2][]int
	// Last two numbers mean number of rows and columns (respectively) occupied by a player
	results[0] = make([]int, len(patterns)+2) // Counts for red
	results[1] = make([]int, len(patterns)+2) // Counts for blue

	var occRows [2][]bool // true indicates that at least one virtual connection is in this row
	var occCols [2][]bool
	occRows[0] = make([]bool, len(grid)) // Occupied rows of the red player
	occRows[1] = make([]bool, len(grid)) // Occupied rows of the blue player
	occCols[0] = make([]bool, len(grid)) // Occupied columns of the red player
	occCols[1] = make([]bool, len(grid)) // Occupied columns of the blue player

	patChecked := 0
	for pi, p := range patterns {
		if usedPat != nil {
			if patChecked >= len(usedPat) {
				break
			} else if usedPat[patChecked] == pi {
				patChecked++
			} else {
				continue
			}
		}
		for xStart := 0; xStart <= len(grid); xStart++ {
			for yStart := 0; yStart <= len(grid); yStart++ {
				for _, r := range p {
					if xStart+r.w > len(grid) || yStart+r.h > len(grid) {
						continue
					}
					found := -1
					matchRow, c := matches(*r, grid, xStart, yStart)
					if matchRow == r.h {
						switch c {
						case Red:
							found = 0
						case Blue:
							found = 1
						}
					}
					if found >= 0 {
						results[found][pi]++
						if !r.excluded {
							for x := xStart; x < xStart+r.w; x++ {
								occCols[found][x] = true
							}
							for y := yStart; y < yStart+r.h; y++ {
								occRows[found][y] = true
							}
						}
					}
				}
			}
		}
	}

	// Last two numbers are for counting rows and columns with at least
	// one virtual connection
	for s := 0; s < len(grid); s++ {
		if occRows[0][s] {
			results[0][len(patterns)]++
		}
		if occRows[1][s] {
			results[1][len(patterns)]++
		}
		if occCols[0][s] {
			results[0][len(patterns)+1]++
		}
		if occCols[1][s] {
			results[1][len(patterns)+1]++
		}
	}

	return results
}

// matches checks whether a subgrid matches the given pattern.
// The first return value tells how many rows did match the pattern.
// The last value tells which player has a match.
func matches(pat pattern, grid []uint32, xStart, yStart int) (int, Color) {
	possibleRed, possibleBlue := true, true
	for y := 0; y < pat.h; y++ {
		rowGrid := grid[yStart+y] >> (2 * uint(xStart))
		rowPat := pat.pat[y]
		for x := 0; x < pat.w; x++ {
			cellGrid := rowGrid & 3
			cellPat := rowPat[x]
			cellColor := GetColorFromBits(cellGrid)

			switch cellPat {
			case cellPlayer: // Marked cell in the pattern (player's color)
				if cellColor == Red {
					possibleBlue = false
				} else if cellColor == Blue {
					possibleRed = false
				} else {
					return y, None
				}
			case cellEmpty: // Empty cell in the pattern
				if cellColor != None {
					return y, None
				}
			case cellOpponent: // Opponent
				if cellColor == Red {
					possibleRed = false
				} else if cellColor == Blue {
					possibleBlue = false
				} else {
					// The cell is empty -> no player can match
					return y, None
				}
			}

			if !possibleRed && !possibleBlue {
				return y, None
			}

			rowGrid = rowGrid >> 2
		}
	}

	if possibleRed && possibleBlue {
		panic("Both players match a pattern")
	}
	if possibleRed {
		return pat.h, Red
	} else if possibleBlue {
		return pat.h, Blue
	}
	return 0, None
}

// patChecker is a goroutine that searches for patterns in grids, sent via
// gridChan. Results are sent via resultChan. stopChan is used to end the
// goroutine.
// It also checks in how many rows and columns each player has at least one
// stoen or virtual connection
func patChecker(filename string, gridChan chan []uint32, patChan chan []int, stopChan chan struct{}, resultChan chan [2][]int) {
	defer close(gridChan)
	defer close(stopChan)
	defer close(resultChan)

	patterns, err := readPatternsFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case grid := <-gridChan:
			usedPats := <-patChan
			results := countPatternsInGrid(patterns, grid, usedPats)
			resultChan <- results
		case <-stopChan:
			return
		}
	}
}
