package crawler

import (
	"log"
	"sync/atomic"
	"time"
)

// processor sequentially processes page crawls
func (c *Crawler) processor() {
out:
	for {
		select {
		case job := <-c.queuedLinksC:
			if _, ok := c.processedLinks[job.Link]; !ok {
				c.processingLinks[job.Link] = struct{}{}
				c.processedLinks[job.Link] = struct{}{}
				go c.processJob(job)
			}
		case result := <-c.processedLinksC:
			delete(c.processingLinks, result.Link)
			if c.resultCallback != nil {
				c.resultCallback(result.Link, result.CollectLinks())
			}
			if atomic.LoadUint64(&c.pagesN) >= c.maxPages {
				break out
			}
			if len(c.processingLinks) == 0 {
				break out
			}
		}
	}
	log.Printf("Pages visited: %d", c.pagesN)
	close(c.doneC)
	c.finished = true
}

// processJob handles single page crawling
func (c *Crawler) processJob(link crawlJob) {
	c.limiterC <- struct{}{}
	start := time.Now()
	log.Printf("Starting link: %s", link.Link)
	task := newTask(link)
	result := task.Process(c.fetcher)
	if result.Error != nil {
		log.Printf("Error processing link %q: %s", link.Link, result.Error)
	} else {
		for i := range result.Links {
			if cleanLink, ok := c.filter.Filter(result.Links[i].String()); ok {
				c.queuedLinksC <- crawlJob{Link: cleanLink, Referrer: link.Referrer}
			}
		}
		log.Printf("Finished link in %s: %s", time.Since(start), link.Link)
	}
	c.processedLinksC <- result
	<-c.limiterC
	atomic.AddUint64(&c.pagesN, 1)
}
