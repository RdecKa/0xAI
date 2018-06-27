package hex

import (
	"bufio"
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
func CreatePatChecker() (chan []uint64, chan struct{}, chan [2][]int) {
	gridChan := make(chan []uint64, 1)
	stopChan := make(chan struct{}, 1)
	resultChan := make(chan [2][]int, 1)

	go patChecker(gridChan, stopChan, resultChan)

	return gridChan, stopChan, resultChan
}

// readPatternsFromFile reads all patterns from a specified file and returns a
// 2D slice of patterns. The first dimension is a pattern, the second dimension
// are all rotations of that pattern.
func readPatternsFromFile() ([][]*pattern, error) {
	f, err := os.Open("mcts/hex/patterns.txt")
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
// rotation) has in the grid.
func countPatternInGrid(pat pattern, grid []uint64) (int, int) {
	countRed, countBlue := 0, 0

	for yStart := 0; yStart <= len(grid)-pat.h; yStart++ {
		for xStart := 0; xStart <= len(grid)-pat.w; xStart++ {
			if matches(pat, grid, xStart, yStart, Red) {
				countRed++
			} else if matches(pat, grid, xStart, yStart, Blue) {
				countBlue++
			}
		}
	}

	return countRed, countBlue
}

// matches checks whether a subgrid matches the given pattern.
func matches(pat pattern, grid []uint64, xStart, yStart int, player color) bool {
	for y := 0; y < pat.h; y++ {
		rowGrid := grid[yStart+y]
		rowPat := pat.pat[y]
		for x := 0; x < pat.w; x++ {
			cellGrid := (rowGrid >> (2 * uint(xStart+x))) & 3
			cellPat := (rowPat >> (2 * uint(x))) & 3

			switch cellPat {
			case 0: // Empty cell in the pattern
				if getColorFromBits(cellGrid) != None {
					return false
				}
			case 3: // Marked cell in the pattern
				if getColorFromBits(cellGrid) != player {
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
func patChecker(gridChan chan []uint64, stopChan chan struct{}, resultChan chan [2][]int) error {
	defer close(gridChan)
	defer close(stopChan)
	defer close(resultChan)

	patterns, err := readPatternsFromFile()
	if err != nil {
		return err
	}

	for {
		select {
		case grid := <-gridChan:
			var results [2][]int
			results[0] = make([]int, len(patterns)) // Counts for red
			results[1] = make([]int, len(patterns)) // Counts for blue
			for pi, p := range patterns {
				for _, r := range p {
					red, blue := countPatternInGrid(*r, grid)
					results[0][pi] += red
					results[1][pi] += blue
				}
			}
			resultChan <- results
		case <-stopChan:
			return nil
		}
	}
}
