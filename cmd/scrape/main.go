package main

import (
	"encoding/json"
	"fmt"
	"github.com/mrbroll/vine/scraper"
	"net/http"
)

func main() {
	url := "https://www.allrecipes.com/recipe/165190/spicy-vegan-potato-curry/"

	sc := scraper.GetScraper(url)

	res, err := sc.Scrape(url)
	if err != nil {
		panic(err)
	}

	if res == nil {
		fmt.Println("nil response")
	} else {
		recipeBytes, err := json.Marshal(res.Recipe)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", recipeBytes)
	}
}
