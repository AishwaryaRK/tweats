package matcher

import "github.com/AishwaryaRK/tweats/datamodel"

// An Item is something we manage in a priority queue.
type Item struct {
	interest string            // The interest.
	tweeps   []datamodel.Tweep // The tweeps that has the interest.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the interest that has the least tweeps.
	return len(pq[i].tweeps) < len(pq[j].tweeps)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push pushes an item into the priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Item))
}

// Pop pops an item from the priority queue.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}
