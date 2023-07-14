package main

import (
	"github.com/VictoriaMetrics/fastcache"
)

// Set thread safe set of bytes
type Set struct {
	memory *fastcache.Cache
}

func NewSet(maxBytes int) *Set {
	return &Set{
		memory: fastcache.New(maxBytes),
	}
}

func (set *Set) Add(value []byte) bool {
	// check before get lock.
	if set.memory.Has(value) {
		return false
	}

	set.memory.Set(value, []byte{})
	return true
}
