package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqFunc(t *testing.T) {
	matcher := EqFunc(func(x int) bool {
		return x%2 == 0
	})

	tests := []struct {
		name  string
		input any
		want  bool
	}{
		{"even int", 4, true},
		{"odd int", 3, false},
		{"not int", "foo", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, matcher.Matches(tt.input))
		})
	}
}
