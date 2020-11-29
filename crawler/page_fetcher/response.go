package page_fetcher

import (
	"io"
	"net/http"
	"net/url"
)

// Response provides only necessary http response values
type Response struct {
	// Requested URL of the page
	URL *url.URL
	// Response status code
	StatusCode int
	// Response headers
	Headers http.Header
	// Page contents
	Body io.ReadCloser
}
