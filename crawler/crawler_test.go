package crawler

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/dmitry-vovk/wcrawler/crawler/page_fetcher"
	"github.com/dmitry-vovk/wcrawler/crawler/types"
	"github.com/stretchr/testify/assert"
)

func TestCrawlerBuild(t *testing.T) {
	const (
		maxPages            uint64 = 2
		maxParallelRequests uint   = 76
	)
	var results []struct {
		Link  string
		Links []string
	}
	c := New(tFetcher, tFilter, func(link string, links []string) {
		results = append(results, struct {
			Link  string
			Links []string
		}{link, links})
	}).MaxPages(maxPages).MaxParallelRequests(maxParallelRequests)
	assert.Equal(t, maxPages, c.maxPages)
	assert.Equal(t, maxParallelRequests, c.maxParallelRequests)
	assert.Equal(t, tFetcher, c.fetcher)
	assert.Equal(t, tFilter, c.filter)
	tFetcher.(*testFetcher).statusCode = 200
	if err := c.Run("http://example.com/"); assert.NoError(t, err) {
		assert.Equal(t, []struct {
			Link  string
			Links []string
		}{
			{
				Link:  "http://example.com/",
				Links: []string{"http://example.com/"},
			},
		}, results)
	}
}

func TestCrawlerIncompleteBuild1(t *testing.T) {
	c := New(nil, nil, nil)
	assert.Error(t, c.Run(""))
}

func TestCrawlerIncompleteBuild2(t *testing.T) {
	c := New(tFetcher, nil, nil)
	assert.Error(t, c.Run(""))
}

func TestCrawler_BadSeed(t *testing.T) {
	c := New(tFetcher, tFilter, nil)
	assert.Equal(t, tFetcher, c.fetcher)
	assert.Equal(t, tFilter, c.filter)
	tFetcher.(*testFetcher).statusCode = 200
	if err := c.Run(""); assert.Error(t, err) {
		assert.Equal(t, "bad seed URL", err.Error())
	}
}

func TestCrawler_Run_200(t *testing.T) {
	var results []struct {
		Link  string
		Links []string
	}
	c := New(tFetcher, tFilter, func(link string, links []string) {
		results = append(results, struct {
			Link  string
			Links []string
		}{link, links})
	})
	assert.Equal(t, tFetcher, c.fetcher)
	assert.Equal(t, tFilter, c.filter)
	tFetcher.(*testFetcher).statusCode = 200
	if err := c.Run("http://example.com/"); assert.NoError(t, err) {
		assert.Equal(t, []struct {
			Link  string
			Links []string
		}{
			{
				Link:  "http://example.com/",
				Links: []string{"http://example.com/"},
			},
		}, results)
	}
}

func TestCrawler_Run_404(t *testing.T) {
	var results []struct {
		Link  string
		Links []string
	}
	c := New(tFetcher, tFilter, func(link string, links []string) {
		results = append(results, struct {
			Link  string
			Links []string
		}{link, links})
	})
	tFetcher.(*testFetcher).statusCode = 404
	if err := c.Run("http://example.com/"); assert.NoError(t, err) {
		assert.Equal(t, []struct {
			Link  string
			Links []string
		}{
			{
				Link:  "http://example.com/",
				Links: []string{},
			},
		}, results)
	}
}

func TestCrawler_Run_Error(t *testing.T) {
	var results []struct {
		Link  string
		Links []string
	}
	c := New(tFetcher, tFilter, func(link string, links []string) {
		results = append(results, struct {
			Link  string
			Links []string
		}{link, links})
	})
	tFetcher.(*testFetcher).statusCode = 200
	tFetcher.(*testFetcher).err = errors.New("some expected error")
	if err := c.Run("http://example.com/"); assert.NoError(t, err) {
		assert.Equal(t, []struct {
			Link  string
			Links []string
		}{
			{
				Link:  "http://example.com/",
				Links: []string{},
			},
		}, results)
	}
}

func TestCrawler_Parse_Error(t *testing.T) {
	var results []struct {
		Link  string
		Links []string
	}
	c := New(tFetcher, tFilter, func(link string, links []string) {
		results = append(results, struct {
			Link  string
			Links []string
		}{link, links})
	})
	tFetcher.(*testFetcher).statusCode = 200
	tFetcher.(*testFetcher).nilBody = true
	if err := c.Run("http://example.com/"); assert.NoError(t, err) {
		assert.Equal(t, []struct {
			Link  string
			Links []string
		}{
			{
				Link:  "http://example.com/",
				Links: []string{},
			},
		}, results)
	}
}

var tFetcher types.Fetcher = &testFetcher{}

type testFetcher struct {
	err        error
	statusCode int
	nilBody    bool
}

func (t testFetcher) Fetch(_ *page_fetcher.Request) (*page_fetcher.Response, error) {
	if t.err != nil {
		return nil, t.err
	}
	body := ioutil.NopCloser(bytes.NewBuffer([]byte(testHTML)))
	if t.nilBody {
		body = ioutil.NopCloser(bytes.NewBuffer([]byte(`<html><body><aef<eqf>>>qq></body></ht>`)))
	}
	return &page_fetcher.Response{
		StatusCode: t.statusCode,
		Headers:    nil,
		Body:       body,
	}, nil
}

var tFilter types.Filter = &testFilter{}

type testFilter struct{}

func (t testFilter) Filter(link string) (string, bool) {
	return link, link != ""
}

// language=HTML
const testHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Foo Bar</title>
</head>
<body>
	<a href="/">Home</a>
	<a href="http://example.com/">Home</a>
</body>
</html>`
