package queue

import (
	"log"

	"github.com/dmitry-vovk/wcrowler/crawler/queue/filter"
)

var (
	// This queue should be large enough not to cause a deadlock
	linksC = make(chan string, 1000)
)

func Enqueue(link string) {
	l, ok := filter.Filter(link)
	if !ok {
		return
	}
	select {
	case linksC <- l:
	default:
		// TODO Queue is full, handle it
		log.Println("Link queue is full")
	}
}

func Next() (string, bool) {
	link, ok := <-linksC
	return link, ok
}
