package main

import (
	"fmt"
	"strconv"

	"github.com/RdecKa/mcts/tree"
)

// BasicNodeValue is a simple node value, used for testing
type BasicNodeValue struct {
	s string
}

func (b BasicNodeValue) String() string {
	return b.s
}

// IntNodeValue is another simple node value, used for testing
type IntNodeValue struct {
	i int
}

func (i IntNodeValue) String() string {
	return strconv.Itoa(i.i)
}

func main() {
	t := tree.NewTree(tree.NewNode(BasicNodeValue{"Mudkip"}))
	t.GetRoot().AddChild(tree.NewNode(BasicNodeValue{"Lapras"}))
	t.GetRoot().AddChild(tree.NewNode(IntNodeValue{12345}))
	fmt.Printf("%s", t)
	for i, c := range t.GetRoot().GetChildren() {
		c.AddChild(tree.NewNode(BasicNodeValue{"Pikachu"}))
		c.AddChild(tree.NewNode(BasicNodeValue{"Chikorita"}))
		c.AddChild(tree.NewNode(IntNodeValue{(i + 1) * 1000}))
		c.GetChildren()[1].AddChild(tree.NewNode(IntNodeValue{5}))
	}
	fmt.Printf("%s", t)
}
