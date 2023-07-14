package main

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSet_Add(t *testing.T) {
	t.Run("race_condition", func(t *testing.T) {
		concurrency := 30
		value := []byte{'t', 'e', 's', 't', '-', 'v', 'a', 'l', 'u', 'e'}

		set := NewSet(1024 * 1024 * 10)
		wg := &sync.WaitGroup{}

		trueCounter := uint32(0)
		falseCounter := uint32(0)

		timeToRun := time.Now().Add(time.Millisecond * 100)

		addToSet := func() {
			time.Until(timeToRun)
			if set.Add(value) {
				atomic.AddUint32(&trueCounter, 1)
			} else {
				atomic.AddUint32(&falseCounter, 1)
			}
			wg.Done()
		}

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go addToSet()
		}
		wg.Wait()
		assert.Equal(t, uint32(1), trueCounter)
		assert.Equal(t, uint32(concurrency-1), falseCounter)
	})
}
