package main

import (
	"encoding/json"
	"fmt"
	"github.com/mrbroll/vine"
	"net/http"
)

func main() {
	url := "https://www.allrecipes.com/recipe/165190/spicy-vegan-potato-curry/"

	scraper := vine.NewScraper(new(http.Client))

	res, err := scraper.Scrape(url)
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
