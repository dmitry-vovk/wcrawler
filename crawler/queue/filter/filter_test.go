package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	testCases := []struct {
		link     string
		expected string
		ok       bool
	}{
		{
			link:     "http://example.com",
			expected: "http://example.com",
			ok:       true,
		},
		{
			link:     "http://example.com",
			expected: "",
			ok:       false,
		},
		{
			link:     "http://example.com/search#id",
			expected: "http://example.com/search",
			ok:       true,
		},
	}
	for _, tt := range testCases {
		normal, ok := Filter(tt.link)
		if assert.Equal(t, tt.ok, ok) && ok {
			assert.Equal(t, tt.expected, normal)
		}
	}
}
