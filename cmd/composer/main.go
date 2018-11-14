package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mrbroll/vine/recipe"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter a name for your recipe: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Please add ingredients:")
	ingredients := []*recipe.Ingredient{}
	for addIngredient := true; addIngredient; {
		ingredient := new(recipe.Ingredient)
		fmt.Println("Ingredient Name:")
		name, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		ingredient.Name = strings.TrimSpace(name)
	MEASURE:
		meas := new(recipe.Measurement)
		fmt.Println("Ingredient Measurement (<number> [<unit>]:")
		m, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		mParts := strings.Split(strings.TrimSpace(m), " ")
		if len(mParts) >= 1 {
			q, err := strconv.ParseFloat(mParts[0], 64)
			if err != nil {
				fmt.Println("Please enter a valid number.")
				goto MEASURE
			}
			meas.Value = q
		}
		if len(mParts) == 2 {
			u := strings.TrimSpace(strings.ToLower(mParts[1]))
			unit, ok := recipe.UnitAliases[u]
			if !ok {
				fmt.Printf("Unrecognized unit %s\n", u)
				goto MEASURE
			}
			meas.Unit = unit
		}
		ingredient.Measurement = meas
		ingredients = append(ingredients, ingredient)

		fmt.Println("Would you like to add another ingredient? (y/n):")
		a, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if strings.TrimSpace(a) == "n" {
			addIngredient = false
		}
	}

	instructions := []*recipe.Instruction{}
	for iNo, addInstruction := 0, true; addInstruction; iNo++ {
		instruction := new(recipe.Instruction)
		fmt.Println("Enter Instruction:")
		txt, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		instruction.Order = iNo
		instruction.Text = strings.TrimSpace(txt)
		instructions = append(instructions, instruction)
		fmt.Println("Would you like to add another instruction? (y/n):")
		a, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if strings.TrimSpace(a) == "n" {
			addInstruction = false
		}
	}

	recipe := &recipe.Recipe{
		Name:         strings.TrimSpace(name),
		Ingredients:  ingredients,
		Instructions: instructions,
	}

	rBytes, err := json.MarshalIndent(recipe, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", rBytes)
}
