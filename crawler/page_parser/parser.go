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

func (p *ParsedPage) addLink(link string) {
	if href := strings.TrimSpace(link); acceptableLink(href) {
		p.Links = append(p.Links, href)
	}
}

func (p *ParsedPage) resolveLinks() {
	if p.BaseURL == "" {
		return
	}
	if bu, err := url.Parse(p.BaseURL); err == nil { // ignore unparsable URL
		for i := range p.Links {
			if lu, err := url.Parse(p.Links[i]); err == nil { // don't handle bad links
				p.Links[i] = bu.ResolveReference(lu).String()
			}
		}
	}
}

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
		page.BaseURL = strings.TrimSpace(href)
		page.resolveLinks()
	}
	return &page, nil
}

func acceptableLink(href string) bool {
	if href == "" {
		return false
	}
	if strings.HasPrefix(href, "#") {
		return false
	}
	return true
}
