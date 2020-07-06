package basket

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApi(usecase *testing.T) {

	e := echo.New()
	mockRepo := &mockRepository{}
	RegisterHandlers(e, mockRepo)
	loadData(mockRepo)

	resource := &resource{
		service: newService(mockRepo),
	}

	usecase.Run("Get Operations", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/basket", nil)

		tests := []struct {
			name       string
			args       string
			wantStatus int
		}{
			{name: "Get Basket", args: "ID_1", wantStatus: http.StatusOK},
			{name: "Get Basket", args: "", wantStatus: http.StatusInternalServerError},
			//TODO :  the resCode should be 404 instead of 500
			{name: "Get Basket", args: "INVALID_BASKET_ID", wantStatus: http.StatusNotFound},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				rec := httptest.NewRecorder()
				ctx := e.NewContext(req, rec)
				ctx.SetPath("/:id")
				ctx.SetParamNames("id")
				ctx.SetParamValues(tt.args)
				res := rec.Result()
				defer res.Body.Close()
				if err := resource.getBasket(ctx); err == nil {
					assert.Equal(t, tt.wantStatus, rec.Code)
					assert.NotNil(t, rec.Body.String())
				} else {
					t.Errorf("getBasket() error = %v", err)
				}
				t.Logf("Response:%v", rec)

			})
		}

	})

	usecase.Run("Delete Operations", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/basket", nil)

		tests := []struct {
			name    string
			args    string
			wantStatus int
		}{
			{name: "Delete Basket", args: "ID_1", wantStatus: http.StatusAccepted},
			{name: "WithEmptyBasketId_ShouldBeFailed", args: "", wantStatus: http.StatusNotFound},
			{name: "WithInvalidBasketId_ShouldBeFailed", args: "INVALID_BASKET_ID", wantStatus: http.StatusInternalServerError},

		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				rec := httptest.NewRecorder()
				ctx := e.NewContext(req, rec)
				ctx.SetPath("/:id")
				ctx.SetParamNames("id")
				ctx.SetParamValues(tt.args)
				res := rec.Result()
				defer res.Body.Close()

				err := resource.deleteBasket(ctx)
				t.Logf("ERR:%v",err)
				if err == nil {
					assert.Equal(t, tt.wantStatus, rec.Code)
					assert.NotNil(t, rec.Body.String())
				}

			})
		}

	})

}
