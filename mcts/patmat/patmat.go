// Package patmat provides functions for pattern matching in grids. Patterns and
// grids must be represented as lists of integers, where one integer represents
// one row, and each two bits in an integer represent one column.
package patmat

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type pattern struct {
	w, h int     // width and heigth of the pattern
	pat  []uint8 // pattern
}

// ReadPatterns reads grid patterns from a file
func ReadPatterns() {
	patterns, _ := readPatternsFromFile()
	for i, p := range patterns {
		fmt.Printf("Pattern %d\n", i)
		for _, c := range p {
			fmt.Printf("---> %v\n", c)
		}
	}
}

func readPatternsFromFile() ([][]*pattern, error) {
	f, err := os.Open("mcts/patmat/patterns.txt")
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

func reverseLine(line []string) {
	for i, j := 0, len(line)-1; i < j; i, j = i+1, j-1 {
		line[i], line[j] = line[j], line[i]
	}
}
