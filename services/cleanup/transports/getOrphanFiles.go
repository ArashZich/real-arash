package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) getOrphanFiles(ctx echo.Context) error {
	orphanFiles, err := r.service.GetOrphanFiles(ctx.Request().Context())
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(orphanFiles))
}
