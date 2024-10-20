package timedmap

import (
	"sync"
	"testing"
	"time"
)

func TestTimedMapBasicCRUD(t *testing.T) {
	tm := New[string, int](time.Minute)
	tm.Put("key", 19, time.Second)
	value, ok := tm.Get("key")
	if !ok || value != 19 {
		t.Errorf("expected value 19, got %d", value)
	}
	tm.Delete("key")
	_, ok = tm.Get("key")
	if ok {
		t.Errorf("expected key to be deleted")
	}
}

func TestTimedMapGetNonExistentKey(t *testing.T) {
	tm := New[string, int](time.Minute)
	_, ok := tm.Get("non-existent-key")
	if ok {
		t.Errorf("expected ok to be false")
	}
}

func TestTimedMapGetExpiredKey(t *testing.T) {
	tm := New[string, int](time.Minute)
	tm.Put("key", 19, 100*time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	_, ok := tm.Get("key")
	if ok {
		t.Errorf("expected ok to be false")
	}
}

func TestTimedMapPutSameKeyMultipleTimes(t *testing.T) {
	tm := New[string, int](time.Minute)
	tm.Put("key", 19, time.Second)
	tm.Put("key", 23, time.Second)
	value, _ := tm.Get("key")
	if value != 23 {
		t.Errorf("expected value 23, got %d", value)
	}
}

func TestTimedMapDeleteNonExistentKey(t *testing.T) {
	tm := New[string, int](time.Minute)
	tm.Delete("non-existent-key")
	if tm.Size() != 0 {
		t.Errorf("expected size 0, got %d", tm.Size())
	}
}

func TestTimedMapSize(t *testing.T) {
	tm := New[string, int](time.Minute)
	if tm.Size() != 0 {
		t.Errorf("expected size 0, got %d", tm.Size())
	}
	tm.Put("key", 19, time.Second)
	if tm.Size() != 1 {
		t.Errorf("expected size 1, got %d", tm.Size())
	}
	tm.Delete("key")
	if tm.Size() != 0 {
		t.Errorf("expected size 0, got %d", tm.Size())
	}
}

func TestTimedMapContains(t *testing.T) {
	tm := New[string, int](time.Minute)
	tm.Put("key", 19, time.Second)
	if !tm.Contains("key") {
		t.Errorf("expected key to be present")
	}
	tm.Delete("key")
	if tm.Contains("key") {
		t.Errorf("expected key to be removed")
	}
}

func TestTimedMapClear(t *testing.T) {
	tm := New[string, int](time.Minute)
	tm.Put("key1", 19, time.Second)
	tm.Put("key2", 23, time.Second)
	tm.Clear()
	if tm.Size() != 0 {
		t.Errorf("expected size 0, got %d", tm.Size())
	}
}

func TestTimedMapExpiration(t *testing.T) {
	tm := New[string, int](time.Minute)
	tm.Put("key1", 19, 300*time.Millisecond)
	tm.Put("key2", 23, 100*time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	_, ok := tm.Get("key1")
	if !ok {
		t.Errorf("expected key1 to still be present")
	}
	_, ok = tm.Get("key2")
	if ok {
		t.Errorf("expected key2 to be expired and removed")
	}
}

func TestTimedMapCleanup(t *testing.T) {
	tm := New[string, int](200 * time.Millisecond)
	tm.Put("key", 19, 100*time.Millisecond)
	time.Sleep(300 * time.Millisecond)
	_, ok := tm.Get("key")
	if ok {
		t.Errorf("expected key to be cleaned up and removed")
	}
}

func TestTimedMapConcurrency(t *testing.T) {
	m := New[string, int](time.Minute)
	var wg sync.WaitGroup
	// Launch multiple goroutines to simulate concurrent access
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Alternate between put and get operations
			if i%2 == 0 {
				m.Put("key", i, time.Second)
			} else {
				_, _ = m.Get("key")
			}
		}(i)
	}
	wg.Wait()
	if value, ok := m.Get("key"); !ok || value == 0 {
		t.Errorf("expected value to exist for key, but it was missing")
	}
}
