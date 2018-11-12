// pacakge ndb contains functionality for interacting with the ndb database: api.nal.usda.gov/ndb
package ndb

type FoodReportResponse struct {
	Report *FoodReport `json:"report"`
}

type FoodSearchResponse struct {
	Results *FoodSearchResults `json:"list"`
}

type FoodSearchResults struct {
	Foods []*Food `json:"item"`
}

type FoodReport struct {
	Food *Food `json:"food"`
}

type Food struct {
	NDBID            string              `json:"ndbno"`
	Name             string              `json:"name"`
	DataSourceAbbrev string              `json:"ds"`
	Group            string              `json:"group,omitempty"`
	Nutrients        []*NutrientQuantity `json:"nutrients,omitempty"`
}

type FoodGroup struct {
	NDBID string `json:"ndbno"`
	Name  string `json:"name"`
}

type Nutrient struct {
	NDBID    string             `json:"nutrient_id"`
	Name     string             `json:"name"`
	Group    string             `json:"group"`
	Unit     string             `json:"unit"`
	Value    float64            `json:"value,string"`
	Measures []*NutrientMeasure `json:"measures"`
}

type NutrientQuantity struct {
	Nutrient
	Unit  string  `json:"unit"`
	Value float64 `json:"value,string"`
}

type NutrientMeasure struct {
	Label       string  `json:"label"`
	EqvQuantity float64 `json:"eqv"`
	EqvUnit     string  `json:"eunit"`
	Quantity    float64 `json:"qty"`
	Value       float64 `json:"value,string"`
}

type DataSource struct {
	Name   string
	Abbrev string
}
