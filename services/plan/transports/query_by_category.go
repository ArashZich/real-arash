package transports

import (
	"net/http"
	"strconv"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) queryByCategoryID(ctx echo.Context) error {
	categoryIDStr := ctx.QueryParam("category_id")
	if categoryIDStr == "" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil, "category_id لازم است"))
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil, "category_id معتبر نیست"))
	}

	plans, errResp := r.service.QueryByCategoryID(ctx.Request().Context(), categoryID)
	if errResp.StatusCode != 0 {
		return ctx.JSON(errResp.StatusCode, errResp)
	}

	return ctx.JSON(http.StatusOK, response.Success(plans))
}
