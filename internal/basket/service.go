package basket

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

// Service encapsulates use case logic for basket.
type Service interface {
	Get(ctx context.Context, id string) (*Basket, error)
	Query(ctx context.Context, offset, limit int) ([]Basket, error)
	Count(ctx context.Context) (int64, error)

	Create(ctx context.Context, buyer string) (*Basket, error)
	Delete(ctx context.Context, id string) (*Basket, error)

	UpdateItem(ctx context.Context, basketId, itemId string, quantity int) error
	AddItem(ctx context.Context, basketId, sku string, quantity int, price int64) (string, error)
	DeleteItem(ctx context.Context, basketId, itemId string) error
}

type service struct {
	repo Repository
}

func (s *service) Get(ctx context.Context, id string) (basket *Basket, err error) {

	basket, err = s.repo.Get(ctx, id)
	if err != nil {
		err = errors.Wrapf(err, "get basket error. Basket Id:%s", id)
	}

	return
}

func (s *service) Query(ctx context.Context, offset, limit int) ([]Basket, error) {

	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, "Service:Failed to querying the basket")
	}

	return items, nil
}

func (s *service) Count(ctx context.Context) (int64, error) {

	return s.repo.Count(ctx)
}

// Create creates a new basket
func (s *service) Create(ctx context.Context, buyer string) (*Basket, error) {

	//TODO validation?
	basket := &Basket{
		Id:        GenerateId(),
		BuyerId:   buyer,
		Items:     nil,
		CreatedAt: time.Now(),
	}
	err := s.repo.Create(ctx, basket)

	if err != nil {
		return nil, errors.Wrap(err, "Service:Failed to create basket")
	}
	return basket, nil

}

func (s *service) AddItem(ctx context.Context, basketId, sku string, quantity int, price int64) (string, error) {

	basket, err := s.repo.Get(ctx, basketId)
	if err != nil {
		return "", errors.Wrap(err, "Service: Get basket error. Basket Id.")
	}
	if basket ==nil{
		return "",errors.New("Service: Basket not found")
	}
	item, err := basket.AddItem(quantity, price, sku)

	if err != nil {
		return "", errors.Wrap(err, "Service: Failed to item added to basket.")
	}
	if err:=s.repo.Update(ctx,basket);err!=nil{
		return "",errors.Wrap(err,"Service: Failed to update basket in data storage.")
	}

	return item.Id, nil
}
func (s *service) UpdateItem(ctx context.Context, basketId, itemId string, quantity int) error {

	basket, err := s.repo.Get(ctx, basketId)
	if err != nil {
		return errors.Wrapf(err, "Service: Get basket error. Basket Id:%s", basketId)
	}
	if basket ==nil{
		return errors.New("Service: Basket not found")
	}
	err = basket.UpdateItem(itemId, quantity)

	if err != nil {
		return errors.Wrapf(err, "Service: Failed to update item")
	}
	if err:=s.repo.Update(ctx,basket);err!=nil{
		return errors.Wrap(err,"Service: Failed to update basket in data storage.")
	}
	return nil
}

func (s *service) DeleteItem(ctx context.Context, basketId, itemId string) error {

	basket, err := s.repo.Get(ctx, basketId)
	if err != nil {
		return errors.Wrapf(err, "Service: Get basket error. Basket Id:%s", basketId)
	}
	if basket ==nil{
		return errors.New("Service: Basket not found")
	}
	err = basket.RemoveItem(itemId)
	if err != nil {
		return errors.Wrap(err, "Service: Basket Item not found.")
	}
	if err:=s.repo.Update(ctx,basket);err!=nil{
		return errors.Wrap(err,"Service: Failed to update basket in data storage.")
	}
	return nil
}

//Delete deletes the basket with the spesified Id
func (s *service) Delete(ctx context.Context, id string) (*Basket, error) {
	basket, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if basket ==nil{
		return nil,errors.New("Service: Basket not found")
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return nil, errors.Wrap(err, "Service:Failed to delete basket")
	}
	return basket, nil
}

// NewService creates a new basket service.
func newService(repo Repository) Service {

	if repo == nil {
		return nil
	}
	return &service{repo}
}
