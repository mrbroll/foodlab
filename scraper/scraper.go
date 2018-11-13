package scraper

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/mrbroll/vine"
	"github.com/mrbroll/vine/scraper/allrecipes"
	"github.com/pkg/errors"
)

var (
	ingredientsRegex  *regexp.Regexp
	instructionsRegex *regexp.Regexp

	scrapers = map[string]Scraper{
		"allrecipes.com": *allrecipes.Scraper,
	}
)

func init() {
	ingredientsRegex = regexp.MustCompile("^.*ingredient.*$/i")
	instructionsRegex = regexp.MustCompile("^.*(instruction|direction).*$/i")
}

// Scraper is a type for scraping a web page believed to contain
// a single recipe.
type Scraper interface {
	Scrape(url string) (*ScrapeResult, error)
}

// ScrapeResult is a type for storing the result of
// scraping a web page containing a recipe.
type ScrapeResult struct {
	Links  []string     `json:"links"`
	Recipe *vine.Recipe `json:"recipe"`
}

// Getter is a client capable of making GET requests.
type Getter interface {
	Get(url string) (resp *http.Response, err error)
}

// GetScraper returns a custom scraper for the given url.
// It returns an error if the given url cannot be parsed,
// or no scraper can be found.
func GetScraper(s string) (Scraper, error) {
	parsedURL, err := url.Parse(s)
	if err != nil {
		return nil, errors.Wrap(err, "Parsing URL.")
	}
	s, ok := scrapers[parsedURL.Host]
	if !ok {
		return nil, errors.Errorf("No scraper for %s.", parsedURL.Host)
	}
	return s, nil
}
