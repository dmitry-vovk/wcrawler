package page_fetcher

import (
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	// URL to visit
	URL *url.URL
	// HTTP Referrer header value
	HTTPReferrer string
	// Valid content types
	AcceptableContentTypes map[string]struct{}
}

func (r *Request) unacceptablePage(resp *http.Response) bool {
	if r.AcceptableContentTypes == nil {
		return true
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return false
	}
	for expectedContentType := range r.AcceptableContentTypes {
		if strings.Contains(contentType, expectedContentType) {
			return true
		}
	}
	return false
}
