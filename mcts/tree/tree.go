package tree

import (
	"fmt"
)

// ---------------------
// |     NodeValue     |
// ---------------------

// NodeValue is an iterface for a value that can be in the node
type NodeValue interface {
	String() string
}

// ----------------
// |     Node     |
// ----------------

// Node represents a node in a tree
type Node struct {
	children []*Node
	value    NodeValue
}

// NewNode creates a new Node with value v and returns a pointer to that node
func NewNode(v NodeValue) *Node {
	return &Node{make([]*Node, 0), v}
}

func (n Node) String() string {
	s := ""
	s += n.value.String()
	s += "("
	for i, c := range n.children {
		if i > 0 {
			s += ", "
		}
		s += c.String()
	}
	s += ")"
	return s
}

// GetChildren returns list of node n's successors
func (n *Node) GetChildren() []*Node {
	return n.children
}

// GetValue returns node n's value
func (n *Node) GetValue() NodeValue {
	return n.value
}

// AddChild adds a child newNode to Node n (to the end of the list of children)
func (n *Node) AddChild(newNode *Node) {
	n.children = append(n.children, newNode)
}

// ----------------
// |     Tree     |
// ----------------

// Tree represents a tree
type Tree struct {
	root *Node
}

// NewTree creates a tree with the given root node
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
