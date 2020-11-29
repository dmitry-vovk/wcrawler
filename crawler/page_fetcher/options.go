package page_fetcher

import "time"

type Option func(f *Fetcher)

func WithTimeout(timeout time.Duration) Option {
	return func(f *Fetcher) {
		f.timeout = timeout
	}
}

func WithAccept(acceptHeader string) Option {
	return func(f *Fetcher) {
		f.accept = acceptHeader
	}
}

func WithUserAgent(userAgent string) Option {
	return func(f *Fetcher) {
		f.userAgent = userAgent
	}
}

func WithHeadRequests(doHeadRequests bool) Option {
	return func(f *Fetcher) {
		f.doHeadRequests = doHeadRequests
	}
}
