package recipe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// RecipeAdder is a type for adding recipes to an underlying storage medium.
type RecipeAdder interface {
	// AddRecipe adds the given recipe to the store.
	// It returns an error if the operation was unsuccessful.
}

// DgraphStore is a recipe store backed by dgraph.
type DgraphStore struct {
	host string
}

// NewDgraphStore returns a new recipe store for the given Dgraph host.
func NewDgraphStore(host string) *DgraphStore {
	return &DgraphStore{host: host}
}

// AddRecipe adds the given recipe to the store.
// It returns an error if the operation was unsuccessful.
func (s *DgraphStore) AddRecipe(r *Recipe) error {
	conn, err := grpc.Dial(s.host, grpc.WithInsecure())
	if err != nil {
		return errors.Wrapf(err, "Dialing grpc host %s.", s.host)
	}
	defer conn.Close()
	client := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	rb, err := json.Marshal(r)
	if err != nil {
		return errors.Wrap(err, "Marshaling recipe to JSON.")
	}
	mu := &api.Mutation{
		SetJson:   rb,
		CommitNow: true,
	}
	if _, err := client.NewTxn().Mutate(context.Background(), mu); err != nil {
		return errors.Wrap(err, "Committing mutation txn.")
	}
	return nil
}

// GetOrCreateFood gets a food node matching the ndb food,
// or creates one if it does not exist in the store.
// It returns an error if the store could not be reached.
func (s *DgraphStore) GetOrCreateFood(f *Food) (*Food, error) {
	conn, err := grpc.Dial(s.host, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "Dialing grpc host %s.", s.host)
	}
	defer conn.Close()
	client := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	vars := map[string]string{"$name": f.Name}
	q := `
	query Food($name: string) {
		food(func: eq(name@., $name)) {
			uid
			ndb.id
			name
			measurement {
				uid
				unit
				value
				eq.unit
				eq.value
				nutrient.measurement {
					uid
					unit
					value
					nutrient {
						uid
						ndb.id
						ndb.group
						name
					}
				}
			}
		}
	}`
	qResp, err := client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return nil, errors.Wrap(err, "Querying for food node.")
	}

	type Root struct {
		Foods []*Food `json:"food"`
	}
	r := new(Root)
	if err := json.Unmarshal(qResp.Json, r); err != nil {
		return nil, errors.Wrap(err, "Unmarshaling respons JSON.")
	}
	if r != nil && len(r.Foods) > 0 {
		return r.Foods[0], nil
	}

	// deep-copy f and set uids so we can reference them in the repsonse.
	food := &Food{
		UID:          fmt.Sprintf("_:%s", f.Hash()),
		NDBID:        f.NDBID,
		Name:         f.Name,
		Measurements: make([]*FoodMeasurement, len(f.Measurements)),
	}

	for iFM, fMeas := range f.Measurements {
		food.Measurements[iFM] = &FoodMeasurement{
			UID:     fmt.Sprintf("_:%s", fMeas.Hash()),
			Unit:    fMeas.Unit,
			Value:   fMeas.Value,
			EqUnit:  fMeas.EqUnit,
			EqValue: fMeas.EqValue,
			NutrientMeasurements: make(
				[]*NutrientMeasurement, len(fMeas.NutrientMeasurements),
			),
		}
		for iNM, nMeas := range fMeas.NutrientMeasurements {
			nutrient, err := s.GetOrCreateNutrient(nMeas.Nutrient)
			if err != nil {
				return nil, errors.Wrap(err, "Getting or creating nutrient.")
			}
			food.Measurements[iFM].NutrientMeasurements[iNM] = &NutrientMeasurement{
				UID:      fmt.Sprintf("_:%s", nMeas.Hash()),
				Nutrient: nutrient,
				Unit:     nMeas.Unit,
				Value:    nMeas.Value,
			}
		}
	}

	fb, err := json.Marshal(food)
	if err != nil {
		return nil, errors.Wrap(err, "Marshaling food JSON.")
	}
	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   fb,
	}
	mResp, err := client.NewTxn().Mutate(context.Background(), mu)
	if err != nil {
		return nil, errors.Wrap(err, "Adding food to store.")
	}

	// replace uids with those from the mutation.
	hash := food.Hash()
	food.UID = fmt.Sprintf("%s:%s", mResp.Uids[hash], hash)
	for _, fMeas := range food.Measurements {
		hash = fMeas.Hash()
		fMeas.UID = fmt.Sprintf("%s:%s", mResp.Uids[hash], hash)
		for _, nMeas := range fMeas.NutrientMeasurements {
			hash = nMeas.Hash()
			nMeas.UID = fmt.Sprintf("%s:%s", mResp.Uids[hash], hash)
		}
	}
	return food, nil
}

// GetOrCreateNutrient idempotently adds the given nutrient to the store,
// and returns it with its UID field populated.
// It returns an error if the operation was unsuccessful.
func (s *DgraphStore) GetOrCreateNutrient(n *Nutrient) (*Nutrient, error) {
	conn, err := grpc.Dial(s.host, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "Dialing grpc host %s.", s.host)
	}
	defer conn.Close()
	client := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	vars := map[string]string{"$name": n.Name}
	q := `
	query Nutrient($name: string) {
		nutrient(func: (name@., $name)) {
			uid
			ndb.id
			ndb.group
			name
		}
	}`
	qResp, err := client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return nil, errors.Wrap(err, "Querying for food node.")
	}

	type Root struct {
		Nutrients []*Nutrient `json:"nutrient"`
	}
	r := new(Root)
	if err := json.Unmarshal(qResp.Json, r); err != nil {
		return nil, errors.Wrap(err, "Unmarshaling nutrient JSON.")
	}
	if r != nil && len(r.Nutrients) > 0 {
		return r.Nutrients[0], nil
	}

	// copy n and set uid so we can reference it in mutation response
	hash := n.Hash()
	nutrient := &Nutrient{
		UID:      fmt.Sprintf("_:%s", hash),
		NDBID:    n.NDBID,
		NDBGroup: n.NDBGroup,
		Name:     n.Name,
	}
	nb, err := json.Marshal(nutrient)
	if err != nil {
		return nil, errors.Wrap(err, "Marshaling nutrient JSON.")
	}
	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   nb,
	}
	mResp, err := client.NewTxn().Mutate(context.Background(), mu)
	if err != nil {
		return nil, errors.Wrap(err, "Adding nutrient to store.")
	}
	nutrient.UID = fmt.Sprintf("%s:%s", mResp.Uids[hash], hash)
	return nutrient, nil
}

// AlterSchema alters the schema of a Dgraph store.
// It returns an error if the operator was unsuccessful.
func (s *DgraphStore) AlterSchema(schema string) error {
	conn, err := grpc.Dial(s.host, grpc.WithInsecure())
	if err != nil {
		return errors.Wrapf(err, "Dialing grpc host %s.", s.host)
	}
	defer conn.Close()
	client := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	op := &api.Operation{Schema: schema}
	if err := client.Alter(context.Background(), op); err != nil {
		return errors.Wrap(err, "Altering schema.")
	}
	return nil
}
