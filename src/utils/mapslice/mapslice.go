// Package mapslice
// Created by RTT.
// Author: teocci@yandex.com on 2023-8월-18
package mapslice

import (
	"fmt"
	"sort"
	"strings"
)

type Entry[T any] struct {
	Key   string
	Value T
}

type SortedMap[T any] struct {
	entries []Entry[T]
}

func New[T any](m map[string]T) SortedMap[T] {
	entries := make([]Entry[T], 0, len(m))
	for key, value := range m {
		entries = append(entries, Entry[T]{Key: key, Value: value})
	}

	return SortedMap[T]{entries: entries}
}

func (sm SortedMap[T]) SortBy(comparator func(a, b Entry[T]) bool) {
	sort.SliceStable(sm.entries, func(i, j int) bool {
		return comparator(sm.entries[i], sm.entries[j])
	})
}

func (sm SortedMap[T]) Keys() []string {
	keys := make([]string, len(sm.entries))
	for i, entry := range sm.entries {
		keys[i] = entry.Key
	}
	return keys
}

func (sm SortedMap[T]) Values() []T {
	values := make([]T, len(sm.entries))
	for i, entry := range sm.entries {
		values[i] = entry.Value
	}
	return values
}

func (sm SortedMap[T]) Entries() []Entry[T] {
	return sm.entries
}

func (sm SortedMap[T]) Len() int {
	return len(sm.entries)
}

func (sm SortedMap[T]) Less(i, j int) bool {
	return sm.entries[i].Key < sm.entries[j].Key
}

func (sm SortedMap[T]) Swap(i, j int) {
	sm.entries[i], sm.entries[j] = sm.entries[j], sm.entries[i]
}

func (sm SortedMap[T]) Great(i, j int) bool {
	return sm.entries[i].Key > sm.entries[j].Key
}

func (sm SortedMap[T]) Equal(i, j int) bool {
	return sm.entries[i].Key == sm.entries[j].Key
}

func (sm SortedMap[T]) PrintKeys() {
	keys := sm.Keys()
	fmt.Printf("Keys: [%s]", strings.Join(keys, ", "))
}

func (sm SortedMap[T]) PrintValues() {
	values := sm.Values()
	fmt.Printf("Values: {%v}", values)
}

func (sm SortedMap[T]) PrintEntries() {
	entries := sm.Entries()
	fmt.Printf("Entries: {%v}", entries)
}
