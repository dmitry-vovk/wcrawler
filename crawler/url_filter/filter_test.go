package url_filter

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
			expected: "http://example.com/",
			ok:       true,
		},
		{
			link:     "http://example.com/",
			expected: "http://example.com/",
			ok:       true,
		},
		{
			link: "http://www.example.com",
			ok:   false,
		},
		{
			link:     "http://example.com/search#id",
			expected: "http://example.com/search",
			ok:       true,
		},
		{
			link:     "https://example.com/about",
			expected: "https://example.com/about",
			ok:       true,
		},
		{
			link: "http://another.domain.tld",
			ok:   false,
		},
		{
			link: "https://example.com/fail",
			ok:   false,
		},
	}
	f := NewFilter("example.com").WithRobots(&testRobot{}, "")
	for _, tt := range testCases {
		normal, ok := f.Filter(tt.link)
		if assert.Equal(t, tt.ok, ok, tt.link) && ok {
			assert.Equal(t, tt.expected, normal, tt.link)
		}
	}
	{
		f.AllowWWWPrefix(true)
		_, ok := f.Filter("http://www.example.com")
		assert.True(t, ok)
	}
	{
		f.WithRobots(&testRobot{}, "Failer")
		_, ok := f.Filter("http://example.com")
		assert.False(t, ok)
	}
	{
		_, ok := f.Filter("😀")
		assert.False(t, ok)
	}
	{ // invalid URL character
		_, ok := f.Filter(string(rune(0x7f)))
		assert.False(t, ok)
	}
}

type testRobot struct{}

func (t *testRobot) TestAgent(path, agent string) bool {
	return !(path == "/fail" || agent == "Failer")
}
