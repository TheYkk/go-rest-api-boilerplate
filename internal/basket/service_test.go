package basket

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestService(usecase *testing.T) {

	usecase.Run("NewService", func(t *testing.T) {
		type args struct {
			repo Repository
		}
		tests := []struct {
			name string
			args args
			want Service
		}{
			{name: "WithValidArgs_ShouldSuccess", args: args{repo: &mockRepository{}}, want: &service{repo: &mockRepository{} }},
			{name: "WithNullRepo_ShouldReturnNil", args: args{nil}, want: nil},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				if got := newService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewService() = %v, want %v", got, tt.want)
				}
			})
		}
	})
	usecase.Run("ReadMethods", func(t *testing.T) {

		givenBasket := Basket{
			Id:        "ID_1",
			BuyerId:   "Buyer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo := &mockRepository{items: []Basket{givenBasket}}
		loadData(mockRepo)
		s := newService(mockRepo)

		t.Run("Get Method Tests", func(t *testing.T) {

			tests := []struct {
				name       string
				args       string
				wantBasket *Basket
				wantErr    bool
			}{
				{name: "WithEmptyBasket_ShouldNotFoundError", args:  "INVALID_ID", wantBasket: nil, wantErr: false},
				{name: "WithEmptyBasket",args:  "ID_1", wantBasket: &givenBasket, wantErr: false},
			}
			ctx := context.Background()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					gotBasket, err := s.Get(ctx, tt.args)
					if (err != nil) != tt.wantErr {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if diff := cmp.Diff(gotBasket, tt.wantBasket); diff != "" {
						t.Errorf("Get() mismatch (-want +got):\n%s", diff)
					}
				})
			}
		})
		t.Run("Query Method Tests", func(t *testing.T) {

			type args struct {
				offset int
				limit  int
			}
			tests := []struct {
				name    string
				args    []int
				want    int
				wantErr bool
			}{
				{name: "First Page", args: []int{0, 5}, want: 5, wantErr: false},
				{name: "Second Page", args: []int{5, 5}, want: 5, wantErr: false},
				{name: "Fourt Page", args: []int{15, 5}, want: 5, wantErr: false},
				{name: "Near to end", args: []int{98, 5}, want: 4, wantErr: false},
				{name: "End Of offset", args: []int{102, 5}, want: 0, wantErr: false},
				{name: "No Rows", args: []int{150, 155}, want: 0, wantErr: false},
				{name: "Get fetch error", args: []int{-1, 1}, want: 0, wantErr: true},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					ctx := context.Background()
					got, err := s.Query(ctx, tt.args[0], tt.args[1])
					if (err != nil) != tt.wantErr {
						t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					t.Log(len(got))
					if len(got) != tt.want {
						t.Errorf("Query() got = %v, want %v", got, tt.want)
					}
				})
			}

		})
		t.Run("Count Tests", func(t *testing.T) {

			 var want int64= 102
			ctx := context.Background()
			got, err := s.Count(ctx)
			if err != nil {
				t.Errorf("Count() error = %v", err)
				return
			}
			if got != want {
				t.Errorf("Count() got = %v, want %v", got, want)
			}

		})
	})
	usecase.Run("Crud Operations", func(t *testing.T) {

		givenBasket := Basket{
			Id:        "CantDelete",
			BuyerId:   "Buyer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo := &mockRepository{items: []Basket{givenBasket}}
		loadData(mockRepo)
		s := newService(mockRepo)
		ctx := context.Background()

		t.Run("CreateBasket", func(t *testing.T) {
			t.Run("WithValidBuyer_ShouldBeSuccess", func(t *testing.T) {
				got, err := s.Create(ctx, "Buyer-X")
				if err != nil {
					t.Errorf("Count() error = %v", err)
					return
				}
				assert.NotNil(t, got)

			})
			t.Run("WithErrorBuyer_ShouldBeFailed", func(t *testing.T) {

				_, err := s.Create(ctx, "error")
				assert.Equal(t, errCRUD, errors.Cause(err))

			})
		})
		t.Run("DeleteBasket", func(t *testing.T) {
			t.Run("WithValidBasket_ShouldBeSuccsess", func(t *testing.T) {
				given, _ := s.Get(ctx, "ID_1")
				got, err := s.Delete(ctx, "ID_1")
				assert.NoError(t, err)
				if diff := cmp.Diff(got, given); diff != "" {
					t.Errorf("Get() mismatch (-want +got):\n%s", diff)
				}

			})
			t.Run("WithNotFoundBasketId_ShouldBeFailed", func(t *testing.T) {
				_, err := s.Delete(ctx, "NotFound")
				t.Log(err)
				assert.Equal(t, err.Error(),"Service: Basket not found")

			})
			t.Run("WithemptyBasketId_ShouldBeFailed", func(t *testing.T) {
				_, err := s.Delete(ctx, "")
				t.Log(err)
				assert.Equal(t, "sql: no rows in result set",errors.Cause(err).Error())

			})

			t.Run("WithExistBasketIdButCantDelete_ShouldBeFailed", func(t *testing.T) {
				_, err := s.Delete(ctx, "CantDelete")
				t.Log(err)
				assert.Equal(t, errCRUD, errors.Cause(err))

			})
		})
		t.Run("AddItem", func(t *testing.T) {

			tests := []struct {
				name    string
				args    []string
				want    string
				wantErr bool
				wantStr string
			}{
				{name: "WithValidItem_ShouldBeAdded", args: []string{"ID_5", "SKU_2", "5", "10"}, want: "", wantErr: false, wantStr: ""},
				{name: "WithExistItem_ShouldBeFailed", args: []string{"ID_5", "SKU_5", "5", "10"}, want: "", wantErr: true, wantStr: "Service: Item already added"},
				{name: "WithNonExistBasket_ShouldBeFailed", args: []string{"INVALID_BASKET_ID", "SKU_1", "5", "10"}, want: "", wantErr: true, wantStr: "Service: Basket not found"},
				{name: "WithNonExistBasket_ShouldBeFailed", args: []string{"", "SKU_1", "5", "10"}, want: "", wantErr: true, wantStr: "sql: no rows in result set"},

			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					_, err := s.AddItem(ctx, tt.args[0], tt.args[1], castStrToInt(tt.args[2]), int64(castStrToInt(tt.args[3])))
					if (err != nil) != tt.wantErr {
						t.Errorf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
					}
					if (err != nil) && tt.wantErr && errors.Cause(err).Error() != tt.wantStr {
						t.Errorf("AddItem() error = %v, wantErr %v", errors.Cause(err).Error(), tt.wantStr)
					}
				})
			}
		})
		t.Run("UpdateItem", func(t *testing.T) {

			tests := []struct {
				name    string
				args    []string
				wantErr bool
				wantStr string
			}{
				{name: "WithValidItemParameters_ShouldBeSuccess", args: []string{"ID_UPDATE", "ITEM_UPDATE", "9"}, wantErr: false, wantStr: ""},
				{name: "WithInvalidItemIdParameters_ShouldBeSuccess", args: []string{"ID_UPDATE", "INVALID_ITEM_ID", "9"}, wantErr: true, wantStr: "Item can not found. ItemId : INVALID_ITEM_ID"},
				{name: "WithInvalidBasketIdParameters_ShouldBeSuccess", args: []string{"ID_INVALID_ID", "ITEM_UPDATE", "9"}, wantErr: true, wantStr: "Service: Basket not found"},
				{name: "WithInvalidBasketIdParameters_ShouldBeSuccess", args: []string{"", "ITEM_UPDATE", "9"}, wantErr: true, wantStr: "sql: no rows in result set"},

			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					err := s.UpdateItem(ctx, tt.args[0], tt.args[1], castStrToInt(tt.args[2]))
					if (err != nil) != tt.wantErr {
						t.Errorf("UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
					}
					if (err != nil) && tt.wantErr && errors.Cause(err).Error() != tt.wantStr {
						t.Errorf("AddItem() error = %v, wantErr %v", errors.Cause(err).Error(), tt.wantStr)
					}
				})
			}
		})
		t.Run("DeleteItem", func(t *testing.T) {

			tests := []struct {
				name    string
				args    []string
				wantErr bool
			}{
				{name: "WithValidItem_ShouldBeSuccess", args: []string{"ID_DELETE", "ITEM_DELETE"}, wantErr: false},
				{name: "WithValidItem_ShouldBeSuccess", args: []string{"ID_DELETE", "INVALID_ITEM_DELETE"}, wantErr: true},
				{name: "WithValidItem_ShouldBeSuccess", args: []string{"INVALID_ID_DELETE", "ITEM_DELETE"}, wantErr: true},
				{name: "WithValidItem_ShouldBeSuccess", args: []string{"", "ITEM_DELETE"}, wantErr: true},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {



					if err := s.DeleteItem(ctx, tt.args[0], tt.args[1]); (err != nil) != tt.wantErr {
						t.Errorf("DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
					}
				})
			}

		})
	})
}

