package tree

import (
	"fmt"
)

// ----------------
// |     Node     |
// ----------------

// Node represents a node in a tree
type Node struct {
	children []*Node
	value    interface{}
}

// NewNode creates a new Node with value v and returns a pointer to that node
func NewNode(v interface{}) *Node {
	return &Node{nil, v}
}

func (n *Node) String() string {
	s := ""
	for i, c := range (*n).children {
		if i > 0 {
			s += ", "
		}
		s += c.String()
	}

	return fmt.Sprintf("%v (%s)", n.value, s)
}

// GetChildren returns list of node n's successors
func (n *Node) GetChildren() []*Node {
	return n.children
}

// GetValue returns node n's value
func (n *Node) GetValue() interface{} {
	return n.value
}

// SetChildren sets children of a node n
func (n *Node) SetChildren(children []*Node) {
	n.children = children
}

// ----------------
// |     Tree     |
// ----------------

// Tree represents a tree
type Tree struct {
	root *Node
}

// NewTree creates a tree with given root node
func NewTree(root *Node) *Tree {
	return &Tree{root}
}

func (t Tree) String() string {
	return fmt.Sprintf("%v\n", t.root)
}

// GetRoot returns root node of the tree
func (t Tree) GetRoot() *Node {
	return t.root
}
