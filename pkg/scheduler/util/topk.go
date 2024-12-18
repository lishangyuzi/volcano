package util

import (
	"container/heap"
	"volcano.sh/volcano/pkg/scheduler/api"
)

type ScoreNode struct {
	Score float64
	Node  *api.NodeInfo
}

type MaxHeap []ScoreNode

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].Score > h[j].Score }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(ScoreNode))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	if n < 1 {
		return nil
	}
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type TopK struct {
	k    int
	data *MaxHeap
}

func NewTopK(k int) *TopK {
	h := &MaxHeap{}
	heap.Init(h)
	return &TopK{k: k, data: h}
}

func (t *TopK) Add(scoreNode ScoreNode) {
	if t.data.Len() < t.k {
		heap.Push(t.data, scoreNode)
	} else if scoreNode.Score > (*t.data)[0].Score {
		heap.Pop(t.data)
		heap.Push(t.data, scoreNode)
	}
}

func (t *TopK) GetTopK() (ScoreNode, bool) {
	if t.data.Len() < 1 {
		return ScoreNode{}, false
	}
	return heap.Pop(t.data).(ScoreNode), true
}

func (t *TopK) IsEmpty() bool {
	return t.data.Len() == 0
}
