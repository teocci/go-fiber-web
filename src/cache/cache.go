// Package cache
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-12
package cache

import (
	"sync"
	"time"
)

type Entry struct {
	Data      interface{}
	ExpiresAt time.Time
}

type SimpleCache struct {
	data  map[string]Entry
	mutex sync.RWMutex
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{
		data: make(map[string]Entry),
	}
}

func (sc *SimpleCache) Set(key string, value interface{}, duration time.Duration) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	expiration := time.Now().Add(duration)
	sc.data[key] = Entry{
		Data:      value,
		ExpiresAt: expiration,
	}
}

func (sc *SimpleCache) Get(key string) (interface{}, bool) {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	entry, found := sc.data[key]
	if !found || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Data, true
}

func (sc *SimpleCache) Has(key string) bool {
	_, found := sc.Get(key)
	return found
}

func (sc *SimpleCache) Delete(key string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	delete(sc.data, key)
}

func (sc *SimpleCache) Clear() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.data = make(map[string]Entry)
}
