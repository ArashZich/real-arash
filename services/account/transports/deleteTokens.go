package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) deleteTokens(ctx echo.Context) error {
	var accessTokens []string
	e := ctx.Bind(&accessTokens)
	if e != nil {
		r.logger.With(ctx.Request().Context()).Error(e)
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil))
	}
	deleteTokens, err := r.service.DeleteTokens(ctx.Request().Context(), accessTokens)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil, "خطایی در حذف توکن ها رخ داده است"))
	}
	return ctx.JSON(http.StatusOK, response.Success(deleteTokens))
}
