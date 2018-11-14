package main

import (
	"log"
	"net/http"

	"github.com/mrbroll/foodlab/ndb"
	"github.com/mrbroll/foodlab/recipe"
)

const (
	schema string = `<ingredient>: uid @reverse .
<instruction>: uid @reverse .
<name>: string @index(fulltext, hash, term, trigram) @lang .
<order>: int @index(int) .
<preparation>: string @index(fulltext, term, trigram) @lang .
<measurement.unit>: string @index(term) .
<measurement.value>: float @index(float) .
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
