package main

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestStringQueue_Add(t *testing.T) {
	t.Run("cut_queue", func(t *testing.T) {
		const count = StringQueueOffsetMax + 10
		queue := NewStringQueue(count + 10)

		for i := 0; i < count; i++ {
			queue.Add("string" + strconv.Itoa(i))
		}

		for i := 0; i < count; i++ {
			if !assert.Equal(t, "string"+strconv.Itoa(i), queue.GetNext()) {
				break
			}
		}
	})
}
