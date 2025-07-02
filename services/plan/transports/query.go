package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/plan/endpoints"
	"github.com/ARmo-BigBang/kit/pagination"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/ARmo-BigBang/kit/restypes"
	"github.com/labstack/echo/v4"
)

func (r *resource) query(ctx echo.Context) error {
	var params = &endpoints.PlanQueryRequestParams{}

	errs := params.ValidateQueries(ctx.Request())
	if len(errs) > 0 {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errs, "خطا در داده ورودی"))
	}
	if params.OrderBy == "" {
		params.OrderBy = "id"
	}

	if params.Order == "" {
		params.Order = "desc"
	}
	count, err := r.service.Count(ctx.Request().Context(), *params)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	pages := pagination.NewFromRequest(ctx.Request(), int(count))

	plans, er := r.service.Query(
		ctx.Request().Context(),
		pages.Offset(), pages.Limit(), *params,
	)
	if er.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	pages.Items = plans

	result := restypes.QueryResponse{
		Limit:      pages.Limit(),
		Offset:     pages.Offset(),
		Page:       pages.Page,
		TotalRows:  int64(pages.TotalCount),
		TotalPages: pages.PageCount,
		Items:      plans,
	}

	return ctx.JSON(http.StatusOK, response.Success(result))
}
