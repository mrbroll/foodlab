package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mrbroll/foodlab/ndb"
)

var (
	query    string
	pageSize int
	token    string
)

func main() {
	input := bufio.NewReader(os.Stdin)
	flagSet := flag.NewFlagSet("NDB Search Flags", flag.ExitOnError)
	flagSet.IntVar(&pageSize, "page-size", 10, "Maximum page size for results.")
	flagSet.StringVar(&query, "query", "", "Search query. (Required)")
	flagSet.StringVar(&token, "token", "", "NDB API Token. (Required)")
	flagSet.Parse(os.Args[1:])
	if token == "" {
		fmt.Println("Must provide a valid token.")
		os.Exit(2)
	}
	if query == "" {
		fmt.Println("Must provide a valid query.")
		os.Exit(2)
	}
	client := ndb.NewHTTPClient(new(http.Client), token)
	iter := client.FoodSearch(query)
ITERATE:
	for food := iter.Next(); food != nil; food = iter.Next() {
		fBytes, err := json.Marshal(food)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", fBytes)
	}

	err := iter.Err()

	if err != nil && err != io.EOF {
		fmt.Println(err)
		os.Exit(1)
	}

	if err != nil && err == io.EOF {
		fmt.Println("No more results.")
		fmt.Println("Exiting...")
		os.Exit(0)
	}

ASK:
	fmt.Println("Show more results? (y/n)")

	answer, err := input.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if answer == "n" {
		fmt.Println("Exiting...")
		os.Exit(0)
	} else if answer == "y" {
		goto ITERATE
	} else {
		fmt.Println("Unrecognized input.")
		goto ASK
	}
}
