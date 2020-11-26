package page_parser

import (
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ParsedPage struct {
	// List of collected links
	Links []string
	// Canonical URL: <link rel="canonical" href="...">
	CanonicalURL string
	// Base URL: <base href="...">
	BaseURL string
}

func Parse(body io.Reader) (*ParsedPage, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	var page ParsedPage
	doc.Find("a[href]").Each(func(i int, selection *goquery.Selection) {
		if href, ok := selection.Attr("href"); ok && href != "" {
			page.Links = append(page.Links, strings.TrimSpace(href))
		}
	})
	if cURL, ok := doc.Find(`link[rel=canonical]`).First().Attr("href"); ok {
		page.CanonicalURL = strings.TrimSpace(cURL)
	}
	if bURL, ok := doc.Find(`base[href]`).First().Attr("href"); ok {
		page.BaseURL = strings.TrimSpace(bURL)
		if bu, err := url.Parse(page.BaseURL); err == nil { // ignore unparsable URL
			for i := range page.Links {
				if lu, err := url.Parse(page.Links[i]); err == nil { // don't handle bad links
					page.Links[i] = bu.ResolveReference(lu).String()
				}
			}
		}
	}
	return &page, nil
}
