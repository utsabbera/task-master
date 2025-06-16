package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStruct is a sample struct used for testing
type TestStruct struct {
	ID   string
	Name string
	Age  int
}

func TestPtrMatcherMatches(t *testing.T) {
	t.Run("should match pointers with equal values", func(t *testing.T) {
		val1 := &TestStruct{ID: "1", Name: "Test", Age: 25}
		val2 := &TestStruct{ID: "1", Name: "Test", Age: 25}

		matcher := ptrMatcher{val1}
		assert.True(t, matcher.Matches(val2))
	})

	t.Run("should not match pointers with different values", func(t *testing.T) {
		val1 := &TestStruct{ID: "1", Name: "Test", Age: 25}
		val2 := &TestStruct{ID: "1", Name: "Test", Age: 30}

		matcher := ptrMatcher{val1}
		assert.False(t, matcher.Matches(val2))
	})

	t.Run("should match against non-pointer expected value", func(t *testing.T) {
		val1 := TestStruct{ID: "1", Name: "Test", Age: 25}
		val2 := &TestStruct{ID: "1", Name: "Test", Age: 25}

		matcher := ptrMatcher{val1}
		assert.True(t, matcher.Matches(val2))
	})

	t.Run("should reject non-pointer actual value", func(t *testing.T) {
		val1 := &TestStruct{ID: "1", Name: "Test", Age: 25}
		val2 := TestStruct{ID: "1", Name: "Test", Age: 25}

		matcher := ptrMatcher{val1}
		assert.False(t, matcher.Matches(val2))
	})

	t.Run("should handle nil values correctly", func(t *testing.T) {
		nilPtrMatcher := ptrMatcher{nil}
		assert.True(t, nilPtrMatcher.Matches(nil))

		var nilPtr *TestStruct = nil
		assert.True(t, (nilPtrMatcher).Matches(nilPtr))

		var anotherNilPtr *TestStruct = nil
		nilStructPtrMatcher := ptrMatcher{nilPtr}
		assert.True(t, nilStructPtrMatcher.Matches(anotherNilPtr))

		nonNilPtr := &TestStruct{ID: "1"}
		assert.False(t, nilPtrMatcher.Matches(nonNilPtr))

		nonNilPtrMatcher := ptrMatcher{nonNilPtr}
		assert.False(t, nonNilPtrMatcher.Matches(nil))
	})
}

func TestPtrMatcherString(t *testing.T) {
	t.Run("should format non-nil pointer expected value correctly", func(t *testing.T) {
		val := &TestStruct{ID: "1", Name: "Test", Age: 25}
		matcher := ptrMatcher{val}
		expected := "points to &{1 Test 25}"
		assert.Equal(t, expected, matcher.String())
	})

	t.Run("should format non-pointer expected value correctly", func(t *testing.T) {
		val := TestStruct{ID: "1", Name: "Test", Age: 25}
		matcher := ptrMatcher{val}
		expected := "points to {1 Test 25}"
		assert.Equal(t, expected, matcher.String())
	})

	t.Run("should format nil expected value correctly", func(t *testing.T) {
		matcher := ptrMatcher{nil}
		expected := "points to <nil>"
		assert.Equal(t, expected, matcher.String())
	})
}

func TestPtrTo(t *testing.T) {
	t.Run("should create matcher for non-nil pointer", func(t *testing.T) {
		val := &TestStruct{ID: "1", Name: "Test", Age: 25}
		matcher := PtrTo(val)
		assert.Equal(t, &ptrMatcher{val}, matcher)
	})

	t.Run("should create matcher for nil pointer", func(t *testing.T) {
		matcher := PtrTo(nil)
		assert.Equal(t, &ptrMatcher{nil}, matcher)
	})

	t.Run("should create matcher for non-pointer value", func(t *testing.T) {
		val := TestStruct{ID: "1", Name: "Test", Age: 25}
		matcher := PtrTo(val)
		assert.Equal(t, &ptrMatcher{val}, matcher)
	})
}
