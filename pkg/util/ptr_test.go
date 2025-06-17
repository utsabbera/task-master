package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	type testCase[T any] struct {
		name  string
		input T
	}
	t.Run("int", func(t *testing.T) {
		tests := []testCase[int]{
			{name: "should return pointer to int value", input: 42},
			{name: "should return pointer to zero int", input: 0},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				ptr := Ptr(tc.input)
				assert.NotNil(t, ptr)
				assert.Equal(t, tc.input, *ptr)
			})
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []testCase[string]{
			{name: "should return pointer to non-empty string", input: "hello"},
			{name: "should return pointer to empty string", input: ""},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				ptr := Ptr(tc.input)
				assert.NotNil(t, ptr)
				assert.Equal(t, tc.input, *ptr)
			})
		}
	})

	t.Run("struct", func(t *testing.T) {
		type sample struct {
			A int
			B string
		}
		val := sample{A: 1, B: "test"}
		t.Run("should return pointer to struct value", func(t *testing.T) {
			ptr := Ptr(val)
			assert.NotNil(t, ptr)
			assert.Equal(t, val, *ptr)
		})
	})
}
