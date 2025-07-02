package transports

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r resource) serveDocument(ctx echo.Context) error {
	id := ctx.Param("id")

	documentHtml, err := r.service.ServeDocument(ctx.Request().Context(), id)
	if err.StatusCode != 0 {
		return ctx.HTML(err.StatusCode, err.Message)
	}
	return ctx.HTML(http.StatusOK, documentHtml)
}
