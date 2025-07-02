package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) getDocumentsByProductUID(ctx echo.Context) error {
	productUID := ctx.Param("product_uid")
	documents, err := r.service.GetDocumentsByProductUID(ctx.Request().Context(), productUID)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusOK, response.Success(documents))
}
