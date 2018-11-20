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
<food.measurement>: uid @reverse .
<nutrient.measurement>: uid @reverse .
<nutrient>: uid @reverse .
<name>: string @index(fulltext, hash, term, trigram) @lang .
<unit>: string @index(term) .
<value>: float @index(float) .
<eq.unit>: string @index(term) .
<eq.value>: float @index(float) .
<ndb.id>: string @index(hash) .
<ndb.group>: string @index(hash, term) .
<order>: int @index(int) .
<text>: string @index(fulltext, term, trigram) @lang .`
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
