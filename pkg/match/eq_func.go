package match

import (
	"go.uber.org/mock/gomock"
)

// EqFunc returns a matcher that matches if the provided function returns true for the argument.
//
// Example:
//
//	mockRepo.EXPECT().Save(match.EqFunc(func(t *Task) bool {
//	    return t.ID == "123" && t.Title == "Example"
//	})).Return(nil)
//
// This will match any argument for which the function returns true.
func EqFunc[T any](f func(T) bool) gomock.Matcher {
	return &eqFuncMatcher[T]{f: f}
}

type eqFuncMatcher[T any] struct {
	f func(T) bool
}

func (e *eqFuncMatcher[T]) Matches(x any) bool {
	if e.f == nil {
		return false
	}
	v, ok := x.(T)
	if !ok {
		return false
	}
	return e.f(v)
}

func (e *eqFuncMatcher[T]) String() string {
	return "matches custom function"
}
