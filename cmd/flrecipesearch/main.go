package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mrbroll/foodlab/recipe"
)

const (
	dgraphHost string = "localhost:9080"
)

var (
	query string
)

func main() {
	flagSet := flag.NewFlagSet("Recipe Search Flags", flag.ExitOnError)
	flagSet.StringVar(&query, "query", "", "Search query. (Required)")

	flagSet.Parse(os.Args[1:])
	if query == "" {
		fmt.Println("Must provide a non-empty query.")
		os.Exit(2)
	}

	store := recipe.NewDgraphStore(dgraphHost)

	recipes, err := store.SearchRecipe(query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, recipe := range recipes {
		rNutr := recipe.AggregateNutrition()
		rnb, err := json.Marshal(rNutr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", rnb)
	}
}
