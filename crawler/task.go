package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/pkg/errors"

	"github.com/dmitry-vovk/wcrawler/crawler/page_fetcher"
	"github.com/dmitry-vovk/wcrawler/crawler/page_parser"
	"github.com/dmitry-vovk/wcrawler/crawler/types"
)

type crawlJob struct {
	Link     string
	Referrer string
}

type crawlResult struct {
	Link          string
	CanonicalLink string
	Links         []*url.URL
	Error         error
}

func (cr crawlResult) CollectLinks() []string {
	links := make([]string, 0, len(cr.Links))
	uniqueLinks := make(map[string]struct{})
	for i := range cr.Links {
		link := cr.Links[i].String()
		if _, ok := uniqueLinks[link]; !ok {
			links = append(links, link)
			uniqueLinks[link] = struct{}{}
		}
	}
	sort.Strings(links)
	return links
}

type task struct {
	job crawlJob
}

func newTask(link crawlJob) *task {
	return &task{job: link}
}

func (t *task) Process(fetcher types.Fetcher) (result crawlResult) {
	result.Link = t.job.Link
	u, err := url.Parse(t.job.Link)
	if err != nil {
		result.Error = errors.Wrap(err, "URL parse error")
		return
	}
	request := page_fetcher.Request{
		URL:          u,
		HTTPReferrer: t.job.Referrer,
		AcceptableContentTypes: map[string]struct{}{
			"text/html": {},
		},
	}
	response, err := fetcher.Fetch(&request)
	if err != nil {
		result.Error = errors.Wrap(err, "fetch")
		return
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		result.Error = fmt.Errorf("got status code %d", response.StatusCode)
		return
	}
	page, err := page_parser.Parse(response.Body)
	if err != nil {
		result.Error = errors.Wrap(err, "parse")
		return
	}
	for i := range page.Links {
		if pageLink, err := url.Parse(page.Links[i]); err == nil {
			result.Links = append(result.Links, u.ResolveReference(pageLink))
		}
	}
	result.CanonicalLink = page.CanonicalURL
	return
}
