package basket

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

var (
	maxQuantityPerProduct          = 10
	minCartAmountForCheckout int64 = 50

	ErrNoBuyer  = errors.New("Buyer field can not null or empty")
	ErrNotFound = errors.New("Item not found")
)

type (
	Basket struct {
		Id        string    `json:"id"`
		BuyerId   string    `json:"buyerId"`
		Items     []Item    `json:"items"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	Item struct {
		Id        string `json:"id"`
		Sku       string `json:"sku"`
		UnitPrice int64  `json:"unitPrice"`
		Quantity  int    `json:"quantity"`
	}
)

func Create(buyer string) (*Basket, error) {

	if len(buyer) == 0 {
		return nil, ErrNoBuyer
	}

	return &Basket{
		Id:        GenerateId(),
		BuyerId:   buyer,
		Items:     nil,
		CreatedAt: time.Now(),
	}, nil
}

func (b *Basket) AddItem(quantity int, price int64, sku string) (*Item, error) {

	if quantity >= maxQuantityPerProduct {
		return nil, errors.Errorf("You can't add more item. Item count can be less then %d", maxQuantityPerProduct)
	}
	_, item := b.SearchItemBySku(sku)
	if item != nil {
		return item, errors.New("Service: Item already added")
	}

	item = &Item{
		Id:        GenerateId(),
		Sku:       sku,
		UnitPrice: price,
		Quantity:  quantity,
	}
	b.Items = append(b.Items, *item)
	b.UpdatedAt = time.Now()
	return item, nil
}
func (b *Basket) UpdateItem(itemId string, quantity int) (err error) {

	if index, item := b.SearchItem(itemId); index != -1 {

		if quantity >= maxQuantityPerProduct {
			return errors.Errorf("You can't add more item. Item count can be less then %d", maxQuantityPerProduct)
		}

		item.Quantity = quantity
	} else {
		return errors.Errorf("Item can not found. ItemId : %s", itemId)
	}
	b.UpdatedAt = time.Now()
	return
}

func (b *Basket) RemoveItem(itemId string) (err error) {

	if index, _ := b.SearchItem(itemId); index != -1 {
		b.Items = append(b.Items[:index], b.Items[index+1:]...)
	} else {
		return ErrNotFound
	}
	b.UpdatedAt = time.Now()
	return
}

func (b *Basket) SearchItem(itemId string) (int, *Item) {

	for i, n := range b.Items {
		if n.Id == itemId {
			return i, &n
		}
	}
	return -1, nil
}
func (b *Basket) SearchItemBySku(sku string) (int, *Item) {

	for i, n := range b.Items {
		if n.Sku == sku {
			return i, &n
		}
	}
	return -1, nil
}

func (b *Basket) ValidateBasket() error {

	totalPrice := calculateBasketAmount(b)

	if totalPrice <= minCartAmountForCheckout {
		return errors.Errorf("Total basket amount must be greater then %d", minCartAmountForCheckout)
	}
	return nil
}

func calculateBasketAmount(b *Basket) (totalPrice int64) {

	for _, item := range b.Items {
		totalPrice += int64(item.Quantity) * item.UnitPrice
	}
	return
}

// GenerateId generates a unique ID that can be used as an identifier for an entity.
func GenerateId() string {
	return uuid.New().String()
}
