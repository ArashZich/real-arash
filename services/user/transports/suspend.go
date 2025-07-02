package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) suspend(ctx echo.Context) error {
	id := ctx.Param("id")
	user, err := r.service.Suspend(ctx.Request().Context(), id)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusOK, response.Success(user))
}
