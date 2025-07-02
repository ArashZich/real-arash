// File: transports/getDocumentsSummaryByProductUID.go

package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) getDocumentSummaryByProductUID(ctx echo.Context) error {
	productUID := ctx.Param("product_uid")
	documentSummary, err := r.service.GetDocumentSummaryByProductUID(ctx.Request().Context(), productUID)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusOK, response.Success(documentSummary))
}
