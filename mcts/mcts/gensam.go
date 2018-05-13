package mcts

import (
	"fmt"
	"os"
	"time"

	"github.com/RdecKa/mcts/hex"
	"github.com/RdecKa/mcts/tree"
)

// GenSamples traverses the MCTS tree and writes samples (nodes that have been
// visited at least treasholdN times) to a new file
func (mcts *MCTS) GenSamples(folder, filePrefix string, treasholdN uint) error {
	// Create a new file
	t := time.Now()
	fileName := folder + filePrefix + "_" + t.Format("2006-01-02T15:04:05") + ".in"
	fmt.Printf("Writing samples to file %s ... ", fileName)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write samples to a file
	root := mcts.mcTree.GetRoot()
	genSamples(root, f, treasholdN)

	fmt.Println("Done!")
	return nil
}

// genSamples traverses the MCTS tree starting from Node node and writes samples
// to a File file
func genSamples(node *tree.Node, file *os.File, treasholdN uint) {
	mnv := node.GetValue().(*mctsNodeValue)
	if mnv.n >= treasholdN {
		file.WriteString(mnv.genSample())
		for _, c := range node.GetChildren() {
			genSamples(c, file, treasholdN)
		}
	}
}

// genSample returns a string representation of a single sample
func (mnv *mctsNodeValue) genSample() string {
	s := mnv.state.(hex.State)
	red, blue, empty := s.GetNumOfStones()
	return fmt.Sprintf("%f,%d,%d,%d\n", mnv.q, red, blue, empty)
}
