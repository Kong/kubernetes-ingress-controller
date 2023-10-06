package util

// StringHeap is an array of strings that implements container/heap.Interface.
type StringHeap []string

func (h StringHeap) Len() int {
	return len(h)
}

func (h StringHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h StringHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *StringHeap) Push(s any) {
	*h = append(*h, s.(string))
}

func (h *StringHeap) Pop() any {
	old := *h
	n := len(old)
	s := old[n-1]
	*h = old[0 : n-1]
	return s
}
