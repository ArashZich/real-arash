package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/ARmo-BigBang/kit/restypes"
	"github.com/labstack/echo/v4"
)

func (r resource) delete(ctx echo.Context) error {
	var ids []int

	e := ctx.Bind(&ids)
	if e != nil {
		r.logger.With(ctx.Request().Context()).Error(e)
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil))
	}

	deletedUsers, err := r.service.Delete(ctx.Request().Context(), ids)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	result := restypes.DeleteResponse{
		IDs: deletedUsers,
	}
	return ctx.JSON(http.StatusOK, response.Success(result))
}
