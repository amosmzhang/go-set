package set

import (
	"encoding/json"
	"sync"
)

// Set uses map with a rwmutex for thread safety
type Set struct {
	items map[interface{}]struct{}
	lock  *sync.RWMutex
}

// basic ops
func New(items ...interface{}) *Set {
	s := &Set{
		items: make(map[interface{}]struct{}),
		lock:  &sync.RWMutex{},
	}

	s.Add(items...)
	return s
}

func (s *Set) Add(items ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range items {
		s.items[item] = struct{}{}
	}
}

func (s *Set) Remove(items ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range items {
		delete(s.items, item)
	}
}

func (s *Set) Contains(item interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, found := s.items[item]
	return found
}

func (s *Set) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.items)
}

func (s *Set) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = make(map[interface{}]struct{})
}

// get slice from set
func (s *Set) Slice() []interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	slice := make([]interface{}, 0, len(s.items))

	for item := range s.items {
		slice = append(slice, item)
	}

	return slice
}

func (s *Set) StringSlice() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	slice := make([]string, 0)
	for item := range s.items {
		if v, ok := item.(string); !ok {
			continue
		} else {
			slice = append(slice, v)
		}
	}
	return slice
}

func (s *Set) IntSlice() []int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	slice := make([]int, 0)
	for item := range s.items {
		if v, ok := item.(int); !ok {
			continue
		} else {
			slice = append(slice, v)
		}
	}
	return slice
}

// JSON ops
func (s *Set) MarshalJSON() ([]byte, error) {
	slice := s.Slice()
	return json.Marshal(slice)
}

func (s *Set) UnmarshalJSON(b []byte) error {
	var items []interface{}
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}

	s.Add(items...)
	return nil
}
