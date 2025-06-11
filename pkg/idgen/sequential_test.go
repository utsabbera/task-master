package idgen

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSequential_Next(t *testing.T) {
	t.Run("should generate sequential IDs with proper padding", func(t *testing.T) {
		gen := NewSequential("TASK-", 1, 3)
		
		assert.Equal(t, "TASK-001", gen.Next())
		assert.Equal(t, "TASK-002", gen.Next())
		assert.Equal(t, "TASK-003", gen.Next())
	})
	
	t.Run("should handle different prefix and width", func(t *testing.T) {
		gen := NewSequential("USER-", 42, 4)
		
		assert.Equal(t, "USER-0042", gen.Next())
		assert.Equal(t, "USER-0043", gen.Next())
	})
	
	t.Run("should be thread-safe for concurrent access", func(t *testing.T) {
		gen := NewSequential("CONC-", 1, 3)
		numRoutines := 10
		idsPerRoutine := 100
		var wg sync.WaitGroup
		
		var mutex sync.Mutex
		allIDs := make(map[string]bool)
		
		wg.Add(numRoutines)
		for range numRoutines {
			go func() {
				defer wg.Done()
				
				localIDs := make([]string, idsPerRoutine)
				for j := 0; j < idsPerRoutine; j++ {
					localIDs[j] = gen.Next()
				}
				
				mutex.Lock()
				defer mutex.Unlock()
				for _, id := range localIDs {
					allIDs[id] = true
				}
			}()
		}
		wg.Wait()
		
		assert.Equal(t, numRoutines*idsPerRoutine, len(allIDs), "All generated IDs should be unique")
	})
}

func TestSequential_Current(t *testing.T) {
	t.Run("should return current ID without incrementing", func(t *testing.T) {
		gen := NewSequential("INV-", 42, 4)
		
		assert.Equal(t, "INV-0042", gen.Current())
		assert.Equal(t, "INV-0042", gen.Current())
		
		gen.Next()
		
		assert.Equal(t, "INV-0043", gen.Current())
	})
}

func TestSequential_Reset(t *testing.T) {
	t.Run("should reset counter to specified value", func(t *testing.T) {
		gen := NewSequential("ORD-", 1, 2)
		
		assert.Equal(t, "ORD-01", gen.Next())
		assert.Equal(t, "ORD-02", gen.Next())
		
		gen.Reset(100)
		
		assert.Equal(t, "ORD-100", gen.Current())
		assert.Equal(t, "ORD-100", gen.Next())
		assert.Equal(t, "ORD-101", gen.Next())
	})
}

func TestNewSequential(t *testing.T) {
	t.Run("should create generator with correct format", func(t *testing.T) {
		
		gen := NewSequential("TEST-", 5, 3)
		
		
		require.NotNil(t, gen)
		assert.Equal(t, "TEST-005", gen.Current())
	})
}
