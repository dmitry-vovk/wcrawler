package page_fetcher

import "time"

type Option func(f *Fetcher)

// WithTimeout sets default timeouts for network operations
func WithTimeout(timeout time.Duration) Option {
	return func(f *Fetcher) {
		f.timeout = timeout
	}
}

// WithAccept sets the value of 'Accept' http header to use
func WithAccept(acceptHeader string) Option {
	return func(f *Fetcher) {
		f.accept = acceptHeader
	}
}

// WithUserAgent sets User-Agent header value
func WithUserAgent(userAgent string) Option {
	return func(f *Fetcher) {
		f.userAgent = userAgent
	}
}

// WithHeadRequests sets the flag to issue HEAD before GET
func WithHeadRequests(doHeadRequests bool) Option {
	return func(f *Fetcher) {
		f.doHeadRequests = doHeadRequests
	}
}
