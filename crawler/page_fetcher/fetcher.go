package page_fetcher

import (
	"errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const defaultTimeout = time.Second * 10

var (
	// HTTP client that will be used for requests
	client = &http.Client{
		Transport: &http.Transport{
			// Here an below, setting sensible timeouts helps to avoid blocking
			TLSHandshakeTimeout:   defaultTimeout,
			ResponseHeaderTimeout: defaultTimeout,
			ExpectContinueTimeout: defaultTimeout,
		},
		Timeout: defaultTimeout,
	}
)

func init() {
	// cookiejar.New without options does not return an error
	client.Jar, _ = cookiejar.New(nil)
}

func Fetch(r *Request) (*Response, error) {
	log.Printf("Fetching %s", r.URL)
	link := r.URL.String()
	if r.DoHeadRequest {
		httpRequest, _ := http.NewRequest("HEAD", link, nil)
		httpRequest.Header.Add("User-Agent", r.UserAgent)
		httpRequest.Header.Add("Referer", r.HTTPReferrer)
		resp, err := client.Do(httpRequest)
		if err == nil {
			if !r.unacceptablePage(resp) {
				return nil, errors.New("unacceptable page type")
			}
		} else {
			log.Printf("HEAD request error: %s", err)
		}
	}
	httpRequest, _ := http.NewRequest("GET", link, nil)
	httpRequest.Header.Add("User-Agent", r.UserAgent)
	httpRequest.Header.Add("Referer", r.HTTPReferrer)
	resp, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	return buildResponse(r, resp), nil
}

func buildResponse(req *Request, resp *http.Response) *Response {
	return &Response{
		Error:       nil,
		OriginalURL: req.URL,
		ActualURL:   resp.Request.URL,
		StatusCode:  resp.StatusCode,
		Headers:     resp.Header,
		Body:        resp.Body,
	}
}
