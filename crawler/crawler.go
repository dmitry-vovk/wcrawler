package crawler

import (
	"errors"
	"log"

	"github.com/dmitry-vovk/wcrawler/crawler/types"
)

type Crawler struct {
	maxPages            uint64
	maxParallelRequests uint
	fetcher             types.Fetcher
	filter              types.Filter
	limiterC            chan struct{}                  // Simple limiter
	queuedLinksC        chan CrawlJob                  // URLs to be processed
	processedLinksC     chan CrawlResult               // URLs that done processing
	processingLinks     map[string]struct{}            // Links that are currently being processed
	processedLinks      map[string]struct{}            // Visited links
	collectedLinks      map[string]map[string]struct{} // Collection of visited pages with found links
	doneC               chan struct{}                  // Done signal
	pagesN              uint64
}

const defaultMaxParallelRequests = 1

// New creates an instance of Crawler
func New() *Crawler {
	return &Crawler{
		queuedLinksC:    make(chan CrawlJob),
		processedLinksC: make(chan CrawlResult),
		processingLinks: make(map[string]struct{}),
		processedLinks:  make(map[string]struct{}),
		collectedLinks:  make(map[string]map[string]struct{}),
		doneC:           make(chan struct{}),
	}
}

func (c *Crawler) MaxPages(maxPages uint64) *Crawler {
	c.maxPages = maxPages
	return c
}

func (c *Crawler) MaxParallelRequests(maxParallelRequests uint) *Crawler {
	c.maxParallelRequests = maxParallelRequests
	return c
}

func (c *Crawler) WithFetcher(fetcher types.Fetcher) *Crawler {
	c.fetcher = fetcher
	return c
}

func (c *Crawler) WithFilter(filter types.Filter) *Crawler {
	c.filter = filter
	return c
}

// Run starts the crawling
func (c *Crawler) Run(seedURL string) error {
	if c.maxParallelRequests == 0 {
		c.limiterC = make(chan struct{}, defaultMaxParallelRequests)
	} else {
		c.limiterC = make(chan struct{}, c.maxParallelRequests)
	}
	seed, ok := c.filter.Filter(seedURL)
	if !ok {
		return errors.New("bad seed URL")
	}
	log.Printf("Starting from %s", seed)
	go c.processor()
	c.queuedLinksC <- CrawlJob{Link: seed}
	<-c.doneC
	return nil
}

func (c *Crawler) Results() map[string]map[string]struct{} {
	return c.collectedLinks
}
