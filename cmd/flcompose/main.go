package main

import (
	"log"
	"net/http"

	"github.com/mrbroll/foodlab/ndb"
	"github.com/mrbroll/foodlab/recipe"
)

const (
	schema string = `<recipe>: uid .
<ingredient>: uid @reverse .
<instruction>: uid @reverse .
<food>: uid @reverse .
<measurement>: uid @reverse .
<nutrient_measurement>: uid @reverse .
<nutrient>: uid @reverse .
<name>: string @index(fulltext, hash, term, trigram) .
<unit>: string @index(term) .
<value>: float @index(float) .
<eq_unit>: string @index(term) .
<eq_value>: float @index(float) .
<ndb_id>: string @index(hash) .
<ndb_group>: string @index(hash, term) .
<order>: int @index(int) .
<text>: string @index(fulltext, term, trigram) .`
)

const (
	ndbToken   string = "O7G3lWJlaYkMAnCeaFFhd7rM0wmdmTR2xkJmOslZ"
	dgraphHost string = "localhost:9080"
)

func main() {
	ndbClient := ndb.NewHTTPClient(new(http.Client), ndbToken)
	store := recipe.NewDgraphStore(dgraphHost)
	if err := store.AlterSchema(schema); err != nil {
		log.Fatal(err)
	}
	composer := recipe.NewCLIComposer(ndbClient, store)
	if err := composer.Compose(); err != nil {
		log.Fatal(err)
	}
}
