package timedmap

import (
	"sync"
	"time"
)

// [TimedMap] is a map that automatically removes entries that have expired.
// It is useful for caching data that expires after a certain period of time.
// This implementation uses a [sync.RWMutex] to synchronize access to the map and hence is thread-safe.
type TimedMap[K comparable, V any] struct {
	mu    sync.RWMutex
	i     time.Duration
	store map[K]*entry[V]
}

// New creates a new [TimedMap] with the default cleanup interval of 1 minute.
func New[K comparable, V any]() *TimedMap[K, V] {
	return NewWithCleanupInterval[K, V](time.Minute)
}

// NewWithCleanupInterval creates a new [TimedMap] with the given cleanup interval.
func NewWithCleanupInterval[K comparable, V any](interval time.Duration) *TimedMap[K, V] {
	tm := &TimedMap[K, V]{
		i:     interval,
		store: make(map[K]*entry[V]),
	}
	go tm.cleanup()
	return tm
}

// Put adds a value and its time-to-live duration to the [TimedMap] for the given key.
func (tm *TimedMap[K, V]) Put(key K, value V, ttl time.Duration) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.store[key] = &entry[V]{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Get returns the value associated with the given key and a boolean indicating if the key exists.
// If the key does not exist, it returns a zero value and false.
// If the key exists but has expired, it returns a zero value and false.
// If the key exists and has not expired, it returns the value and true.
func (tm *TimedMap[K, V]) Get(key K) (V, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	e, ok := tm.store[key]
	if !ok {
		return *new(V), false
	}
	if time.Now().After(e.expiration) {
		delete(tm.store, key)
		return *new(V), false
	}
	return e.value, true
}

// Contains returns true if the [TimedMap] contains the given key, false otherwise.
func (tm *TimedMap[K, V]) Contains(key K) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	_, ok := tm.store[key]
	return ok
}

// Delete removes the value associated with the given key regardless of its expiration time.
func (tm *TimedMap[K, V]) Delete(key K) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	delete(tm.store, key)
}

// Clear removes all entries from the [TimedMap].
func (tm *TimedMap[K, V]) Clear() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	clear(tm.store)
}

// Len returns the number of entries in the [TimedMap].
func (tm *TimedMap[K, V]) Len() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return len(tm.store)
}

type entry[V any] struct {
	value      V
	expiration time.Time
}

// cleanup removes expired entries from the [TimedMap]. It runs in a separate goroutine.
func (tm *TimedMap[K, V]) cleanup() {
	for {
		time.Sleep(tm.i)
		tm.mu.Lock()
		now := time.Now()
		for k, e := range tm.store {
			if now.After(e.expiration) {
				delete(tm.store, k)
			}
		}
		tm.mu.Unlock()
	}
}
