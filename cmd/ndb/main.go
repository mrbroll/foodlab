package main

import (
	"encoding/json"
	"fmt"
	"github.com/mrbroll/vine/ndb"
	"net/http"
)

const (
	apiToken string = "O7G3lWJlaYkMAnCeaFFhd7rM0wmdmTR2xkJmOslZ"
)

func main() {
	httpClient := new(http.Client)
	ndbClient := ndb.NewHTTPClient(httpClient, apiToken)

	foods, err := ndbClient.FoodSearch("potato")
	if err != nil {
		panic(err)
	}

	foodsJSON, err := json.Marshal(foods)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", foodsJSON)

	/*
		foodID := foods[0].NDBID
		food, err := ndbClient.FoodReport(foodID)
		if err != nil {
			panic(err)
		}

		foodJSON, err := json.Marshal(food)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", foodJSON)
	*/
}
