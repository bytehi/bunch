package bunch

import (
	"container/heap"
	"time"
)

// heapNode 是堆中的元素
type heapNode struct {
	interval      time.Duration
	executionTime time.Time
	callback      func()
	canceled      bool
}

// 实现 minHeap.Interface 接口的 minHeap
type minHeap []*heapNode

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].executionTime.Before(h[j].executionTime) }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(*heapNode))
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Peek 返回堆顶元素，但不弹出它
func (h *minHeap) Peek() interface{} {
	return (*h)[0]
}

type Timer struct {
	h minHeap
}

func (t *Timer) Add(duration time.Duration, callback func()) func() {
	node := &heapNode{
		interval:      duration,
		executionTime: time.Now().Add(duration),
		callback:      callback,
	}
	heap.Push(&t.h, node)
	return func() {
		node.canceled = true
	}
}

func (t *Timer) Timeout(now time.Time) {
	for {
		first := t.h.Peek()
		if first == nil {
			break
		}
		node := first.(*heapNode)
		if node.executionTime.After(now) {
			break
		}
		heap.Pop(&t.h)
		if node.canceled {
			continue
		}
		node.callback()
		node.executionTime = now.Add(node.interval)
		heap.Push(&t.h, node)
	}
}
