package recipe

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mrbroll/foodlab/ndb"
	"github.com/pkg/errors"
)

type CLIComposer struct {
	sr    FoodSearcherReporter
	store RecipeAdder
}

type FoodSearcher interface {
	FoodSearch(query string) ([]*ndb.Food, error)
}

type FoodReporter interface {
	FoodReport(ndbno string) (*ndb.Food, error)
}

type FoodSearcherReporter interface {
	FoodSearcher
	FoodReporter
}

func NewCLIComposer(sr FoodSearcherReporter, store RecipeAdder) *CLIComposer {
	return &CLIComposer{
		sr:    sr,
		store: store,
	}
}

func (c *CLIComposer) Compose() error {
	reader := bufio.NewReader(os.Stdin)
	// add name
	fmt.Println("Enter a name for your recipe: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return errors.Wrap(err, "Reading recipe name.")
	}
	name = strings.TrimSpace(name)
	// add ingredients
	fmt.Println("Please add ingredients:")
	ingredients := []*Ingredient{}
	for addIngredient := true; addIngredient; {
		ingredient := new(Ingredient)
		fmt.Println("Ingredient Name:")
		keywords, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "Reading ingredient name.")
		}
		keywords = strings.TrimSpace(keywords)
		// search for suggestions
		foods, err := c.sr.FoodSearch(keywords)
		if err != nil {
			return errors.Wrap(err, "Searching for food.")
		}
		foodMap := map[string]*ndb.Food{}
		for _, food := range foods {
			foodMap[food.NDBID] = food
			fmt.Printf("%s: %s\n", food.NDBID, food.Name)
		}
		fmt.Println("Please choose one of the suggested foods above by entering its id:")
	SELECT_SUGGESTION:
		foodID, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "Reading selected food id.")
		}
		foodID = strings.TrimSpace(foodID)
		food, ok := foodMap[foodID]
		if !ok {
			fmt.Printf("ID %s not found, please select a valid food ID:\n", foodID)
			goto SELECT_SUGGESTION
		}
		ingredient.Name = food.Name
		// get food report

		foodReport, err := c.sr.FoodReport(food.NDBID)
		if err != nil {
			return errors.Wrapf(err, "Getting food report for ndb id %s.", food.NDBID)
		}
		foodJSON, _ := json.MarshalIndent(foodReport, "", "  ")
		fmt.Printf("%s\n", foodJSON)
	MEASURE:
		meas := new(Measurement)
		fmt.Println("Ingredient Measurement (<number> [<unit>]:")
		m, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "Reading ingredient measurement.")
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
			unit, ok := UnitAliases[u]
			if !ok {
				fmt.Printf("Unrecognized unit %s\n", u)
				goto MEASURE
			}
			meas.Unit = unit
		} else {
			// get quantity for nutrients

		}
		ingredient.Measurement = meas
		ingredients = append(ingredients, ingredient)

		fmt.Println("Would you like to add another ingredient? (y/n):")
		a, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "Reading next instruction answer.")
		}
		if strings.TrimSpace(a) == "n" {
			addIngredient = false
		}
	}

	// add instructions
	instructions := []*Instruction{}
	for iNo, addInstruction := 1, true; addInstruction; iNo++ {
		instruction := new(Instruction)
		fmt.Println("Enter Instruction:")
		txt, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "Reading instruction.")
		}
		instruction.Order = iNo
		instruction.Text = strings.TrimSpace(txt)
		instructions = append(instructions, instruction)
		fmt.Println("Would you like to add another instruction? (y/n):")
		a, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "Reading next instruction answer.")
		}
		if strings.TrimSpace(a) == "n" {
			addInstruction = false
		}
	}

	// store recipe
	if err := c.store.AddRecipe(&Recipe{
		Name:         name,
		Ingredients:  ingredients,
		Instructions: instructions,
	}); err != nil {
		return errors.Wrap(err, "Adding recipe to store")
	}
	return nil
}
