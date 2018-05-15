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
func (mcts *MCTS) GenSamples(folder, filePrefix string, treasholdN uint) ([]*tree.Node, error) {
	// Create a new file
	t := time.Now()
	fileName := folder + filePrefix + "_" + t.Format("2006-01-02T15:04:05") + ".in"
	fmt.Printf("Writing samples to file %s ... ", fileName)
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Write samples to a file
	root := mcts.mcTree.GetRoot()
	expandCandidates := genSamples(root, f, treasholdN)

	fmt.Println("Done!")
	return expandCandidates, nil
}

// genSamples traverses the MCTS tree starting from Node node and writes samples
// to a File file
func genSamples(node *tree.Node, file *os.File, treasholdN uint) []*tree.Node {
	mnv := node.GetValue().(*mctsNodeValue)
	expandCandidates := make([]*tree.Node, 0, 20)
	if mnv.n >= treasholdN {
		file.WriteString(mnv.genSample())
		for _, c := range node.GetChildren() {
			g := genSamples(c, file, treasholdN)
			expandCandidates = append(expandCandidates, g...)
		}
	} else {
		// Add the node to a list of nodes that will possibly be expanded in the
		// following MCTS
		expandCandidates = append(expandCandidates, node)
	}
	return expandCandidates
}

// genSample returns a string representation of a single sample
func (mnv *mctsNodeValue) genSample() string {
	s := mnv.state.(hex.State)
	red, blue, empty := s.GetNumOfStones()
	return fmt.Sprintf("%f,%d,%d,%d\n", mnv.q, red, blue, empty)
}
