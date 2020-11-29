package url_filter

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/purell"
)

// RobotsChecker defines a method to tell if this path allowed to be visited by this agent
type RobotsChecker interface {
	TestAgent(path, agent string) bool
}

// NormalizingFilter is a URL filter with basic normalisation rules
type NormalizingFilter struct {
	baseDomain string
	allowWWW   bool
	robots     RobotsChecker
	userAgent  string
}

var (
	// Set some sensible normalisation flags, may be tweaked to suit a specific site requirements
	URLNormalisationFlags = purell.FlagsUsuallySafeGreedy |
		purell.FlagRemoveDuplicateSlashes |
		purell.FlagRemoveFragment
)

// NewFilter returns and instance of NormalizingFilter with sane defaults
func NewFilter(baseDomain string) *NormalizingFilter {
	f := NormalizingFilter{
		baseDomain: baseDomain,
	}
	return &f
}

// AllowWWWPrefix sets switch for handling commonly used "www" prefix
func (f *NormalizingFilter) AllowWWWPrefix(allow bool) *NormalizingFilter {
	f.allowWWW = allow
	return f
}

// WithRobots injects RobotsChecker
func (f *NormalizingFilter) WithRobots(r RobotsChecker, agent string) *NormalizingFilter {
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
	link = purell.NormalizeURL(u, URLNormalisationFlags)
	if u.Path == "" {
		return link + "/", true
	}
	return link, true
}
