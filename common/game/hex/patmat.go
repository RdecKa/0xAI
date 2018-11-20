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
	lspRow   []int        // longest suffix-prefix for rows (KMP algorithm)
}

// return values:
//	-1: columns are different
//	0: non-indefinite cells in columns are the same
//	1: columns are the same
func (p *pattern) columnsSame(c1, c2 int) int {
	possiblySame := false
	for _, row := range p.pat {
		if row[c1] != row[c2] {
			if row[c1] == cellIndefinite || row[c2] == cellIndefinite {
				possiblySame = true
			} else {
				return -1
			}
		}
	}
	if possiblySame {
		return 0
	}
	return 1
}

func (p *pattern) rowsSame(r1, r2 int) int {
	possiblySame := false
	row1 := p.pat[r1]
	row2 := p.pat[r2]
	for col := 0; col < len(p.pat[r1]); col++ {
		if row1[col] != row2[col] {
			if row1[col] == cellIndefinite || row2[col] == cellIndefinite {
				possiblySame = true
			} else {
				return -1
			}
		}
	}
	if possiblySame {
		return 0
	}
	return 1
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
	return fmt.Sprintf("%ssize: (%d, %d), excluded: %v\nlspRow: %v\n",
		s, p.w, p.h, p.excluded, p.lspRow)
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
func CreatePatChecker(fileName string) (chan []uint32, chan struct{}, chan [2][]int) {
	gridChan := make(chan []uint32, 1)
	stopChan := make(chan struct{}, 1)
	resultChan := make(chan [2][]int, 1)

	go patChecker(fileName, gridChan, stopChan, resultChan)

	return gridChan, stopChan, resultChan
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
			patterns[patC] = append(patterns[patC], &pattern{0, 0, make([][]cellType, 0, 3), exclude, nil})
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

	// Calculate LSP tables
	for _, p := range patterns {
		for _, r := range p {
			lspRow := make([]int, r.h)
			for d := 1; d < len(lspRow); d++ {
				j := lspRow[d-1]
				for j > 0 && r.rowsSame(d, j) == -1 {
					j--
				}
				if r.rowsSame(d, j) > -1 {
					j++
				}
				lspRow[d] = j
			}
			r.lspRow = lspRow
		}
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

// reverseLine reverses a slice of strings in-place.
func reverseLine(line []string) {
	for i, j := 0, len(line)-1; i < j; i, j = i+1, j-1 {
		line[i], line[j] = line[j], line[i]
	}
}

// countPatternInGrid counts how many occurences the given pattern (with given
// rotation) has in the grid. It also counts how many rows and columns each
// player has occupied.
func countPatternInGrid(pat pattern, grid []uint32) (int, int, [2][]bool, [2][]bool) {
	countRed, countBlue := 0, 0

	var occRows [2][]bool
	var occCols [2][]bool
	occRows[0] = make([]bool, len(grid))
	occRows[1] = make([]bool, len(grid))
	occCols[0] = make([]bool, len(grid))
	occCols[1] = make([]bool, len(grid))

	for xStart := 0; xStart <= len(grid)-pat.w; xStart++ {
		for yStart := 0; yStart <= len(grid)-pat.h; {
			found := -1
			matchRow, c := matches(pat, grid, xStart, yStart)
			if matchRow == pat.h {
				switch c {
				case Red:
					countRed++
					found = 0
				case Blue:
					countBlue++
					found = 1
				}
			}
			if found >= 0 && !pat.excluded {
				for x := xStart; x < xStart+pat.w; x++ {
					occCols[found][x] = true
				}
				for y := yStart; y < yStart+pat.h; y++ {
					occRows[found][y] = true
				}
			}
			if matchRow == 0 {
				yStart++
			} else {
				yStart += matchRow - pat.lspRow[matchRow-1]
			}
		}
	}

	return countRed, countBlue, occRows, occCols
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
func patChecker(filename string, gridChan chan []uint32, stopChan chan struct{}, resultChan chan [2][]int) {
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

			for pi, p := range patterns {
				for _, r := range p {
					red, blue, occR, occC := countPatternInGrid(*r, grid)
					results[0][pi] += red
					results[1][pi] += blue

					for x := 0; x < len(occC[0]); x++ {
						occCols[0][x] = occCols[0][x] || occC[0][x]
						occCols[1][x] = occCols[1][x] || occC[1][x]
					}
					for y := 0; y < len(occR[0]); y++ {
						occRows[0][y] = occRows[0][y] || occR[0][y]
						occRows[1][y] = occRows[1][y] || occR[1][y]
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

			resultChan <- results
		case <-stopChan:
			return
		}
	}
}
