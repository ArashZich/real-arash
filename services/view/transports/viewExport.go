package transports

import (
	"fmt"
	"net/http"
	"time"

	"gitag.ir/armogroup/armo/services/reality/services/view/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) exportViews(ctx echo.Context) error {
	var request endpoints.ExportViewsRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil, "Invalid request body"))
	}

	if err := request.Validate(); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil, err.Error()))
	}

	params := endpoints.ViewQueryRequestParams{
		Duration: request.Duration,
		Filters:  request.Filters,
		Order:    "desc", // Set default order
		OrderBy:  "id",   // Set default orderBy
	}

	data, contentType, err := r.service.ExportViews(ctx.Request().Context(), request.Format, params)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("views-export-%s.%s", timestamp, request.Format)

	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Response().Header().Set("Content-Type", contentType)

	return ctx.Blob(http.StatusOK, contentType, data)
}
