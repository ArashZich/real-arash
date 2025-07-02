package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/pagination"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/ARmo-BigBang/kit/restypes"
	"github.com/labstack/echo/v4"
)

func (r resource) query(ctx echo.Context) error {
	query := ctx.QueryParam("query")
	order := ctx.QueryParam("order")
	orderBy := ctx.QueryParam("order_by")
	if orderBy == "" {
		orderBy = "id"
	}
	if order == "" {
		order = "desc"
	}

	c := ctx.Request().Context()

	count, err := r.service.Count(c, query)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	pages := pagination.NewFromRequest(ctx.Request(), int(count))
	invites, er := r.service.Query(c, pages.Offset(), pages.Limit(), orderBy, order, query)
	if er.StatusCode != 0 {
		return ctx.JSON(er.StatusCode, er)
	}
	result := restypes.QueryResponse{
		Limit:      pages.Limit(),
		Offset:     pages.Offset(),
		Page:       pages.Page,
		TotalRows:  int64(pages.TotalCount),
		TotalPages: pages.PageCount,
		Items:      invites,
	}
	return ctx.JSON(http.StatusOK, response.Success(result))
}