func castStrToInt(s string) int {
	if i, ok := strconv.Atoi(s); ok == nil {
		return i
	}
	return 0
}

/*
MockRepository here
*/
var errCRUD = errors.New("Mock: Error crud operation")

type mockRepository struct {
	items []Basket
}

func (m mockRepository) Get(ctx context.Context, id string) (*Basket, error) {

	if len(id) == 0{
		return nil,sql.ErrNoRows
	}

	for _, item := range m.items {
		if item.Id == id {
			return &item, nil
		}
	}
	return nil, nil
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int) ([]Basket, error) {

	if offset == -1 {
		return nil, errCRUD
	}
	if offset > len(m.items) {
		offset = len(m.items)
	}

	end := offset + limit
	if end > len(m.items) {
		end = len(m.items)
	}

	return m.items[offset:end], nil
}

func (m *mockRepository) Create(ctx context.Context, basket *Basket) error {
	if basket.BuyerId == "error" {
		return errCRUD
	}
	m.items = append(m.items, *basket)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, basket *Basket) error {
	if basket.BuyerId == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.Id == basket.Id {
			m.items[i] = *basket
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	if id == "CantDelete" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.Id == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}

func loadData(repo *mockRepository) {
	ctx := context.TODO()

	for i := 1; i < 100; i++ {

		basket := &Basket{
			Id:        fmt.Sprintf("ID_%v", i),
			BuyerId:   fmt.Sprintf("Buyer_%v", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		basket.AddItem(1, 5, fmt.Sprintf("SKU_%v", i))
		repo.Create(ctx, basket)
	}

	//for UpdateItem
	repo.Create(ctx, &Basket{
		Id:      "ID_UPDATE",
		BuyerId: "Buyer",
		Items: []Item{{
			Id:        "ITEM_UPDATE",
			Sku:       "SKU",
			UnitPrice: 5,
			Quantity:  2,
		}},
		CreatedAt: time.Now(),
	})
	//for UpdateItem
	repo.Create(ctx, &Basket{
		Id:      "ID_DELETE",
		BuyerId: "Buyer",
		Items: []Item{{
			Id:        "ITEM_DELETE",
			Sku:       "SKU",
			UnitPrice: 5,
			Quantity:  2,
		}},
		CreatedAt: time.Now(),
	})

}
