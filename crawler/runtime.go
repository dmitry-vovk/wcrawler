package crawler

import (
	"log"
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
			c.collectedLinks[result.Link] = result.CollectLinks()
			if c.pagesN >= c.maxPages {
				break out
			}
			if len(c.processingLinks) == 0 {
				break out
			}
		}
	}
	log.Printf("Processor exited")
	log.Printf("Pages visited: %d", c.pagesN)
	close(c.doneC)
}

// processLink handles single page crawling
func (c *Crawler) processJob(link CrawlJob) {
	c.limiterC <- struct{}{}
	log.Printf("Starting: %s", link)
	task := NewTask(link)
	result := task.Process(c.fetcher)
	for i := range result.Links {
		if cleanLink, ok := c.filter.Filter(result.Links[i].String()); ok {
			c.queuedLinksC <- CrawlJob{Link: cleanLink, Referrer: link.Referrer}
		}
	}
	c.processedLinksC <- result
	<-c.limiterC
	c.pagesN++
	log.Printf("Finished: %s", link)
}
