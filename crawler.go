package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dmitry-vovk/wcrawler/crawler"
	"github.com/dmitry-vovk/wcrawler/crawler/page_fetcher"
	"github.com/dmitry-vovk/wcrawler/crawler/url_filter"
	"github.com/temoto/robotstxt"
)

const (
	defaultConfigFile = "config.json"
	configFileEnv     = "CRAWLER_CONFIG"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg := mustLoadConfig()
	start := time.Now()
	if err := mustBuildCrawler(cfg, printResults).Run(cfg.SeedURL); err != nil {
		log.Printf("Error running crawler: %s\n", err)
	} else {
		log.Printf("Crawler finished in %s\n", time.Since(start))
	}
}

func mustBuildCrawler(cfg *Config, resultsCallback func(string, []string)) *crawler.Crawler {
	// Validate seed URL
	u, err := url.Parse(cfg.SeedURL)
	if err != nil {
		log.Printf("Error parsing seed URL: %s", err)
		os.Exit(2)
	} else if u.Hostname() == "" {
		log.Print("Invalid seed URL")
		os.Exit(2)
	}
	// Initialize dependencies
	filter := url_filter.
		NewFilter(u.Hostname()).
		AllowWWWPrefix(cfg.AllowWWWPrefix)
	if !cfg.IgnoreRobotsTxt {
		u.Path = "/robots.txt"
		if robots := fetchRobots(u.String()); robots != nil {
			filter.WithRobots(robots, cfg.UserAgent)
		}
	}
	fetcher := page_fetcher.NewFetcher(
		page_fetcher.WithUserAgent(cfg.UserAgent),
		page_fetcher.WithHeadRequests(cfg.DoHeadRequests),
	)
	// Assemble a crawler instance
	return crawler.
		New(fetcher, filter, resultsCallback).
		MaxPages(cfg.MaxPages).
		MaxParallelRequests(cfg.MaxParallelRequests)
}

// Config contains all the variables needed for crawler
type Config struct {
	// Starting URL
	SeedURL string `json:"seed_url"`
	// List of upper level domains to allow:
	// e.g. with "www" treat example.com and www.example.com as the same domain
	AllowWWWPrefix bool `json:"allow_www_prefix"`
	// Whether to take in account robots.txt rules
	IgnoreRobotsTxt bool `json:"ignore_robots_txt"`
	// Do HEAD requests before GET requests to avoid fetching inappropriate links
	DoHeadRequests bool `json:"do_head_requests"`
	// HTTP user agent string to use
	UserAgent string `json:"user_agent"`
	// Do not crawl more than this number of pages
	MaxPages uint64 `json:"max_pages"`
	// How many requests to allow to run in parallel
	MaxParallelRequests uint `json:"max_parallel_requests"`
}

func mustLoadConfig() *Config {
	configFile := defaultConfigFile
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	} else if cfgPath := os.Getenv(configFileEnv); cfgPath != "" {
		configFile = cfgPath
	}
	cfg, err := readConfig(configFile)
	if err != nil {
		log.Printf("Error reading config file %q: %s", configFile, err)
		os.Exit(1)
	}
	return cfg
}

// readConfig returns config read from file or error
func readConfig(filePath string) (*Config, error) {
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

// fetchRobots tries to get 'robots.txt' file for the seed URL
// for simplicity, it is ok if 'robots.txt' could not be fetched
// see https://developers.google.com/search/reference/robots_txt#handling-http-result-codes
func fetchRobots(link string) *robotstxt.RobotsData {
	resp, err := http.Get(link)
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

func printResults(link string, links []string) {
	fmt.Printf("Links found on the page %s\n", link)
	for i := range links {
		fmt.Printf("\t%s\n", links[i])
	}
}
