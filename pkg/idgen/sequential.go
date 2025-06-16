package idgen

import (
	"fmt"
	"sync"
)

// Sequential creates sequential IDs with a given prefix
// It implements the Generator interface for generating predictable ID sequences
type Sequential struct {
	// prefix is the string that appears at the beginning of each generated ID
	prefix string
	// nextID is the numeric portion that will be used for the next ID
	nextID int
	// idFormat is the format string used to create IDs with consistent width
	idFormat string
	// mu protects concurrent access to the generator state
	mu sync.Mutex
}

// NewSequential creates a new sequential ID generator with the given prefix
// and starting from the specified ID. The width parameter determines the minimum
// width of the numeric portion of the ID, padded with zeros if necessary.
func NewSequential(prefix string, startID, width int) *Sequential {
	return &Sequential{
		prefix:   prefix,
		nextID:   startID,
		idFormat: "%s%0" + fmt.Sprintf("%dd", width),
		mu:       sync.Mutex{},
	}
}

// Next generates the next sequential ID
func (g *Sequential) Next() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	id := fmt.Sprintf(g.idFormat, g.prefix, g.nextID)
	g.nextID++
	return id
}

// Current returns the current ID without advancing the counter
func (g *Sequential) Current() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	return fmt.Sprintf(g.idFormat, g.prefix, g.nextID)
}

// Reset sets the next ID to the specified value
func (g *Sequential) Reset(nextID int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.nextID = nextID
}
