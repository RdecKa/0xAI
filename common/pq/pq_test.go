package pq

import (
	"container/heap"
	"testing"
)

func TestPQ(t *testing.T) {
	p := New(50)

	input := []int{5, 9, 7, 2, 1, 1, 10}
	output := []int{1, 1, 2, 5, 7, 9, 10}

	for _, el := range input {
		i := NewItem(el, el)
		heap.Push(p, i)
	}

	for _, el := range output {
		val := heap.Pop(p)
		if val != el {
			t.Fatalf("Expected %v, got %v", el, val)
		}
	}

}
