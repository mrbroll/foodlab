package vine

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"net/http"
)

// Scraper is a type for scraping a web page believed to contain
// a single recipe.
type Scraper struct {
	getter Getter
}

// ScrapeResult is a type for storing the result of
// scraping a web page containing a recipe.
type ScrapeResult struct {
	OutboundLinks []string `json:"outboundLinks"`
	Recipe        *Recipe  `json:"recipe"`
}

// Getter is a client capable of making GET requests.
type Getter interface {
	Get(url string) (resp *http.Response, err error)
}

// NewScraper returns a Scraper with the given client.
func NewScraper(getter Getter) *Scraper {
	return &Scraper{
		getter: getter,
	}
}

// Scrape loads the page for the given URL and returns a scrape result
// if successful, or an error otherwise.
func (s *Scraper) Scrape(url string) (*ScrapeResult, error) {
	resp, err := s.getter.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "Issuing GET request to %s.", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("Status code %d from url %s.", resp.StatusCode, url)
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Parsing HTML.")
	}
	res := new(ScrapeResult)

	bfs([]*html.Node{node})

	return res, nil
}

func bfs(nodes []*html.Node) {
	if len(nodes) == 0 {
		return
	}
	cNodes := []*html.Node{}
	for _, node := range nodes {
		fmt.Println(node.Data)
		for n := node.FirstChild; n != nil; n = n.NextSibling {
			cNodes = append(cNodes, n)
		}
	}

	bfs(cNodes)
}
