package mcts

import (
	"os"

	"github.com/RdecKa/bachleor-thesis/common/tree"
)

// GenSamples traverses the MCTS tree and writes samples (nodes that have been
// visited at least treasholdN times) to an outputFile. It returns possible
// candidates for later MCTS
func (mcts *MCTS) GenSamples(outputFile *os.File, treasholdN uint, gridChan chan []uint32,
	patChan chan []int, resultChan chan [2][]int) ([]*tree.Node, error) {

	// Write samples to a file
	root := mcts.mcTree.GetRoot()
	expandCandidates := genSamples(root, outputFile, treasholdN, gridChan, patChan, resultChan)
	return expandCandidates, nil
}

// genSamples traverses the MCTS tree starting from Node node and writes samples
// to a File file. It returns possible candidates for later MCTS
func genSamples(node *tree.Node, outputFile *os.File, treasholdN uint, gridChan chan []uint32,
	patChan chan []int, resultChan chan [2][]int) []*tree.Node {
	mnv := node.GetValue().(*mctsNodeValue)
	expandCandidates := make([]*tree.Node, 0, 20)
	if mnv.n >= treasholdN {
		outputFile.WriteString(mnv.state.GenSample(mnv.q, gridChan, patChan, resultChan))
		for _, c := range node.GetChildren() {
			g := genSamples(c, outputFile, treasholdN, gridChan, patChan, resultChan)
			expandCandidates = append(expandCandidates, g...)
		}
	} else {
		// Add the node to a list of nodes that will possibly be expanded in the
		// following MCTS
		expandCandidates = append(expandCandidates, node)
	}
	return expandCandidates
}
