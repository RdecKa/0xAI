package hex

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// This file provides functions for pattern matching in hex grids. Patterns
// and grids must be represented as lists of integers, where one integer
// represents one row, and each two bits in an integer represent one column.

type pattern struct {
	w, h int     // width and heigth of the pattern
	pat  []uint8 // pattern
}

// CreatePatChecker creates a go routine that will serach for patterns in grids.
// It returns channels for communicatin with this goroutine.
func CreatePatChecker(fileName string) (chan []uint64, chan struct{}, chan [2][]int) {
	gridChan := make(chan []uint64, 1)
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

	patterns := make([][]*pattern, 0, 10)
	patC, rotC, lineC := -1, -1, 0 // Counters of patterns, rotations for each pattern, and lines in each pattern

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lineSplit := strings.Fields(line)
		if lineSplit[0] == "###" {
			patC++
			rotC = -1
			patterns = append(patterns, make([]*pattern, 0, 6))
		} else if lineSplit[0] == "---" {
			rotC++
			lineC = 0
			patterns[patC] = append(patterns[patC], &pattern{0, 0, make([]uint8, 0)})
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
func lineToNumber(lineSplit []string) (uint8, int) {
	var num uint8
	width := 0

	reverseLine(lineSplit)
	for _, ls := range lineSplit {
		width++

		num = num << 2
		switch ls {
		case ".":
			num += 0
		case "*":
			num += 3
		default:
			num += 2
		}
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
func countPatternInGrid(pat pattern, grid []uint64) (int, int, [2][]bool, [2][]bool) {
	countRed, countBlue := 0, 0

	var occRows [2][]bool
	var occCols [2][]bool
	occRows[0] = make([]bool, len(grid))
	occRows[1] = make([]bool, len(grid))
	occCols[0] = make([]bool, len(grid))
	occCols[1] = make([]bool, len(grid))

	for yStart := 0; yStart <= len(grid)-pat.h; yStart++ {
		for xStart := 0; xStart <= len(grid)-pat.w; xStart++ {
			found := -1
			if matches(pat, grid, xStart, yStart, Red) {
				countRed++
				found = 0
			} else if matches(pat, grid, xStart, yStart, Blue) {
				countBlue++
				found = 1
			}
			if found >= 0 {
				for x := xStart; x < xStart+pat.w; x++ {
					occCols[found][x] = true
				}
				for y := yStart; y < yStart+pat.h; y++ {
					occRows[found][y] = true
				}
			}
		}
	}

	return countRed, countBlue, occRows, occCols
}

// matches checks whether a subgrid matches the given pattern.
func matches(pat pattern, grid []uint64, xStart, yStart int, player Color) bool {
	for y := 0; y < pat.h; y++ {
		rowGrid := grid[yStart+y]
		rowPat := pat.pat[y]
		for x := 0; x < pat.w; x++ {
			cellGrid := (rowGrid >> (2 * uint(xStart+x))) & 3
			cellPat := (rowPat >> (2 * uint(x))) & 3

			switch cellPat {
			case 0: // Empty cell in the pattern
				if GetColorFromBits(cellGrid) != None {
					return false
				}
			case 3: // Marked cell in the pattern
				if GetColorFromBits(cellGrid) != player {
					return false
				}
			}
		}
	}
	return true
}

// patChecker is a goroutine that searches for patterns in grids, sent via
// gridChan. Results are sent via resultChan. stopChan is used to end the
// goroutine.
// It also checks in how many rows and columns each player has at least one
// virtual connection
func patChecker(filename string, gridChan chan []uint64, stopChan chan struct{}, resultChan chan [2][]int) {
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
