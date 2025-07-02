package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/invoice/service"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) pay(ctx echo.Context) error {
	var input = &service.PayInvoiceRequest{}

	errors := input.Validate(ctx.Request())
	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	invoice, err := r.Invoice.Pay(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusCreated, response.Created(invoice))
}
