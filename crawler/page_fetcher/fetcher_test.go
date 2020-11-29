package page_fetcher

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetch_Head(t *testing.T) {
	s := startServer()
	f := NewFetcher(
		WithTimeout(time.Second),
		WithUserAgent("Bot/1"),
		WithAccept("application/binary"),
		WithHeadRequests(true),
	)
	r := &Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   s.listener.Addr().String(),
		},
		AcceptableContentTypes: map[string]struct{}{
			"text/html": {},
		},
		HTTPReferrer: "/referrer",
	}
	s.responseCode = 202
	if resp, err := f.Fetch(r); assert.NoError(t, err) {
		if body, err := ioutil.ReadAll(resp.Body); assert.NoError(t, err) {
			assert.Equal(t, []byte(testHTML), body)
		}
		_ = resp.Body.Close()
		assert.Equal(t, []string{"HEAD", "GET"}, s.methods())
		assert.Equal(t, []string{"Bot/1", "Bot/1"}, s.userAgents())
		assert.Equal(t, []string{"application/binary", "application/binary"}, s.accepts())
		assert.Equal(t, []string{"/referrer", "/referrer"}, s.referrers())
		assert.Equal(t, s.responseCode, resp.StatusCode)
	}
	_ = s.listener.Close()
}

func TestFetch_Get(t *testing.T) {
	s := startServer()
	f := NewFetcher(
		WithTimeout(time.Second),
		WithUserAgent("Bot/1"),
		WithAccept("application/binary"),
		WithHeadRequests(false),
	)
	r := &Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   s.listener.Addr().String(),
		},
		AcceptableContentTypes: map[string]struct{}{
			"text/html": {},
		},
	}
	s.responseCode = 403
	if resp, err := f.Fetch(r); assert.NoError(t, err) {
		if body, err := ioutil.ReadAll(resp.Body); assert.NoError(t, err) {
			assert.Equal(t, []byte(testHTML), body)
		}
		_ = resp.Body.Close()
		assert.Equal(t, []string{"GET"}, s.methods())
		assert.Equal(t, []string{"Bot/1"}, s.userAgents())
		assert.Equal(t, []string{"application/binary"}, s.accepts())
		assert.Equal(t, []string{""}, s.referrers())
		assert.Equal(t, s.responseCode, resp.StatusCode)
	}
	_ = s.listener.Close()
}

func TestFailedHead(t *testing.T) {
	s := startServer()
	req := &Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   s.listener.Addr().String(),
			Path:   "/",
		},
	}
	s.failHeadRequest = true
	f := NewFetcher(WithTimeout(time.Second))
	if resp, err := f.Fetch(req); assert.NoError(t, err) {
		if body, err := ioutil.ReadAll(resp.Body); assert.NoError(t, err) {
			assert.Equal(t, []byte(testHTML), body)
		}
		_ = resp.Body.Close()
		assert.Equal(t, []string{"GET"}, s.methods())
		assert.Equal(t, s.responseCode, resp.StatusCode)
	}
	_ = s.listener.Close()
}

func TestSlowHead(t *testing.T) {
	s := startServer()
	req := &Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   s.listener.Addr().String(),
			Path:   "/slow",
		},
	}
	f := NewFetcher(WithTimeout(time.Second))
	_, err := f.Fetch(req)
	assert.Error(t, err)
	_ = s.listener.Close()
}

func TestSlowGet(t *testing.T) {
	s := startServer()
	req := &Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   s.listener.Addr().String(),
			Path:   "/slow",
		},
	}
	f := NewFetcher(WithTimeout(time.Second))
	_, err := f.Fetch(req)
	assert.Error(t, err)
	_ = s.listener.Close()
}

func TestAccept(t *testing.T) {
	s := startServer()
	req := &Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   s.listener.Addr().String(),
		},
	}
	f := NewFetcher(WithAccept("application/binary"))
	if resp, err := f.Fetch(req); assert.NoError(t, err) {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
		assert.Equal(t, []string{"application/binary"}, s.accepts())
	}
	_ = s.listener.Close()
}

func TestContentType(t *testing.T) {
	s := startServer()
	req := &Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   s.listener.Addr().String(),
		},
		AcceptableContentTypes: map[string]struct{}{"application/binary": {}},
	}
	// Acceptable content type makes sense with HEAD requests only
	f := NewFetcher(WithHeadRequests(true))
	_, err := f.Fetch(req)
	assert.Equal(t, ErrBadContentType, err)
	_ = s.listener.Close()
}

type testServer struct {
	listener        net.Listener
	contentType     string
	responseCode    int
	failHeadRequest bool
	seenMethods     []string
	seenUAs         []string
	seenReferrers   []string
	seenAccept      []string
}

func (t *testServer) methods() []string {
	m := make([]string, len(t.seenMethods))
	copy(m, t.seenMethods)
	t.seenMethods = t.seenMethods[0:0]
	return m
}

func (t *testServer) userAgents() []string {
	m := make([]string, len(t.seenUAs))
	copy(m, t.seenUAs)
	t.seenUAs = t.seenUAs[0:0]
	return m
}

func (t *testServer) referrers() []string {
	m := make([]string, len(t.seenReferrers))
	copy(m, t.seenReferrers)
	t.seenReferrers = t.seenReferrers[0:0]
	return m
}

func (t *testServer) accepts() []string {
	m := make([]string, len(t.seenAccept))
	copy(m, t.seenAccept)
	t.seenAccept = t.seenAccept[0:0]
	return m
}

func (t *testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.seenMethods = append(t.seenMethods, r.Method)
	t.seenUAs = append(t.seenUAs, r.Header.Get("User-Agent"))
	t.seenReferrers = append(t.seenReferrers, r.Header.Get("Referer"))
	t.seenAccept = append(t.seenAccept, r.Header.Get("Accept"))
	if t.failHeadRequest && r.Method == "HEAD" {
		_ = t.listener.Close()
		return
	}
	if r.URL.Path == "/slow" {
		time.Sleep(time.Second * 2)
		return
	}
	w.Header().Add("Content-Type", t.contentType)
	w.WriteHeader(t.responseCode)
	if t.contentType == "text/css" {
		if _, err := w.Write([]byte(testCSS)); err != nil {
			panic(err)
		}
	} else if t.contentType == "text/html" {
		if _, err := w.Write([]byte(testHTML)); err != nil {
			panic(err)
		}
	} else {
		panic("content type is not set")
	}
}

func startServer() *testServer {
	serv := testServer{
		contentType:  "text/html",
		responseCode: 201,
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	serv.listener = l
	go func() { _ = http.Serve(l, &serv) }()
	return &serv
}

// language=CSS
const testCSS = `body {
   background-color: red;
}
`

// language=HTML
const testHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Foo Bar</title>
</head>
<body>
</body>
</html>
`
