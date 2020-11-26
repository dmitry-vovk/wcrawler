package crawler

import (
	"encoding/json"
	"os"
)

type Config struct {
	// Starting URL
	SeedURL string `json:"seed_url"`
	// List of upper level domains to allow:
	// e.g. with "www" treat example.com and www.example.com as the same domain
	AllowPrefixes []string `json:"allow_prefixes"`
	// Whether to take in account robots.txt rules, better 'yes' to be polite
	IgnoreRobotsTxt bool `json:"ignore_robots_txt"`
	// Do HEAD requests before GET requests to avoid fetching inappropriate files
	DoHeadRequests bool `json:"do_head_requests"`
	// HTTP user agent string to use
	UserAgent string `json:"user_agent"`
	// Do not crawl more than this number of pages
	MaxPages int `json:"max_pages"`
	// How many requests to allow to run in parallel
	MaxParallelRequests int `json:"max_parallel_requests"`
}

// Read returns config read from file or error
func Read(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		// We do not care about errors here
		_ = f.Close()
	}()
	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}
