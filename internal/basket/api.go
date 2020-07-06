package basket

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

var (
	baseUrl = "basket"
)

func RegisterHandlers(instance *echo.Echo, repo Repository) {

	res := &resource{
		service: newService(repo),
	}

	instance.GET(fmt.Sprintf("%s/:id", baseUrl), res.getBasket)
	instance.POST(fmt.Sprintf("%s", baseUrl), res.createBasket)
	instance.DELETE(fmt.Sprintf("%s/:id", baseUrl), res.deleteBasket)

	instance.POST(fmt.Sprintf("%s/item", baseUrl), res.addItem)
	instance.DELETE(fmt.Sprintf("%s/:id/item/:itemId", baseUrl), res.deleteItem)
	instance.PUT(fmt.Sprintf("%s/:id/item/:item/quantity/:quantity", baseUrl), res.updateItem)

	instance.Validator = &CustomValidator{validator: validator.New()}

}

type resource struct {
	service Service
}

func (r *resource) getBasket(ctx echo.Context) error {

	id := ctx.Param("id")
	result, err := r.service.Get(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if result == nil {
		return ctx.JSON(http.StatusNotFound, "")
	}
	return ctx.JSON(http.StatusOK, result)

}
func (r *resource) createBasket(ctx echo.Context) error {

	entity := new(CreateBasketRequest)

	if err := ctx.Bind(entity); err != nil {
		return ctx.JSON(http.StatusBadRequest,err.Error())
	}
	if err := ctx.Validate(entity); err != nil {
		return ctx.JSON(http.StatusBadRequest,err.Error())
	}
	if b,err :=r.service.Create(ctx.Request().Context(),entity.Buyer); err!=nil{
		return ctx.JSON(http.StatusBadRequest,err.Error())
	}else {

		return  ctx.JSON(http.StatusCreated,map[string]string{"id":b.Id})
	}
}

func (r *resource) deleteBasket(ctx echo.Context) error {

	id := ctx.Param("id")
	_, err := r.service.Delete(ctx.Request().Context(), id)

	if errors.Cause(err) == sql.ErrNoRows{
		return ctx.JSON(http.StatusNotFound, err.Error())

	}
	if err != nil  {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusAccepted,"")

}
func (r *resource) addItem(ctx echo.Context) error {
	req := new(AddItemRequest)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest,err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest,err.Error())
	}
	if itemId,err := r.service.AddItem(ctx.Request().Context(),req.BasketId,req.Sku,req.Quantity,req.Price);err!=nil{
		return ctx.JSON(http.StatusBadRequest,err.Error())
	}else {
		return  ctx.JSON(http.StatusCreated,map[string]string{"id":itemId})
	}

}
func (r *resource) updateItem(ctx echo.Context) error{

	id := ctx.Param("id")
	itemId := ctx.Param("itemId")
	quantity,err := strconv.Atoi(ctx.Param("quantity"))

	if len(id)==0 || len(itemId)==0 || err!=nil || quantity<=0{
		return ctx.JSON(http.StatusBadRequest,"Failed to update item. BasketId or BasketItem Id is null or empty.")
	}
	if err:=r.service.UpdateItem(ctx.Request().Context(),id,itemId,quantity); err!=nil{
		return ctx.JSON(http.StatusBadRequest,err.Error())
	}
	return ctx.JSON(http.StatusAccepted,"")

}

func (r *resource) deleteItem(ctx echo.Context) error {

	id := ctx.Param("id")
	itemId := ctx.Param("itemId")

	if len(id)==0 || len(itemId)==0{
		return ctx.JSON(http.StatusBadRequest,"Failed to delete item. BasketId or BasketItem Id is null or empty.")
	}
	if err :=r.service.DeleteItem(ctx.Request().Context(),id,itemId);err!=nil{
		return ctx.JSON(http.StatusBadRequest,err.Error())
	}
	return ctx.JSON(http.StatusAccepted,"")
}



type (

	CreateBasketRequest struct {
		Buyer string `json:"buyer" validate:"required"`
	}

	AddItemRequest struct {
		BasketId	string 		`json:"basketId"  validate:"required"`
		Sku 		string		`json:"sku"  validate:"required"`
		Quantity 	int			`json:"quantity" validate:"required,gte=0,lte=20"`
		Price 		int64		`json:"price" validate:"required,gte=0"`
	}

)



//https://github.com/go-playground/validator framework for validation:
type CustomValidator struct {
validator *validator.Validate
}
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}