package page_fetcher

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Fetcher can go and fetch remote web page
type Fetcher struct {
	timeout        time.Duration // network operations timeout value
	accept         string        // default 'Accept' http header
	userAgent      string        // default 'User-Agent' http header
	doHeadRequests bool          // whether to perform HEAD requests before GET requests
	client         *http.Client  // http client to use for requests
}

type method string

const (
	// Timeout value for network operations
	defaultTimeout = time.Second * 10
	// Tell servers we are interested in this content types
	defaultAcceptHeader        = `text/html,application/xhtml+xml`
	methodGET           method = `GET`
	methodHEAD          method = `HEAD`
)

// NewFetcher creates an instance of Fetcher with options
func NewFetcher(options ...Option) *Fetcher {
	f := Fetcher{}
	// apply options
	for _, fn := range options {
		fn(&f)
	}
	// set defaults
	if f.timeout == 0 {
		f.timeout = defaultTimeout
	}
	if f.accept == "" {
		f.accept = defaultAcceptHeader
	}
	// build http client
	f.client = &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout:   f.timeout,
			ResponseHeaderTimeout: f.timeout,
			ExpectContinueTimeout: f.timeout,
		},
		Timeout: f.timeout,
	}
	// cookiejar.New() without options does not return an error
	f.client.Jar, _ = cookiejar.New(nil)
	return &f
}

// Fetch performs http requests and build response object
func (f *Fetcher) Fetch(r *Request) (*Response, error) {
	if f.doHeadRequests {
		resp, err := f.client.Do(f.buildRequest(r, methodHEAD))
		if err != nil {
			// Error on HEAD request is not critical, let's do GET anyway
			log.Printf("HEAD request error: %s", err)
		} else if !r.acceptableResponse(resp) {
			return nil, ErrBadContentType
		}
	}
	resp, err := f.client.Do(f.buildRequest(r, methodGET))
	if err != nil {
		return nil, err
	}
	if !r.acceptableResponse(resp) {
		return nil, ErrBadContentType
	}
	return buildResponse(r, resp), nil
}

// buildRequest assembles http.Request according to parameters
func (f Fetcher) buildRequest(r *Request, method method) *http.Request {
	link := r.URL.String()
	// http.NewRequest will not return an error with this set of arguments
	httpRequest, _ := http.NewRequest(string(method), link, nil)
	if f.userAgent != "" {
		httpRequest.Header.Add("User-Agent", f.userAgent)
	}
	httpRequest.Header.Add("Referer", r.HTTPReferrer)
	httpRequest.Header.Add("Accept", f.accept)
	return httpRequest
}

func buildResponse(req *Request, resp *http.Response) *Response {
	return &Response{
		URL:        req.URL,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       resp.Body,
	}
}
