package set

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert" // Assertion package
)

func TestNewBasicSet(t *testing.T) {
	s1 := NewBasicSet()
	if testing.Verbose() {
		fmt.Println(s1)
	}
	assert.Equal(t, 0, s1.Size(), "s1 should have size 0")

	s2 := NewBasicSet(1, 2, 3)
	assert.Equal(t, 3, s2.Size(), "s2 should have size 3")

	s3 := NewBasicSet(1, "2", 3.4, true, nil)
	assert.Equal(t, 5, s3.Size(), "s3 should have size 5")
}

func TestAdd(t *testing.T) {
	s1 := NewBasicSet()
	s1.Add(1)
	assert.Equal(t, 1, s1.Size(), "s1 should have size 1 after add")
	s1.Add(1)
	assert.Equal(t, 1, s1.Size(), "s1 should still have size 1 after adding same item")
	s1.Add(2)
	assert.Equal(t, 2, s1.Size(), "s1 should have size 2 after add")
}

func TestRemove(t *testing.T) {
	s1 := NewBasicSet(1, 2, 3)
	s1.Remove(2)
	assert.Equal(t, 2, s1.Size(), "s1 should have size 2 after remove")
	s1.Remove(3.0)
	assert.Equal(t, 2, s1.Size(), "s1 should have size 2 after remove not found item")
}

func TestContains(t *testing.T) {
	s1 := NewBasicSet(1, 2, 3)
	assert.True(t, s1.Contains(1), "s1 should contain 1")
	assert.False(t, s1.Contains(2.0), "s1 should not contain 2.0")
	assert.False(t, s1.Contains(nil), "s1 should not contain nil")
}

func TestClear(t *testing.T) {
	s1 := NewBasicSet(1, 2, 3)
	s1.Clear()
	assert.Equal(t, 0, s1.Size(), "s1 should have size 0 after clear")
}

func TestSlice(t *testing.T) {
	s1 := NewBasicSet(1, 2, 3)
	slice1 := s1.Slice()
	assert.Equal(t, 3, len(slice1), "slice1 should have size 3")

	s2 := NewBasicSet(1, "2")
	slice2 := s2.StringSlice()
	assert.Equal(t, 1, len(slice2), "slice2 should have size 1 (only 1 string)")
	assert.Equal(t, "2", slice2[0], "slice2 should contain 2")

	slice3 := s2.IntSlice()
	assert.Equal(t, 1, len(slice3), "slice3 should have size 1 (only 1 int)")
	assert.Equal(t, 1, slice3[0], "slice3 should contain 1")
}

func TestJSON(t *testing.T) {
	s1 := NewBasicSet(1, 2, 3, 4)
	b, err := json.Marshal(s1)
	assert.Nil(t, err, "should not error marshaling")
	if testing.Verbose() {
		fmt.Println(string(b))
	}

	s2 := NewBasicSet()
	err = json.Unmarshal(b, s2)
	assert.Nil(t, err, "should not error unmarshaling")
	assert.Equal(t, 4, s2.Size(), "s2 should have size 4")
	if testing.Verbose() {
		fmt.Println(s2.Slice())
	}

	s3 := NewOrderedSet(5, 6, 7, 8)
	b2, err := json.Marshal(s3)
	assert.Nil(t, err, "should not error marshaling")
	if testing.Verbose() {
		fmt.Println(string(b2))
	}

	s4 := NewOrderedSet()
	err = json.Unmarshal(b2, s4)
	assert.Nil(t, err, "should not error unmarshaling")
	assert.Equal(t, 4, s4.Size(), "s4 should have size 4")
	if testing.Verbose() {
		fmt.Println(s4.Slice())
	}
}

func TestJSON2(t *testing.T) {
	type TestStruct struct {
		A string    `json:"a"`
		B *BasicSet `json:"b"`
	}

	test := `{
		"a": "hello",
		"b": ["good", "bye"]
	}`

	var ts TestStruct
	err := json.Unmarshal([]byte(test), &ts)
	assert.Nil(t, err, "should not error unmarshalling")
	if testing.Verbose() {
		fmt.Println(ts.B)
	}
}

func TestOrderedSetAdd(t *testing.T) {
	s1 := NewOrderedSet()
	s1.Add(1, 2)
	if testing.Verbose() {
		fmt.Println(s1.Slice())
	}
	s1.Add(3)
	if testing.Verbose() {
		fmt.Println(s1.Slice())
	}
	assert.Equal(t, 3, s1.Size(), "s1 should have size 3")
	assert.Equal(t, []int{1, 2, 3}, s1.IntSlice(), "s1 intslice should match")
}

func TestOrderedSetRemove(t *testing.T) {
	s1 := NewOrderedSet(1, 2, 3)
	s1.Remove(2)
	assert.Equal(t, 2, s1.Size(), "s1 should have size 2")
	assert.Equal(t, []int{1, 3}, s1.IntSlice(), "s1 intslice should match")

	s1.Remove(1)
	assert.Equal(t, 1, s1.Size(), "s1 should have size 1")
	assert.Equal(t, []int{3}, s1.IntSlice(), "s1 intslice should match")

	s1.Remove(3)
	assert.Equal(t, 0, s1.Size(), "s1 should have size 0")
	assert.Equal(t, []int{}, s1.IntSlice(), "s1 intslice should match")

	s1.Add(4)
	assert.Equal(t, 1, s1.Size(), "s1 should have size 1")
	assert.Equal(t, []int{4}, s1.IntSlice(), "s1 intslice should match")
}

func TestSetUnion(t *testing.T) {
	s1 := NewBasicSet(1, 2, 3)
	s2 := NewOrderedSet(2, 3, 4)

	s3 := s2.Union(s1)
	if testing.Verbose() {
		fmt.Println(s3.Slice())
	}
	assert.Equal(t, 4, s3.Size(), "s3 should have size 4")
	assert.Equal(t, []int{2, 3, 4, 1}, s3.IntSlice(), "s3 intslice should match")

	s4 := s1.Union(s2)
	if testing.Verbose() {
		fmt.Println(s4.Slice())
	}
	assert.Equal(t, 4, s4.Size(), "s4 should have size 4")
}

func TestSetIntersection(t *testing.T) {
	s1 := NewBasicSet(1, 2, 3)
	s2 := NewOrderedSet(2, 3, 4)

	s3 := s2.Intersection(s1)
	if testing.Verbose() {
		fmt.Println(s3.Slice())
	}
	assert.Equal(t, 2, s3.Size(), "s3 should have size 2")
	assert.Equal(t, []int{2, 3}, s3.IntSlice(), "s3 intslice should match")

	s4 := s1.Intersection(s2)
	if testing.Verbose() {
		fmt.Println(s4.Slice())
	}
	assert.Equal(t, 2, s4.Size(), "s4 should have size 2")
}

func TestSetDifference(t *testing.T) {
	s1 := NewBasicSet(1, 2, 3)
	s2 := NewOrderedSet(2, 3, 4)

	s3 := s2.Difference(s1)
	if testing.Verbose() {
		fmt.Println(s3.Slice())
	}
	assert.Equal(t, 2, s3.Size(), "s3 should have size 2")
	assert.Equal(t, []int{4, 1}, s3.IntSlice(), "s3 intslice should match")

	s4 := s1.Difference(s2)
	if testing.Verbose() {
		fmt.Println(s4.Slice())
	}
	assert.Equal(t, 2, s4.Size(), "s4 should have size 2")
}
