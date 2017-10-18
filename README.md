# go-set

```go
type Set interface {
	Add(items ...interface{})
	Remove(items ...interface{})
	Contains(item interface{}) bool
	Size() int
	Clear()

	IsSubsetOf(other Set) bool
	IsSupersetOf(other Set) bool
	Union(other Set) Set
	Intersection(other Set) Set
	Difference(other Set) Set

	Slice() []interface{}
	StringSlice() []string
	IntSlice() []int

	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
}
```

Implemented in `BasicSet` and `OrderedSet`. See tests for sample usage.
