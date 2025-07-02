package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/view/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) getProductView(ctx echo.Context) error {

	var params = &endpoints.ViewGetProductRequestParams{}

	errs := params.ValidateQueries(ctx.Request())
	if len(errs) > 0 {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errs, "خطا در داده ورودی"))
	}
	if params.OrderBy == "" {
		params.OrderBy = "id"
	}

	if params.Order == "" {
		params.Order = "desc"
	}

	product, err := r.service.GetProductView(ctx.Request().Context(), *params)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusOK, response.Success(product))
}
