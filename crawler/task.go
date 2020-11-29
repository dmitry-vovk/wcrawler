package crawler

import (
	"log"
	"net/http"
	"net/url"

	"github.com/dmitry-vovk/wcrawler/crawler/page_fetcher"
	"github.com/dmitry-vovk/wcrawler/crawler/page_parser"
	"github.com/dmitry-vovk/wcrawler/crawler/types"
)

type CrawlJob struct {
	Link     string
	Referrer string
}

type CrawlResult struct {
	Link  string
	Links []*url.URL
}

func (cr CrawlResult) CollectLinks() map[string]struct{} {
	links := make(map[string]struct{})
	for i := range cr.Links {
		links[cr.Links[i].String()] = struct{}{}
	}
	return links
}

type Task struct {
	job CrawlJob
}

func NewTask(link CrawlJob) *Task {
	return &Task{job: link}
}

func (t *Task) Process(fetcher types.Fetcher) (result CrawlResult) {
	u, err := url.Parse(t.job.Link)
	if err != nil {
		// if parsing fails here, we have a bug somewhere before
		panic(err)
	}
	request := page_fetcher.Request{
		URL:          u,
		HTTPReferrer: t.job.Referrer,
		AcceptableContentTypes: map[string]struct{}{
			"text/html": {},
		},
	}
	result.Link = t.job.Link
	response, err := fetcher.Fetch(&request)
	if err != nil {
		log.Printf("Error fetching page %q: %s", t.job, err)
		return
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		log.Printf("Got %d status code from %q", response.StatusCode, t.job)
		return
	}
	page, err := page_parser.Parse(response.Body)
	if err != nil {
		log.Printf("Error parsing response from %q: %s", t.job, err)
		return
	}
	for i := range page.Links {
		if pageLink, err := url.Parse(page.Links[i]); err == nil {
			result.Links = append(result.Links, u.ResolveReference(pageLink))
		}
	}
	return
}
