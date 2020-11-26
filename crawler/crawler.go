package crawler

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/dmitry-vovk/wcrowler/crawler/limiter"
	"github.com/dmitry-vovk/wcrowler/crawler/page_fetcher"
	"github.com/dmitry-vovk/wcrowler/crawler/page_parser"
	"github.com/dmitry-vovk/wcrowler/crawler/queue"
	"github.com/dmitry-vovk/wcrowler/crawler/queue/filter"
	"github.com/temoto/robotstxt"
)

type Crawler struct {
	cfg    *Config
	doneC  chan struct{}
	doStop bool
}

const defaultMaxParallelRequests = 1

// New creates an instance of Crawler
func New(cfg *Config) *Crawler {
	return &Crawler{
		cfg:   cfg,
		doneC: make(chan struct{}),
	}
}

// Run starts the crawling
func (c *Crawler) Run() error {
	if err := c.validateConfig(); err != nil {
		return err
	}
	if err := c.setup(); err != nil {
		return err
	}
	log.Printf("Starting from %s", c.cfg.SeedURL)
	go c.processor()
	c.run()
	<-c.doneC
	for i, link := range filter.GetCollectedLinks() {
		log.Printf("%d. %s", i+1, link)
	}
	return nil
}

// setup initializes dependencies according to configuration
func (c *Crawler) setup() error {
	u, err := url.Parse(c.cfg.SeedURL)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		return errors.New("missing URL scheme")
	}
	if u.Hostname() == "" {
		return errors.New("invalid URL")
	}
	filter.SetBaseDomain(u.Hostname())
	filter.SetMaxPagesCallback(c.cfg.MaxPages, c.stop)
	if c.cfg.MaxParallelRequests == 0 {
		limiter.SetLimit(defaultMaxParallelRequests)
	} else {
		limiter.SetLimit(c.cfg.MaxParallelRequests)
	}
	if !c.cfg.IgnoreRobotsTxt {
		if robots := c.fetchRobots(); robots != nil {
			filter.SetRobots(robots, c.cfg.UserAgent)
		}
	}
	return nil
}

// run kicks off the crawling
func (c *Crawler) run() {
	queue.Enqueue(c.cfg.SeedURL)
}

// stop would interrupt crawling when called
func (c *Crawler) stop() {
	c.doStop = true
}

// processor sequentially processes page crawls
func (c *Crawler) processor() {
	for {
		if c.doStop {
			log.Printf("Requested to stop")
			break
		}
		link, ok := queue.Next()
		if !ok {
			log.Println("No more links")
			break
		}
		log.Printf("Starting processing %s", link)
		go c.processLink(link)
	}
	close(c.doneC)
}

// processLink handles single page crawling
func (c *Crawler) processLink(link string) {
	limiter.Start()
	defer limiter.Finish()
	u, err := url.Parse(link)
	if err != nil {
		// if parsing fails here, we have a bug somewhere before
		panic(err)
	}
	request := page_fetcher.Request{
		URL:           u,
		HTTPReferrer:  "", // TODO populate referrer
		UserAgent:     c.cfg.UserAgent,
		DoHeadRequest: c.cfg.DoHeadRequests,
		AcceptableContentTypes: map[string]struct{}{
			"text/html": {},
		},
	}
	response, err := page_fetcher.Fetch(&request)
	if err != nil {
		log.Printf("Error fetching page %q: %s", link, err)
		return
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Got %d status code from %q", response.StatusCode, link)
		return
	}
	page, err := page_parser.Parse(response.Body)
	_ = response.Body.Close()
	if err != nil {
		log.Printf("Error parsing response from %q: %s", link, err)
		return
	}
	for i := range page.Links {
		if pageLink, err := url.Parse(page.Links[i]); err == nil {
			resolvedURL := u.ResolveReference(pageLink)
			queue.Enqueue(resolvedURL.String())
		}
	}
}

// validateConfig makes sure the config makes sense
func (c *Crawler) validateConfig() error {
	if c.cfg == nil {
		return errors.New("empty config")
	}
	if c.cfg.SeedURL == "" {
		return errors.New("empty seed URL")
	}
	return nil
}

// fetchRobots tries to get 'robots.txt' file for the seed URL
func (c *Crawler) fetchRobots() *robotstxt.RobotsData {
	u, err := url.Parse(c.cfg.SeedURL)
	if err != nil {
		panic("should've checked seed URL!")
	}
	u.Path = "/robots.txt"
	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("Error fetching robots.txt: %s", err)
		return nil
	}
	r, err := robotstxt.FromResponse(resp)
	if err != nil {
		log.Printf("Error parsing robots.txt: %s", err)
		return nil
	}
	return r
}
