// Package pq provides a priority queue, implemented as a heap
package pq

import (
	"container/heap"
)

// ----------------
// |     Item     |
// ----------------

// Item is an element of a priority queue.
type Item struct {
	priority int
	index    int
	value    interface{}
}

// NewItem creates a new item with given priority, index and value
func NewItem(priority int, value interface{}) *Item {
	return &Item{
		priority: priority,
		value:    value,
	}
}

// UpdatePriority updates priorty of Item i in PriorityQueue pq
func (i *Item) UpdatePriority(newPriority int, pq *PriorityQueue) {
	i.priority = newPriority
	heap.Fix(pq, i.index)
}

// GetValue returns value of Item i
func (i *Item) GetValue() interface{} {
	return i.value
}

// GetPriority returns priority of Item i
func (i *Item) GetPriority() int {
	return i.priority
}

// -------------------------
// |     PriorityQueue     |
// -------------------------

// PriorityQueue is stored as a list of Items
type PriorityQueue []*Item

// New returns an empty priority queue with initial size 0 and reserved size
// reservedSize
func New(reservedSize int) *PriorityQueue {
	pq := make(PriorityQueue, 0, reservedSize)
	heap.Init(&pq)
	return &pq
}

// Len returns number of elements in the priority queue pq
func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

// Swap swaps elements with indices i and j in pq
func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

// Less returns true if priority of element on index i is smaller than priority
// of element on index j
func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].priority < (*pq)[j].priority
}

// Push adds el to the priority queue pq
func (pq *PriorityQueue) Push(el interface{}) {
	item := el.(*Item) // Ensure the right type of el
	item.index = len(*pq)
	*pq = append(*pq, item)
}

// Pop returns the first element is the priority queue pq
func (pq *PriorityQueue) Pop() interface{} {
	hLen := len(*pq)
	item := (*pq)[hLen-1]
	*pq = (*pq)[0 : hLen-1]
	return item.value
}

// Contains checks whether PriorityQueue pq contains element el, using function
// same for comparison
func (pq *PriorityQueue) Contains(el interface{}, same func(a, b interface{}) bool) *Item {
	for _, item := range *pq {
		i := item.value
		if same(el, i) {
			return item
		}
	}
	return nil
}
