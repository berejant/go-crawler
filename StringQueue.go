package main

import (
	"sync"
	"sync/atomic"
)

// StringQueueOffsetMax - When reach max - it needs cut processed items from Queue
const StringQueueOffsetMax = 255

type StringQueue struct {
	queue      []string
	queueMutex sync.Mutex
	Counter    uint32

	offset      int
	offsetMutex sync.Mutex
}

func NewStringQueue(capacity int) *StringQueue {
	return &StringQueue{
		queue:       make([]string, 0, capacity),
		queueMutex:  sync.Mutex{},
		offset:      0,
		offsetMutex: sync.Mutex{},
	}
}

func (queue *StringQueue) Add(val string) {
	queue.queueMutex.Lock()
	defer queue.queueMutex.Unlock()

	queue.queue = append(queue.queue, val)
	atomic.AddUint32(&queue.Counter, 1)
}

func (queue *StringQueue) GetNext() (val string) {
	if queue.offset < len(queue.queue) {
		queue.offsetMutex.Lock()
		defer queue.offsetMutex.Unlock()
		val = queue.queue[queue.offset]
		queue.offset++

		if queue.offset == StringQueueOffsetMax {
			queue.queueMutex.Lock()
			defer queue.queueMutex.Unlock()

			// rewrite slice (array) is operation with high amortization cost (because of memory). Do it only when reach limit
			// also we will have a lot write (append) operations into this slice, so do not block its
			queue.queue = queue.queue[queue.offset:]
			queue.offset = 0
		}
	}

	return val
}
