package page_fetcher

import (
	"io"
	"net/http"
	"net/url"
)

type Response struct {
	// Requested URL of the page
	OriginalURL *url.URL
	// Actual URL of the page
	ActualURL *url.URL
	// Response status code
	StatusCode int
	// Response headers
	Headers http.Header
	// Page contents
	Body io.ReadCloser
}
