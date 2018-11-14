// pacakge ndb contains functionality for interacting with the ndb database: api.nal.usda.gov/ndb
package ndb

// FoodReportResponse is a type for NDB food report API responses.
type FoodReportResponse struct {
	Report *FoodReport `json:"report"`
}

// FoodSearchResponse is a type for NDB food search API responses.
type FoodSearchResponse struct {
	Results *FoodSearchResults `json:"list"`
}

// FoodSearchResults is a type for NDB food search API results.
type FoodSearchResults struct {
	Foods []*Food `json:"item"`
}

// FoodReport is a type for NDB food reports.
type FoodReport struct {
	Food *Food `json:"food"`
}

// Food is a type for NDB foods.
type Food struct {
	NDBID            string              `json:"ndbno"`
	Name             string              `json:"name"`
	DataSourceAbbrev string              `json:"ds"`
	Group            string              `json:"group,omitempty"`
	Nutrients        []*NutrientQuantity `json:"nutrients,omitempty"`
}

// Nutrient is a type for NDB nutrients.
type Nutrient struct {
	NDBID    string             `json:"nutrient_id"`
	Name     string             `json:"name"`
	Group    string             `json:"group"`
	Unit     string             `json:"unit"`
	Value    float64            `json:"value,string"`
	Measures []*NutrientMeasure `json:"measures"`
}

// NutrientQuantity is a type for NDB nutrient quanitities.
type NutrientQuantity struct {
	Nutrient
	Unit  string  `json:"unit"`
	Value float64 `json:"value,string"`
}

// NutrientMeasure is a type for NDB nutrient measurements.
type NutrientMeasure struct {
	Label       string  `json:"label"`
	EqvQuantity float64 `json:"eqv"`
	EqvUnit     string  `json:"eunit"`
	Quantity    float64 `json:"qty"`
	Value       float64 `json:"value,string"`
}
