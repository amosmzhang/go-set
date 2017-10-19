package set

import (
	"encoding/json"
	"sync"
)

type element struct {
	val  interface{}
	prev *element
	next *element
}

// OrderedSet uses map with a rwmutex for thread safety
// Element insertion order is kept
type OrderedSet struct {
	items map[interface{}]*element
	root  *element
	lock  *sync.RWMutex
	init  bool
}

// basic ops
func NewOrderedSet(items ...interface{}) *OrderedSet {
	s := &OrderedSet{
		items: make(map[interface{}]*element),
		root:  &element{nil, nil, nil},
		lock:  &sync.RWMutex{},
		init:  true,
	}

	s.Add(items...)
	return s
}

func (s *OrderedSet) Add(items ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range items {
		if _, found := s.items[item]; !found {
			last := s.root.prev
			if last == nil {
				last = s.root
			}
			adding := &element{
				val:  item,
				prev: last,
				next: s.root,
			}
			last.next = adding
			s.root.prev = adding
			s.items[item] = adding
		}
	}
}

func (s *OrderedSet) Remove(items ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range items {
		if deleting, found := s.items[item]; found {
			prev := deleting.prev
			next := deleting.next

			deleting.prev.next = next
			deleting.next.prev = prev
			delete(s.items, item)
		}
	}
}

func (s *OrderedSet) Contains(item interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, found := s.items[item]
	return found
}

func (s *OrderedSet) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.items)
}

func (s *OrderedSet) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = make(map[interface{}]*element)
	s.root = &element{nil, nil, nil}
}

// Set operations
func (s *OrderedSet) IsSubsetOf(other Set) bool {
	for _, item := range s.items {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

func (s *OrderedSet) IsSupersetOf(other Set) bool {
	return other.IsSubsetOf(s)
}

func (s *OrderedSet) Union(other Set) Set {
	res := NewOrderedSet(s.Slice()...)
	res.Add(other.Slice()...)
	return res
}

func (s *OrderedSet) Intersection(other Set) Set {
	res := NewOrderedSet(s.Slice()...)
	for _, item := range s.items {
		if !other.Contains(item.val) {
			res.Remove(item.val)
		}
	}
	return res
}

func (s *OrderedSet) Difference(other Set) Set {
	res := s.Union(other)
	res.Remove(s.Intersection(other).Slice()...)
	return res
}

// get slice from set
func (s *OrderedSet) Slice() []interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	slice := make([]interface{}, 0, len(s.items))
	for i := s.root; i != nil && i != i.next && i.next != s.root; {
		i = i.next
		slice = append(slice, i.val)
	}

	return slice
}

func (s *OrderedSet) StringSlice() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	slice := make([]string, 0)
	for i := s.root; i != nil && i != i.next && i.next != s.root; {
		i = i.next
		if v, ok := i.val.(string); !ok {
			continue
		} else {
			slice = append(slice, v)
		}
	}
	return slice
}

func (s *OrderedSet) IntSlice() []int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	slice := make([]int, 0)
	for i := s.root; i != nil && i != i.next && i.next != s.root; {
		i = i.next
		if v, ok := i.val.(int); !ok {
			continue
		} else {
			slice = append(slice, v)
		}
	}
	return slice
}

// JSON ops
func (s *OrderedSet) MarshalJSON() ([]byte, error) {
	slice := s.Slice()
	return json.Marshal(slice)
}

func (s *OrderedSet) UnmarshalJSON(b []byte) error {
	if !s.init {
		s.items = make(map[interface{}]*element)
		s.root = &element{nil, nil, nil}
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
