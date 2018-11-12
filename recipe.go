package vine

// Recipe represents a recipe.
type Recipe struct {
	Title        string         `json:"title"`
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
	Units MeasuredUnits
}

// Units describes the units used for a quantity of an ingredient.
type MeasuredUnits interface {
	isUnits()
}

type Count struct{}

type Centimeters struct{}
type Millimeters struct{}
type Inches struct{}

type Grams struct{}
type Kilograms struct{}
type Milligrams struct{}

type Ounces struct{}
type Pounds struct{}

type Cups struct{}
type Liters struct{}
type Milliliters struct{}
type Tablespoons struct{}
type Teaspoons struct{}

func (_ Count) isUnits() {}

func (_ Centimeters) isUnits() {}
func (_ Millimeters) isUnits() {}
func (_ Inches) isUnits()      {}

func (_ Grams) isUnits()      {}
func (_ Kilograms) isUnits()  {}
func (_ Milligrams) isUnits() {}

func (_ Ounces) isUnits() {}
func (_ Pounds) isUnits() {}

func (_ Cups) isUnits()        {}
func (_ Liters) isUnits()      {}
func (_ Milliliters) isUnits() {}
func (_ Tablespoons) isUnits() {}
func (_ Teaspoons) isUnits()   {}
