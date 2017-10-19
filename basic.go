package set

import (
	"encoding/json"
	"sync"
)

// BasicSet uses map with a rwmutex for thread safety, is unordered
type BasicSet struct {
	items map[interface{}]struct{}
	lock  *sync.RWMutex
	init  bool
}

func NewBasicSet(items ...interface{}) *BasicSet {
	s := &BasicSet{
		items: make(map[interface{}]struct{}),
		lock:  &sync.RWMutex{},
		init:  true,
	}

	s.Add(items...)
	return s
}

// basic ops
func (s *BasicSet) Add(items ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range items {
		s.items[item] = struct{}{}
	}
}

func (s *BasicSet) Remove(items ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range items {
		delete(s.items, item)
	}
}

func (s *BasicSet) Contains(item interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, found := s.items[item]
	return found
}

func (s *BasicSet) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.items)
}

func (s *BasicSet) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = make(map[interface{}]struct{})
}

// Set operations
func (s *BasicSet) IsSubsetOf(other Set) bool {
	for item := range s.items {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

func (s *BasicSet) IsSupersetOf(other Set) bool {
	return other.IsSubsetOf(s)
}

func (s *BasicSet) Union(other Set) Set {
	res := NewBasicSet(s.Slice()...)
	res.Add(other.Slice()...)
	return res
}

func (s *BasicSet) Intersection(other Set) Set {
	res := NewBasicSet(s.Slice()...)
	for item := range s.items {
		if !other.Contains(item) {
			res.Remove(item)
		}
	}
	return res
}

func (s *BasicSet) Difference(other Set) Set {
	res := s.Union(other)
	res.Remove(s.Intersection(other).Slice()...)
	return res
}

// get slice from set
func (s *BasicSet) Slice() []interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	slice := make([]interface{}, 0, len(s.items))
	for item := range s.items {
		slice = append(slice, item)
	}

	return slice
}

func (s *BasicSet) StringSlice() []string {
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

func (s *BasicSet) IntSlice() []int {
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
func (s *BasicSet) MarshalJSON() ([]byte, error) {
	slice := s.Slice()
	return json.Marshal(slice)
}

func (s *BasicSet) UnmarshalJSON(b []byte) error {
	if !s.init {
		s.items = make(map[interface{}]struct{})
		s.lock = &sync.RWMutex{}
		s.init = true
	}

	var items []interface{}
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}

	s.Add(items...)
	return nil
}
