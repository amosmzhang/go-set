package set

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert" // Assertion package
)

func TestNew(t *testing.T) {
	s1 := New()
	fmt.Println(s1)
	assert.Equal(t, 0, s1.Size(), "s1 should have size 0")

	s2 := New(1, 2, 3)
	assert.Equal(t, 3, s2.Size(), "s2 should have size 3")

	s3 := New(1, "2", 3.4, true, nil)
	assert.Equal(t, 5, s3.Size(), "s3 should have size 5")
}

func TestAdd(t *testing.T) {
	s1 := New()
	s1.Add(1)
	assert.Equal(t, 1, s1.Size(), "s1 should have size 1 after add")
	s1.Add(1)
	assert.Equal(t, 1, s1.Size(), "s1 should still have size 1 after adding same item")
	s1.Add(2)
	assert.Equal(t, 2, s1.Size(), "s1 should have size 2 after add")
}

func TestRemove(t *testing.T) {
	s1 := New(1, 2, 3)
	s1.Remove(2)
	assert.Equal(t, 2, s1.Size(), "s1 should have size 2 after remove")
	s1.Remove(3.0)
	assert.Equal(t, 2, s1.Size(), "s1 should have size 2 after remove not found item")
}

func TestContains(t *testing.T) {
	s1 := New(1, 2, 3)
	assert.True(t, s1.Contains(1), "s1 should contain 1")
	assert.False(t, s1.Contains(2.0), "s1 should not contain 2.0")
	assert.False(t, s1.Contains(nil), "s1 should not contain nil")
}

func TestClear(t *testing.T) {
	s1 := New(1, 2, 3)
	s1.Clear()
	assert.Equal(t, 0, s1.Size(), "s1 should have size 0 after clear")
}

func TestSlice(t *testing.T) {
	s1 := New(1, 2, 3)
	slice1 := s1.Slice()
	assert.Equal(t, 3, len(slice1), "slice1 should have size 3")

	s2 := New(1, "2")
	slice2 := s2.StringSlice()
	assert.Equal(t, 1, len(slice2), "slice2 should have size 1 (only 1 string)")
	assert.Equal(t, "2", slice2[0], "slice2 should contain 2")

	slice3 := s2.IntSlice()
	assert.Equal(t, 1, len(slice3), "slice3 should have size 1 (only 1 int)")
	assert.Equal(t, 1, slice3[0], "slice3 should contain 1")
}

func TestJSON(t *testing.T) {
	s1 := New(1, 2, 3, 4)
	b, err := json.Marshal(s1)
	assert.Nil(t, err, "should not error marshaling")
	fmt.Println(string(b))

	s2 := New()
	err = json.Unmarshal(b, s2)
	assert.Nil(t, err, "should not error unmarshaling")
	assert.Equal(t, 4, s2.Size(), "s2 should have size 4")
	fmt.Println(s2.Slice())
}

func TestJSON2(t *testing.T) {
	type TestStruct struct {
		A string `json:"a"`
		B *Set   `json:"b"`
	}

	test := `{
		"a": "hello",
		"b": ["good", "bye"]
	}`

	var ts TestStruct
	err := json.Unmarshal([]byte(test), &ts)
	assert.Nil(t, err, "should not error unmarshalling")
	fmt.Println(ts.B)
}
