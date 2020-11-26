# wcrawler

Simple web crawler exercise.

Given a starting URL, the crawler should visit each URL it finds on the same domain.
It should print each URL visited, and a list of links found on that page. 
The crawler should be limited to one subdomain -- so when you start with `https://example.com/`,
it would crawl all pages within `example.com`, but not follow external links,
for example to `facebook.com` or `community.example.com`.
