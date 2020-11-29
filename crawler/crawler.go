package crawler

import (
	"errors"
	"log"

	"github.com/dmitry-vovk/wcrawler/crawler/types"
)

// Crawler is a page crawler
type Crawler struct {
	maxPages            uint64
	maxParallelRequests uint
	fetcher             types.Fetcher
	filter              types.Filter
	resultCallback      func(string, []string) // Callback function to send page crawl results
	limiterC            chan struct{}          // Simple limiter
	queuedLinksC        chan crawlJob          // URLs to be processed
	processedLinksC     chan crawlResult       // URLs that done processing
	processingLinks     map[string]struct{}    // Links that are currently being processed
	processedLinks      map[string]struct{}    // Visited links
	doneC               chan struct{}          // Done signal
	pagesN              uint64
	finished            bool
}

const defaultMaxParallelRequests = 1

// New creates an instance of Crawler
func New(fetcher types.Fetcher, filter types.Filter, pageCrawlResultCallback func(string, []string)) *Crawler {
	return &Crawler{
		fetcher:        fetcher,
		filter:         filter,
		resultCallback: pageCrawlResultCallback,
	}
}

// MaxPages sets the maximum number of pages to crawl
func (c *Crawler) MaxPages(maxPages uint64) *Crawler {
	c.maxPages = maxPages
	return c
}

// MaxParallelRequests limit the number or parallel requests
func (c *Crawler) MaxParallelRequests(maxParallelRequests uint) *Crawler {
	c.maxParallelRequests = maxParallelRequests
	return c
}

// Run starts the crawling and blocks until finished
func (c *Crawler) Run(seedURL string) error {
	if c.fetcher == nil {
		return errors.New("fetcher not set")
	}
	if c.filter == nil {
		return errors.New("filter not set")
	}
	seed, ok := c.filter.Filter(seedURL)
	if !ok {
		return errors.New("bad seed URL")
	}
	if c.resultCallback == nil {
		log.Print("Results callback function not set")
	}
	if c.maxParallelRequests == 0 {
		c.limiterC = make(chan struct{}, defaultMaxParallelRequests)
	} else {
		c.limiterC = make(chan struct{}, c.maxParallelRequests)
	}
	c.queuedLinksC = make(chan crawlJob)
	c.processedLinksC = make(chan crawlResult)
	c.processingLinks = make(map[string]struct{})
	c.processedLinks = make(map[string]struct{})
	c.doneC = make(chan struct{})
	log.Printf("Starting from %s", seed)
	go c.processor()
	c.queuedLinksC <- crawlJob{Link: seed}
	<-c.doneC
	return nil
}
