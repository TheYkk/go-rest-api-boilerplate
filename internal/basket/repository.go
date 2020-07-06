package basket

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository encapsulates the logic to access basket from the data source.
type Repository interface {
	// Get returns the album with the specified basket Id.
	Get(ctx context.Context, id string) (*Basket, error)
	// Count returns the number of basket.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of baskets with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]Basket, error)
	// Create saves a new basket in the storage.
	Create(ctx context.Context, basket *Basket) error
	// Update updates the basket with given Is in the storage.
	Update(ctx context.Context, basket *Basket) error
	// Delete removes the basket with given Is from the storage.
	Delete(ctx context.Context, id string) error
}

func (r repository) Get(ctx context.Context, id string) (basket *Basket,err error) {

	err = r.collection.FindOne(ctx, getId(id)).Decode(&basket)

	return basket,err
}

func (r repository) Count(ctx context.Context) (int64, error) {

	return r.collection.CountDocuments(ctx,bson.M{})
}

func (r repository) Query(ctx context.Context, offset, limit int) ([]Basket, error) {
	return nil, nil
}

func (r repository) Create(ctx context.Context, basket *Basket) error {

	_,err :=r.collection.InsertOne(ctx,basket)
	return err
}

func (r repository) Update(ctx context.Context, basket *Basket) error {

	replaceOptions := options.Replace().SetUpsert(true)
	objectId := getId(basket.Id)

	_,err:=r.collection.ReplaceOne(ctx,objectId,basket,replaceOptions)

	return err
}

func (r repository) Delete(ctx context.Context, id string) error {

	_,err := r.collection.DeleteOne(ctx, getId(id) ,options.Delete())
	return err
}

func getId(id string) bson.M{

	return bson.M{ "id" : id }
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {

	col := db.Collection("baskets")

	return &repository{
		collection:col,
	}
}

