// package recipe contains types and functionality for creating, finding, and managing recipes.
package recipe

import (
	"encoding/json"
	"fmt"

	"github.com/mrbroll/foodlab/ndb"
	"github.com/pkg/errors"
)

// Recipe represents a recipe.
type Recipe struct {
	UID          string         `json:"uid,omitempty"`
	Name         string         `json:"name"`
	Ingredients  []*Ingredient  `json:"ingredient"`
	Instructions []*Instruction `json:"instruction"`
}

type RecipeNutrition struct {
	Name      string                       `json:"name"`
	Nutrition []*SimpleNutrientMeasurement `json:"nutrition"`
}

type SimpleNutrientMeasurement struct {
	Name  string  `json:"name"`
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

func (r *Recipe) AggregateNutrition() *RecipeNutrition {
	nutrMap := map[string]*SimpleNutrientMeasurement{}
	for _, ing := range r.Ingredients {
		unit := ing.Unit
		food := ing.Food
		var foodMeasure *FoodMeasurement
		for _, fMeas := range food.Measurements {
			if fMeas.Unit == unit {
				foodMeasure = fMeas
				break
			}
		}
		mul := ing.Value / foodMeasure.Value
		for _, nMeas := range foodMeasure.NutrientMeasurements {
			snm, ok := nutrMap[nMeas.Nutrient.Name]
			if !ok {
				nutrMap[nMeas.Nutrient.Name] = &SimpleNutrientMeasurement{
					Name:  nMeas.Nutrient.Name,
					Unit:  nMeas.Unit,
					Value: nMeas.Value * mul,
				}
			} else {
				snm.Value += nMeas.Value * mul
			}
		}
	}
	simpleMeasures := []*SimpleNutrientMeasurement{}
	for _, snm := range nutrMap {
		simpleMeasures = append(simpleMeasures, snm)
	}
	return &RecipeNutrition{
		Name:      r.Name,
		Nutrition: simpleMeasures,
	}
}

// Ingredient is a single ingredient of a recipe.
type Ingredient struct {
	UID   string  `json:"uid,omitempty"`
	Food  *Food   `json:"food"`
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

func (ingr *Ingredient) UnmarshalJSON(data []byte) error {
	type Ingr struct {
		UID   string  `json:"uid,omitempty"`
		Food  []*Food `json:"food"`
		Unit  string  `json:"unit"`
		Value float64 `json:"value"`
	}

	tmpIngr := new(Ingr)
	if err := json.Unmarshal(data, tmpIngr); err != nil {
		return errors.Wrap(err, "Unmarshaling intermediate ingredient.")
	}
	if len(tmpIngr.Food) > 0 {
		ingr.Food = tmpIngr.Food[0]
	}
	ingr.UID = tmpIngr.UID
	ingr.Unit = tmpIngr.Unit
	ingr.Value = tmpIngr.Value
	return nil
}

// Food is a type for listing a food along with its nutritional value.
type Food struct {
	UID          string             `json:"uid,omitempty"`
	NDBID        string             `json:"ndb_id,omitempty"`
	Name         string             `json:"name"`
	Measurements []*FoodMeasurement `json:"measurement"`
}

// FoodMeasurement is a type for listing a specific measurement of a food
// and its nutritient measurements.
type FoodMeasurement struct {
	UID                  string                 `json:"uid,omitempty"`
	Unit                 string                 `json:"unit"`
	Value                float64                `json:"value"`
	EqUnit               string                 `json:"eq_unit"`
	EqValue              float64                `json:"eq_value"`
	NutrientMeasurements []*NutrientMeasurement `json:"nutrient_measurement"`
}

// NutrientMeasurement is a type for listing a specific measurement of a nutrient.
// This is used to link to a FoodMeasurement in order to aggregate nutrient measurents.
type NutrientMeasurement struct {
	UID      string    `json:"uid,omitempty"`
	Nutrient *Nutrient `json:"nutrient"`
	Unit     string    `json:"unit"`
	Value    float64   `json:"value"`
}

func (m *NutrientMeasurement) UnmarshalJSON(data []byte) error {
	type NMeas struct {
		UID      string      `json:"uid,omitempty"`
		Nutrient []*Nutrient `json:"nutrient"`
		Unit     string      `json:"unit"`
		Value    float64     `json:"value"`
	}
	nm := new(NMeas)
	if err := json.Unmarshal(data, nm); err != nil {
		return errors.Wrap(err, "Unmarshing nutrient measure json.")
	}
	nutrient := new(Nutrient)
	if len(nm.Nutrient) > 0 {
		nutrient = nm.Nutrient[0]
	}
	m.UID = nm.UID
	m.Nutrient = nutrient
	m.Unit = nm.Unit
	m.Value = nm.Value
	return nil
}

// Nutrient is used to globally identify nutrients in the database.
type Nutrient struct {
	UID      string `json:"uid,omitempty"`
	NDBID    string `json:"ndb_id,omitempty"`
	NDBGroup string `json:"ndb_group,omitempty"`
	Name     string `json:"name"`
}

// Instruction is a single instruction for preparing a recipe.
type Instruction struct {
	UID   string `json:"uid,omitempty"`
	Order int    `json:"order"`
	Text  string `json:"text"`
}

// NewFoodFromNDB creates a food from an ndb food.
func NewFoodFromNDB(ndbFood *ndb.Food) *Food {
	food := &Food{
		NDBID:        ndbFood.NDBID,
		Name:         ndbFood.Name,
		Measurements: []*FoodMeasurement{},
	}
	measMap := map[string][]*NutrientMeasurement{}
	for _, nMeas := range ndbFood.Nutrients[0].Measures {
		fMeas := &FoodMeasurement{
			Unit:                 nMeas.Label,
			Value:                nMeas.Quantity,
			EqUnit:               nMeas.EqUnit,
			EqValue:              nMeas.EqValue,
			NutrientMeasurements: []*NutrientMeasurement{},
		}
		food.Measurements = append(food.Measurements, fMeas)
		measMap[fMeas.Unit] = []*NutrientMeasurement{}
	}

	for _, nutr := range ndbFood.Nutrients {
		for _, nMeas := range nutr.Measures {
			measMap[nMeas.Label] = append(
				measMap[nMeas.Label],
				&NutrientMeasurement{
					Unit:  nutr.Unit,
					Value: nMeas.Value,
					Nutrient: &Nutrient{
						NDBID:    nutr.NDBID,
						NDBGroup: nutr.Group,
						Name:     nutr.Name,
					},
				},
			)
		}
	}

	for _, meas := range food.Measurements {
		meas.NutrientMeasurements = measMap[meas.Unit]
	}

	return food
}

// Hash returns a string to reference the node in mutation responses.
func (r *Recipe) Hash() string {
	return r.Name
}

// Hash returns a string to reference the node in mutation responses.
func (i *Ingredient) Hash() string {
	foodName := "_"
	if i.Food != nil {
		foodName = i.Food.Name
	}
	return fmt.Sprintf("%s_%f_%s", i.Unit, i.Value, foodName)
}

// Hash returns a string to reference the node in mutation responses.
func (f *Food) Hash() string {
	return f.Name
}

// Hash returns a string to reference the node in mutation responses.
func (m *FoodMeasurement) Hash() string {
	return fmt.Sprintf("%s_%f", m.Unit, m.Value)
}

// Hash returns a string to reference the node in mutation responses.
func (m *NutrientMeasurement) Hash() string {
	nutrient := "_"
	if m.Nutrient != nil {
		nutrient = m.Nutrient.Name
	}
	return fmt.Sprintf("%s_%f_%s", m.Unit, m.Value, nutrient)
}

// Hash returns a string to reference the node in mutation responses.
func (n *Nutrient) Hash() string {
	return n.Name
}

// Hash returns a string to reference the node in mutation responses.
func (i *Instruction) Hash() string {
	return fmt.Sprintf("%d_%s", i.Order, i.Text)
}
