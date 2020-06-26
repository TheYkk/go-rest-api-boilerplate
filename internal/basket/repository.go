package basket

import "context"

// Repository encapsulates the logic to access basket from the data source.
type Repository interface {
	// Get returns the album with the specified basket Id.
	Get(ctx context.Context, id string) (*Basket, error)
	// Count returns the number of basket.
	Count(ctx context.Context) (int, error)
	// Query returns the list of baskets with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]Basket, error)
	// Create saves a new basket in the storage.
	Create(ctx context.Context, basket *Basket) error
	// Update updates the basket with given Is in the storage.
	Update(ctx context.Context, basket *Basket) error
	// Delete removes the basket with given Is from the storage.
	Delete(ctx context.Context, id string) error
}

func (r repository) Get(ctx context.Context, id string) (*Basket, error) {
	return nil, nil
}

func (r repository) Count(ctx context.Context) (int, error) {
	return 0, nil
}

func (r repository) Query(ctx context.Context, offset, limit int) ([]Basket, error) {
	return nil, nil
}

func (r repository) Create(ctx context.Context, basket *Basket) error {
	return nil
}

func (r repository) Update(ctx context.Context, basket *Basket) error {
	return nil
}

func (r repository) Delete(ctx context.Context, id string) error {
	return nil
}

type repository struct {
}

func NewRepository() Repository {
	return repository{}
}
