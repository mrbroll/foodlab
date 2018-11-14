package recipe

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// RecipeAdder is a type for adding recipes to an underlying storage medium.
type RecipeAdder interface {
	// AddRecipe adds the given recipe to the store.
	// It returns an error if the operation was unsuccessful.
	AddRecipe(r *Recipe) error
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
