package allrecipes

import (
	"fmt"
	"strings"

	"github.com/mrbroll/vine"
	"github.com/mrbroll/vine/scraper"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type Scraper struct {
	getter Getter
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

	links := findLinks(node)
	title := findTitle(node)
	ingredients := findIngredients(node)
	instructions := findInstructions(node)

	return &ScrapeResult{
		Links: links,
		Recipe: &vine.Recipe{
			URL:          url,
			Title:        title,
			Ingredients:  ingredients,
			Instructions: instructions,
		},
	}, nil
}

func findLinks(node *html.Node) (links []string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" && strings.HasPrefix(attr.Val, "http") {
				links = append(links, attr.Val)
			}
		}
	}

	for node = node.FirstChild; node != nil; node = node.NextSibling {
		links = append(links, findLinks(node)...)
	}
	return links
}

func findTitle(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "title" {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				return child.Data
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		t := findTitle(child)
		if t != "" {
			return t
		}
	}
	return ""
}

func findIngredients(node *html.Node) []*vine.Ingredient {
	if node.Type == html.ElementNode && node.Data == "ul" {
		for _, attr := range node.Attr {
			fmt.Println(attr.Val)
			if ingredientsRegex.MatchString(attr.Val) {
				ings := []*vine.Ingredient{}
				lis := findListItems(node)
				for _, li := range lis {
					ing := parseIngredientText(li)
					if ing != nil {
						ings = append(ings, ing)
					}
				}
				return ings
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		ings := findIngredients(child)
		if len(ings) > 0 {
			return ings
		}
	}

	return nil
}

func findInstructions(node *html.Node) (insts []*vine.Instruction) {
	if node.Type == html.ElementNode && node.Data == "ol" {
	}
	return insts
}

func parseIngredientText(t string) *vine.Ingredient {
	// TODO: actually parse
	return &vine.Ingredient{
		Name: t,
	}
}

func findListItems(node *html.Node) []string {
	if node.Type == html.ElementNode && node.Data == "li" {
		liStrs := []string{findInnerText(node)}
		for node = node.NextSibling; node != nil; node = node.NextSibling {
			liStrs = append(liStrs, findInnerText(node))
		}
		return liStrs
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		lis := findListItems(child)
		if lis != nil {
			return lis
		}
	}
	return nil
}

func findInnerText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		t := findInnerText(child)
		if t != "" {
			return t
		}
	}

	return ""
}
