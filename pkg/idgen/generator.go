package idgen

// Generator defines an interface for ID generators
type Generator interface {
	// Next generates the next ID in sequence
	Next() string
	// Current returns the current ID without advancing
	Current() string
	// Reset sets the next ID to the specified value
	Reset(nextID int)
}
