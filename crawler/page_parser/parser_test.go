package page_parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	r := bytes.NewReader([]byte(html))
	if result, err := Parse(r); assert.NoError(t, err) {
		assert.Equal(t, "http://example.com/foo/bar", result.CanonicalURL)
		assert.Equal(t, "http://example.com/foo/bar/", result.BaseURL)
		assert.Equal(t, []string{
			"http://example.com/",
			"http://example.com/foo/bar/page.html",
			"http://example.com/foo/",
			"http://some.other.com",
			"javascript:void(0)",
		}, result.Links)
	}
}

// language=HTML
const html = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
	<base href="http://example.com/foo/bar/">
	<link rel="canonical" href="http://example.com/foo/bar ">
    <title>Foo Bar</title>
</head>
<body>
	<a href="/">Home</a>
	<a href="page.html">Some page</a>
	<a href="..">Up</a>
	<a href="http://some.other.com">External link</a>
	<a href=" javascript:void(0)">Click here!</a>
</body>
</html>
`
