// Package match provides custom matchers for testing with gomock and similar frameworks.
package match

import (
	"fmt"
	"reflect"

	"go.uber.org/mock/gomock"
)

// PtrTo returns a matcher that matches if a pointer argument points to a value
// equal to the expected value. It performs a deep equality check on the values
// rather than comparing pointer addresses. This is useful when you want to verify
// that a function was called with a pointer to data with specific content, but
// you don't care about the exact pointer instance.
//
// For example:
//
//	expectedTask := &Task{ID: "123", Title: "Example"}
//	mockRepo.EXPECT().Save(match.PtrTo(expectedTask)).Return(nil)
//
//	// Even if a different pointer is passed, it will match if content is equal
//	actualTask := &Task{ID: "123", Title: "Example"}
//	mockRepo.Save(actualTask) // This will match the expectation
func PtrTo(expected any) gomock.Matcher {
	return &ptrMatcher{expected}
}

// ptrMatcher is a gomock.Matcher that matches pointer values by comparing
// their dereferenced values rather than the pointer addresses.
type ptrMatcher struct {
	expected any
}

// Matches returns true if x is a pointer to a value that is equal to the expected value.
// It handles nil values and properly compares both pointer and non-pointer expected values.
func (p *ptrMatcher) Matches(x any) bool {
	if p.expected == nil {
		if x == nil {
			return true
		}

		xVal := reflect.ValueOf(x)
		if xVal.Kind() == reflect.Ptr && xVal.IsNil() {
			return true
		}
		return false
	}

	if x == nil {
		return false
	}

	xVal := reflect.ValueOf(x)
	expectedVal := reflect.ValueOf(p.expected)

	if xVal.Kind() != reflect.Ptr {
		return false
	}

	if xVal.IsNil() {
		return expectedVal.Kind() == reflect.Ptr && expectedVal.IsNil()
	}

	xElem := xVal.Elem().Interface()

	if expectedVal.Kind() != reflect.Ptr {
		return reflect.DeepEqual(xElem, p.expected)
	}

	if expectedVal.IsNil() {
		return false
	}

	expectedElem := expectedVal.Elem().Interface()
	return reflect.DeepEqual(xElem, expectedElem)
}

// String returns a description of the matcher.
func (p *ptrMatcher) String() string {
	return fmt.Sprintf("points to %v", p.expected)
}
