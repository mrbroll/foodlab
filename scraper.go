package main

import (
	"net/http"
	"net/url"
)

// Scraper is a type for scraping a web page believed to contain
// a single recipe.
type Scraper struct {
	client GetHeader
}

// ScrapeResult is a type for storing the result of
// scraping a web page containing a recipe.
type ScrapeResult struct {
	OutboundLinks []*url.URL
	Recipe        *Recipe
}

// Getter is a client capable of making GET requests.
type Getter interface {
	Get(url string) (resp *http.Response, err error)
}

// Header is a client capable of making HEAD requests.
type Header interface {
	Head(url string) (resp *http.Response, err error)
}

// GetHeader is a client capable of making both GET and HEAD requests.
type GetHeader interface {
	Getter
	Header
}

// New returns a Scraper with the given client.
func New(client GetHeader) *Scraper {
	return &Scraper{client: client}
}

// Scrape loads the page for the given URL and returns a scrape result
// if successful, or an error otherwise.
func (s *Scraper) Scrape(url *url.URL) (*ScrapeResult, error) {
	return nil, nil
}
