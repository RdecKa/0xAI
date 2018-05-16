package mcts

import (
	"fmt"
	"os"

	"github.com/RdecKa/mcts/hex"
	"github.com/RdecKa/mcts/tree"
)

// GenSamples traverses the MCTS tree and writes samples (nodes that have been
// visited at least treasholdN times) to an outputFile. It returns possible
// candidates for later MCTS
func (mcts *MCTS) GenSamples(outputFile *os.File, treasholdN uint) ([]*tree.Node, error) {
	// Write samples to a file
	root := mcts.mcTree.GetRoot()
	expandCandidates := genSamples(root, outputFile, treasholdN)
	return expandCandidates, nil
}

// genSamples traverses the MCTS tree starting from Node node and writes samples
// to a File file. It returns possible candidates for later MCTS
func genSamples(node *tree.Node, outputFile *os.File, treasholdN uint) []*tree.Node {
	mnv := node.GetValue().(*mctsNodeValue)
	expandCandidates := make([]*tree.Node, 0, 20)
	if mnv.n >= treasholdN {
		outputFile.WriteString(mnv.genSample())
		for _, c := range node.GetChildren() {
			g := genSamples(c, outputFile, treasholdN)
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
