// Package scache
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-12
package scache

import (
	"sync"
	"time"
)

type Entry[T any] struct {
	Data      T
	ExpiresAt time.Time
}

type SimpleCache[T any] struct {
	data            map[string]Entry[T]
	mutex           sync.RWMutex
	defaultDuration time.Duration
}

type CacheInstance[T any] struct {
	instance *SimpleCache[T]
	once     sync.Once
}

var cacheInstances = make(map[string]interface{})

func (ci *CacheInstance[T]) GetInstance() *SimpleCache[T] {
	ci.once.Do(func() {
		ci.instance = &SimpleCache[T]{
			data:            make(map[string]Entry[T]),
			defaultDuration: time.Hour,
		}
	})
	return ci.instance
}

func GetCacheInstance[T any](key string) *SimpleCache[T] {
	if _, ok := cacheInstances[key]; !ok {
		cacheInstances[key] = &CacheInstance[T]{}
	}
	return cacheInstances[key].(*CacheInstance[T]).GetInstance()
}

func (sc *SimpleCache[T]) Set(key string, value T, duration time.Duration) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if duration == 0 {
		duration = sc.defaultDuration
	}

	expiration := time.Now().Add(duration)
	sc.data[key] = Entry[T]{
		Data:      value,
		ExpiresAt: expiration,
	}
}

func (sc *SimpleCache[T]) Get(key string) (*T, bool) {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	entry, found := sc.data[key]
	if !found || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return &entry.Data, true
}

func (sc *SimpleCache[T]) Has(key string) bool {
	_, found := sc.Get(key)
	return found
}

func (sc *SimpleCache[T]) Delete(key string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	delete(sc.data, key)
}

func (sc *SimpleCache[T]) Clear() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.data = make(map[string]Entry[T])
}
