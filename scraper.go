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
	Links  []string `json:"links"`
	Recipe *Recipe  `json:"recipe"`
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

	bfs([]*html.Node{node})

	links, err := findLinks(node)
	if err != nil {
		return nil, errors.Wrap(err, "Finding links.")
	}

	title, err := findTitle(node)
	if err != nil {
		return nil, errors.Wrap(err, "Finding title.")
	}

	ingredients, err := findIngredients(node)
	if err != nil {
		return nil, errors.Wrap(err, "Finding ingredients.")
	}

	instructions, err := findInstructions(node)
	if err != nil {
		return nil, errors.Wrap(err, "Finding instructions.")
	}

	return &ScrapeResult{
		Links: links,
		Recipe: &Recipe{
			Title:        title,
			Ingredients:  ingredients,
			Instructions: instructions,
		},
	}, nil
}

func findLinks(node *html.Node) (links []string, err error) {
	return links, err
}

func findTitle(node *html.Node) (title string, err error) {
	return title, err
}

func findIngredients(node *html.Node) (ingredients []*Ingredient, err error) {
	return ingredients, err
}

func findInstructions(node *html.Node) (instructions []*Instruction, err error) {
	return instructions, err
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
