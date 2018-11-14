package recipe

const (
	Count      MeasuredUnit = ""
	Inch       MeasuredUnit = "in"
	Cup        MeasuredUnit = "c"
	Tablespoon MeasuredUnit = "tbsp"
	Teaspoon   MeasuredUnit = "tsp"
)

var (
	UnitAliases = map[string]MeasuredUnit{
		"":            Count,
		"count":       Count,
		"c":           Cup,
		"cup":         Cup,
		"cups":        Cup,
		"tbsp":        Tablespoon,
		"tbsps":       Tablespoon,
		"tablespoon":  Tablespoon,
		"tablespoons": Tablespoon,
		"tsp":         Teaspoon,
		"tsps":        Teaspoon,
		"teaspoon":    Teaspoon,
		"teaspoons":   Teaspoon,
	}
)

// Recipe represents a recipe.
type Recipe struct {
	Name         string         `json:"name"`
	Ingredients  []*Ingredient  `json:"ingredients"`
	Instructions []*Instruction `json:"instructions"`
}

// Ingredient is a single ingredient of a recipe.
type Ingredient struct {
	Name        string       `json:"name"`
	Measurement *Measurement `json:"measurement"`
	Preparation string       `json:"preparation"`
}

// Instruction is a single instruction for preparing a recipe.
type Instruction struct {
	Order int
	Text  string
}

// Measurement describes the quantity of an ingredient in a recipe.
type Measurement struct {
	Value float64
	Unit  MeasuredUnit
}

type MeasuredUnit string
