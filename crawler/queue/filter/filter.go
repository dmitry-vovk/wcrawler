package filter

import (
	"log"
	"net/url"
	"sort"
	"sync"

	"github.com/PuerkitoBio/purell"
	"github.com/temoto/robotstxt"
)

var (
	baseDomain            string
	maxPages              int
	stopFn                func()
	robots                *robotstxt.RobotsData
	userAgent             string
	seenURLs              = make(map[string]struct{})
	seenURLsM             sync.Mutex
	URLNormalizationFlags = purell.FlagsUsuallySafeGreedy |
		purell.FlagRemoveDuplicateSlashes |
		purell.FlagRemoveFragment
)

func SetBaseDomain(domain string) {
	baseDomain = domain
}

func SetMaxPagesCallback(max int, callbackFn func()) {
	maxPages, stopFn = max, callbackFn
}

func SetRobots(r *robotstxt.RobotsData, agent string) {
	robots, userAgent = r, agent
}

// Filter returns normalized link and/or tells if the link has been seen
func Filter(link string) (string, bool) {
	link = normalize(link)
	u, err := url.Parse(link)
	if err != nil {
		return "", false
	} else if host := u.Hostname(); host != baseDomain {
		return "", false
	}
	if robots != nil {
		if !robots.TestAgent(u.Path, userAgent) {
			log.Printf("Visiting %s disallowed by robots.txt", link)
			return "", false
		}
	}
	seenURLsM.Lock()
	defer seenURLsM.Unlock()
	if _, ok := seenURLs[link]; ok {
		return "", false
	}
	log.Printf("New URL: %s", link)
	seenURLs[link] = struct{}{}
	if maxPages != 0 && maxPages <= len(seenURLs) && stopFn != nil {
		stopFn()
	}
	return link, true
}

func GetCollectedLinks() []string {
	seenURLsM.Lock()
	defer seenURLsM.Unlock()
	var pages []string
	for link := range seenURLs {
		pages = append(pages, link)
	}
	sort.Strings(pages)
	return pages
}

// normalize cleans up and standardizes links
func normalize(link string) string {
	return purell.MustNormalizeURLString(link, URLNormalizationFlags)
}
