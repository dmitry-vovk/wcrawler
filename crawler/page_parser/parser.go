package page_parser

import (
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParsedPage represents contents of a web page
type ParsedPage struct {
	// List of collected links
	Links []string
	// Canonical URL: <link rel="canonical" href="...">
	CanonicalURL string
	// Base URL: <base href="...">
	baseURL string
}

func (p *ParsedPage) addLink(link string) {
	if href := strings.TrimSpace(link); acceptableLink(href) {
		p.Links = append(p.Links, href)
	}
}

// resolveLinks tries to convert relative links into absolute ones
func (p *ParsedPage) resolveLinks() {
	if p.baseURL == "" {
		return
	}
	if bu, err := url.Parse(p.baseURL); err == nil { // we ignore unparseable base URL
		links := make([]string, 0, len(p.Links))
		for i := range p.Links {
			if lu, err := url.Parse(p.Links[i]); err == nil { // don't handle bad links
				links = append(links, bu.ResolveReference(lu).String())
			}
		}
		p.Links = links
	}
}

// Parse reads data from provided reader, extracting data
func Parse(body io.Reader) (*ParsedPage, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	var page ParsedPage
	doc.Find("a[href]").Each(func(i int, selection *goquery.Selection) {
		if href, ok := selection.Attr("href"); ok {
			page.addLink(href)
		}
	})
	if href, ok := doc.Find(`link[rel=canonical]`).First().Attr("href"); ok {
		page.CanonicalURL = strings.TrimSpace(href)
	}
	if href, ok := doc.Find(`base[href]`).First().Attr("href"); ok {
		page.baseURL = strings.TrimSpace(href)
		page.resolveLinks()
	}
	return &page, nil
}

// acceptableLink tells if the link is ok to use
func acceptableLink(href string) bool {
	if href == "" {
		return false
	}
	// 'fragments' are not ok
	if strings.HasPrefix(href, "#") {
		return false
	}
	return true
}
