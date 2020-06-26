package basket

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	baseUrl = "basket"
)

func RegisterHandlers(instance *echo.Echo, repo Repository) {

	res := &resource{
		service: newService(repo),
	}

	instance.GET(fmt.Sprintf("%s/:id", baseUrl), res.getBasket)
	instance.PUT(fmt.Sprintf("%s/:id", baseUrl), res.deleteBasket)
	instance.DELETE(fmt.Sprintf("%s", baseUrl), res.deleteBasket)
	instance.DELETE(fmt.Sprintf("%s/:id/:itemId", baseUrl), res.deleteItem)
	instance.POST(fmt.Sprintf("%s", baseUrl), res.createBasket)
	instance.POST(fmt.Sprintf("%s", baseUrl), res.createBasket)

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
		return ctx.JSON(http.StatusNotFound, err.Error())
	}
	return ctx.JSON(http.StatusOK, result)

}
func (r *resource) createBasket(ctx echo.Context) error {
	return nil
}
func (r *resource) deleteBasket(ctx echo.Context) error {
	return nil
}
func (r *resource) deleteItem(ctx echo.Context) error {
	return nil
}
