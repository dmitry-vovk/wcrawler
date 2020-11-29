<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-98%25-brightgreen.svg?longCache=true&style=flat)</a>

# wcrawler

Simple web crawler exercise.

## Requirements

Given a starting URL, the crawler should visit each URL it finds on the same domain.
It should print each URL visited, and a list of links found on that page. 
The crawler should be limited to one subdomain -- so when you start with `https://example.com/`,
it would crawl all pages within `example.com`, but not follow external links,
for example to `facebook.com` or `community.example.com`.

## Extended requirements

* Scalability: the crawler shall work with single domain and process up to a few thousand pages. 
* Robustness: the crawler shall tolerate errors when fetching pages, handle slow responses by setting a time out threshold.
* Politeness: the crawler shall obey `robots.txt` rules by not visiting pages that are not allowed to visit.
* Extensibility: it should be fairly easy to add new functionality, such as:
  * Additional content processing, e.g. extracting and saving images.
  * Downloading certain file types, e.g. PDF documents. 
