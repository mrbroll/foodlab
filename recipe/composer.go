package recipe

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mrbroll/foodlab/ndb"
	"github.com/pkg/errors"
)

// CLIComposer is a type for composing a recipe via CLI.
type CLIComposer struct {
	ndbStore    NDBStore
	recipeStore RecipeStore
}

// NDBStore is a type for interacting with the NDB.
type NDBStore interface {
	// FoodSearch returns an iterator for the given query.
	FoodSearch(query string) *ndb.FoodIter

	// FoodReport returns an NDB food report for the food with the given ndbno.
	// It returns an error if the NDB was unreachable.
	FoodReport(ndbno string) (*ndb.Food, error)
}

// RecipeStore is a type for interacting with Recipes, Ingredients, and Nutritional information.
type RecipeStore interface {
	// Add recipe adds a recipe to the store.
	// It returns an error if the operation was unsuccessful.
	AddRecipe(r *Recipe) error

	// GetOrCreateFood idempotently creates the food in the store,
	// and returns it with the UID attribute populated.
	// It returns an error if the operation was unsuccessful.
	GetOrCreateFood(food *Food) (*Food, error)
}

// NewCLIComposer returns a cli composer using the given interfaces
// for food searches and reports, as well as recipe storage.
func NewCLIComposer(ndb NDBStore, rs RecipeStore) *CLIComposer {
	return &CLIComposer{
		ndbStore:    ndb,
		recipeStore: rs,
	}
}

// Compose kicks of a cli composition session.
// It returns an error if the composition failed.
func (c *CLIComposer) Compose() error {
	name, err := c.getName(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "Naming recipe")
	}

	ingredients, err := c.getIngredients(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "Adding ingredients to recipe.")
	}

	instructions, err := c.getInstructions(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "Adding instructions to recipe.")
	}

	// store recipe
	if err := c.recipeStore.AddRecipe(&Recipe{
		Name:         name,
		Ingredients:  ingredients,
		Instructions: instructions,
	}); err != nil {
		return errors.Wrap(err, "Adding recipe to store")
	}
	return nil
}

// getName gets a name using input from the given reader.
// It returns an error if unsuccessful.
func (c *CLIComposer) getName(in io.Reader) (string, error) {
	reader := bufio.NewReader(in)
	fmt.Println("Enter a name for your recipe: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.Wrap(err, "Reading recipe name.")
	}
	return strings.TrimSpace(name), nil
}

// getIngredients get ingredients  using input from the given reader.
// It returns an error if unsuccessful.
func (c *CLIComposer) getIngredients(in io.Reader) ([]*Ingredient, error) {
	input := bufio.NewReader(in)
	fmt.Println("Please add ingredients:")
	ingredients := []*Ingredient{}
	for {
		ingredient := new(Ingredient)
		fmt.Println("Ingredient Name:")
		keywords, err := input.ReadString('\n')
		if err != nil {
			return nil, errors.Wrap(err, "Reading ingredient name.")
		}
		keywords = strings.TrimSpace(keywords)
		// search for suggestions
		foodIter := c.ndbStore.FoodSearch(keywords)
		for ndbFood := foodIter.Next(); ndbFood != nil; ndbFood = foodIter.Next() {
		FOOD_SUGGESTION:
			fmt.Printf("Did you mean \"%s\" (y/n)?\n:", ndbFood.Name)
			answer, err := input.ReadString('\n')
			if err != nil {
				return nil, errors.Wrap(err, "Reading answer for food suggestion.")
			}
			answer = strings.ToLower(strings.TrimSpace(answer))
			if answer == "n" {
				continue
			} else if answer == "y" {
				ndbFood, err := c.ndbStore.FoodReport(ndbFood.NDBID)
				if err != nil {
					return nil, errors.Wrap(err, "Getting NDB Food Report.")
				}
			GET_UNIT:
				// get unit of measure
				measures := ndbFood.Nutrients[0].Measures
				fmt.Printf("Please select a unit of measure (%d-%d):\n", 0, len(measures)-1)
				for i, meas := range measures {
					fmt.Printf("%d: %s\n", i, meas.Label)
				}
				measIdxStr, err := input.ReadString('\n')
				if err != nil {
					return nil, errors.Wrap(err, "Getting unit of measure.")
				}
				measIdx, err := strconv.Atoi(strings.TrimSpace(measIdxStr))
				if err != nil {
					fmt.Println(err)
					goto GET_UNIT
				} else if measIdx < 0 || measIdx >= len(measures) {
					fmt.Println("Not a valid option.")
					goto GET_UNIT
				}
				ingredient.Unit = measures[measIdx].Label
			GET_QUANTITY:
				// get measured quantity
				fmt.Println("Enter a quantity as a decimal:")
				qStr, err := input.ReadString('\n')
				if err != nil {
					return nil, errors.Wrap(err, "Getting quantity.")
				}
				quantity, err := strconv.ParseFloat(strings.TrimSpace(qStr), 64)
				if err != nil {
					fmt.Println(err)
					goto GET_QUANTITY
				}
				ingredient.Value = quantity

				food, err := c.recipeStore.GetOrCreateFood(NewFoodFromNDB(ndbFood))
				if err != nil {
					return nil, errors.Wrap(err, "Getting or creating food.")
				}
				ingredient.Food = food
				break
			} else {
				fmt.Printf("Invalid input: \"%s\"\n", answer)
				goto FOOD_SUGGESTION
			}
		}
		ingredient.UID = fmt.Sprintf("_:%s", ingredient.Hash())
		ingredients = append(ingredients, ingredient)

	ANOTHER_INGREDIENT:
		fmt.Println("Add another ingredient? (y/n)")
		a, err := input.ReadString('\n')
		if err != nil {
			return nil, errors.Wrap(err, "Reading next ingredient answer.")
		}
		a = strings.ToLower(strings.TrimSpace(a))
		if a == "n" {
			break
		} else if a != "y" {
			fmt.Printf("Invalid input: \"%s\".\n", a)
			goto ANOTHER_INGREDIENT
		}
	}

	return ingredients, nil
}

// getInstructions gets instructions from the given reader.
// It returns an error if unsuccessful.
func (c *CLIComposer) getInstructions(in io.Reader) ([]*Instruction, error) {
	reader := bufio.NewReader(in)
	instructions := []*Instruction{}
	for iNo, addInstruction := 1, true; addInstruction; iNo++ {
		instruction := new(Instruction)
		fmt.Println("Enter Instruction:")
		txt, err := reader.ReadString('\n')
		if err != nil {
			return nil, errors.Wrap(err, "Reading instruction.")
		}
		instruction.Order = iNo
		instruction.Text = strings.TrimSpace(txt)
		instructions = append(instructions, instruction)
		fmt.Println("Would you like to add another instruction? (y/n):")
		a, err := reader.ReadString('\n')
		if err != nil {
			return nil, errors.Wrap(err, "Reading next instruction answer.")
		}
		if strings.TrimSpace(a) == "n" {
			addInstruction = false
		}
	}
	return instructions, nil
}
