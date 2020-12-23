<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-97%25-brightgreen.svg?longCache=true&style=flat)</a>

# wcrawler

Simple web crawler exercise.

## Requirements

Given a starting URL, the crawler should visit each URL it finds on the same domain.
It should print each URL visited, and a list of links found on that page. 
The crawler should be limited to one subdomain -- so when you start with `https://example.com/`,
it would crawl all pages within `example.com`, but not follow external links,
for example to `facebook.com` or `community.example.com`.

## Extended requirements

* Scalability: the crawler works with single domain and processes relative small number of pages.
* Robustness: the crawler tolerates errors when fetching pages, handle slow responses.
* Politeness: the crawler obeys `robots.txt` rules by not visiting pages that are not allowed to visit.
* Extensibility: it should be fairly easy to add new functionality or alter existing, such as:
  * Additional content processing, e.g. extracting and saving images.
  * Downloading certain file types, e.g. PDF documents.
* The crawler is able to do HEAD requests before downloading a URL to determine if the URL returns HTML, 
  as we do not want to download files of different types. Failure with HEAD request shall not prevent it to do GET requests.
* The crawler handles cookies to avoid certain scenarios when a website sets tracking cookie with redirection,
  which may cause loops and other undesired behaviour.
* The crawler sends `Referer` header to counter "anti-hotlinking" measures.

## Project structure

Packages:
 * `crawler` -- base package containing top-level Crawler type and constructor to build it.
 * `crawler/page_fetcher` -- contains the code needed to perform HTTP requests and return fetched content.
 * `crawler/page_parser` -- contains the code to parse the page contents, extract links, and resolve them against base URL.
 * `types` -- contains types allowing testing `crawler` package.
 * `crawler/url_filter` -- contains the code that filters and normalises found URLs. 

## Configuration

Crawler configuration file has the following sytax:
```json
{
  "seed_url": "https://example.com",
  "ignore_robots_txt": false,
  "allow_www_prefix": true,
  "user_agent": "CrawlBot/0.1",
  "do_head_requests": true,
  "max_pages": 100,
  "max_parallel_requests": 5
}
```
Crawler will search for config file in this order:
1. Command line argument: `crawler config.json`
2. Environment variable: `CRAWLER_CONFIG=config.json crawler`
3. A file named `config.json` in the current PATH: `crawler`

## Building

`go build crawler.go`

## Testing

`go test -race ./...`

## Running

Crawler outputs results into stdOut, logs go into stdErr.
To collect results into a text file, the following command will do:
`crawler > results.txt`.
