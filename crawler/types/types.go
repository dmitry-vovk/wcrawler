/*
These types are required for testing
*/
package types

import "github.com/dmitry-vovk/wcrawler/crawler/page_fetcher"

type Fetcher interface {
	Fetch(*page_fetcher.Request) (*page_fetcher.Response, error)
}

type Filter interface {
	Filter(link string) (string, bool)
}
