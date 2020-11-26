package page_fetcher

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetch_Get(t *testing.T) {
	s := startServer()
	req := &Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   s.listener.Addr().String(),
		},
		AcceptableContentTypes: map[string]struct{}{
			"text/html": {},
		},
	}
	{
		req.UserAgent = "Bot/1"
		req.HTTPReferrer = "/1"
		resp, err := Fetch(req)
		if assert.NoError(t, err) {
			if body, err := ioutil.ReadAll(resp.Body); assert.NoError(t, err) {
				assert.Equal(t, testHTML, string(body))
			}
			_ = resp.Body.Close()
			assert.Equal(t, 201, resp.StatusCode)
		}
		assert.Equal(t, []string{"GET"}, s.methods())
		assert.Equal(t, []string{"Bot/1"}, s.userAgents())
		assert.Equal(t, []string{"/1"}, s.referrers())
	}
	{
		req.UserAgent = "Bot/2"
		req.HTTPReferrer = "/2"
		s.responseCode = 412
		resp, err := Fetch(req)
		if assert.NoError(t, err) {
			_ = resp.Body.Close()
			assert.Equal(t, 412, resp.StatusCode)
		}
		assert.Equal(t, []string{"GET"}, s.methods())
		assert.Equal(t, []string{"Bot/2"}, s.userAgents())
		assert.Equal(t, []string{"/2"}, s.referrers())
	}
	{
		req.UserAgent = "Bot/3"
		req.HTTPReferrer = "/3"
		req.DoHeadRequest = true
		resp, err := Fetch(req)
		if assert.NoError(t, err) {
			if resp.Body != nil {
				_ = resp.Body.Close()
			}
			assert.Equal(t, 412, resp.StatusCode)
		}
		assert.Equal(t, []string{"HEAD", "GET"}, s.methods())
		assert.Equal(t, []string{"Bot/3", "Bot/3"}, s.userAgents())
		assert.Equal(t, []string{"/3", "/3"}, s.referrers())
	}
	{
		req.UserAgent = "Bot/4"
		req.HTTPReferrer = "/4"
		req.DoHeadRequest = false
		s.contentType = "text/css"
		resp, err := Fetch(req)
		if assert.NoError(t, err) {
			if body, err := ioutil.ReadAll(resp.Body); assert.NoError(t, err) {
				assert.Equal(t, testCSS, string(body))
			}
			_ = resp.Body.Close()
			assert.Equal(t, 412, resp.StatusCode)
			assert.Equal(t, "text/css", resp.Headers.Get("Content-Type"))
		}
		assert.Equal(t, []string{"GET"}, s.methods())
		assert.Equal(t, []string{"Bot/4"}, s.userAgents())
		assert.Equal(t, []string{"/4"}, s.referrers())
	}
	_ = s.listener.Close()
}

type testServer struct {
	listener      net.Listener
	contentType   string
	responseCode  int
	seenMethods   []string
	seenUAs       []string
	seenReferrers []string
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

func (t *testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.seenMethods = append(t.seenMethods, r.Method)
	t.seenUAs = append(t.seenUAs, r.Header.Get("User-Agent"))
	t.seenReferrers = append(t.seenReferrers, r.Header.Get("Referer"))
	w.Header().Add("Content-Type", t.contentType)
	w.WriteHeader(t.responseCode)
	if t.contentType == "text/css" {
		if _, err := w.Write([]byte(testCSS)); err != nil {
			panic(err)
		}
	} else {
		if _, err := w.Write([]byte(testHTML)); err != nil {
			panic(err)
		}
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
