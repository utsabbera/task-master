package util

import time "time"

//go:generate mockgen -destination=clock_mock.go -package=util . Clock

// Clock is an interface that provides the current time.
// It allows for abstraction over time retrieval, which is useful for testing and mocking.
type Clock interface {
	Now() time.Time
}

// NewClock creates a new instance of the Clock interface.
// It returns a realClock, which uses the system's current time.
func NewClock() Clock {
	return realClock{}
}

type realClock struct{}

// Now returns the current local time.
// This method satisfies the Clock interface.
func (realClock) Now() time.Time {
	return time.Now()
}
