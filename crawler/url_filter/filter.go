package url_filter

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/purell"
)

type RobotsTxt interface {
	TestAgent(path, agent string) bool
}

type NormalizingFilter struct {
	baseDomain string
	allowWWW   bool
	robots     RobotsTxt
	userAgent  string
}

var (
	URLNormalizationFlags = purell.FlagsUsuallySafeGreedy |
		purell.FlagRemoveDuplicateSlashes |
		purell.FlagRemoveFragment
)

func NewFilter(baseDomain string) *NormalizingFilter {
	f := NormalizingFilter{
		baseDomain: baseDomain,
	}
	return &f
}

func (f *NormalizingFilter) AllowWWWPrefix(allow bool) *NormalizingFilter {
	f.allowWWW = allow
	return f
}

func (f *NormalizingFilter) WithRobots(r RobotsTxt, agent string) *NormalizingFilter {
	f.robots, f.userAgent = r, agent
	return f
}

// Filter returns normalized link and/or tells if the link is ok to use
func (f *NormalizingFilter) Filter(link string) (string, bool) {
	u, err := url.Parse(link)
	if err != nil {
		return "", false
	}
	if u.Path == "" {
		u.Path = "/"
	}
	if f.baseDomain != "" {
		host := u.Hostname()
		if f.allowWWW {
			if strings.TrimPrefix(host, "www.") != strings.TrimPrefix(f.baseDomain, "www.") {
				return "", false
			}
		} else if host != f.baseDomain {
			return "", false
		}
	}
	if f.robots != nil {
		if !f.robots.TestAgent(u.Path, f.userAgent) {
			log.Printf("Visiting %s disallowed by robots.txt", link)
			return "", false
		}
	}
	link = purell.NormalizeURL(u, URLNormalizationFlags)
	if u.Path == "" {
		return link + "/", true
	}
	return link, true
}
