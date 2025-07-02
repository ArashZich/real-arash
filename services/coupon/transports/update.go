package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/coupon/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) update(ctx echo.Context) error {
	id := ctx.Param("id")
	var input = &endpoints.UpdateCouponRequest{}
	errors := input.Validate(ctx.Request())

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	coupon, err := r.service.Update(ctx.Request().Context(), id, *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusCreated, response.Success(coupon))
}
