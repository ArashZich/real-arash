// package transports

// import (
// 	"net/http"

// 	"gitag.ir/armogroup/armo/services/reality/services/view/endpoints"
// 	"github.com/ARmo-BigBang/kit/pagination"
// 	"github.com/ARmo-BigBang/kit/response"
// 	"github.com/ARmo-BigBang/kit/restypes"
// 	"github.com/labstack/echo/v4"
// )

// func (r *resource) query(ctx echo.Context) error {
// 	var params = &endpoints.ViewQueryRequestParams{}

// 	errs := params.ValidateQueries(ctx.Request())
// 	if len(errs) > 0 {
// 		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errs, "خطا در داده ورودی"))
// 	}
// 	if params.OrderBy == "" {
// 		params.OrderBy = "id"
// 	}

// 	if params.Order == "" {
// 		params.Order = "desc"
// 	}
// 	count, err := r.service.Count(ctx.Request().Context(), *params)
// 	if err.StatusCode != 0 {
// 		return ctx.JSON(err.StatusCode, err)
// 	}
// 	pages := pagination.NewFromRequest(ctx.Request(), int(count))

// 	views, er := r.service.Query(
// 		ctx.Request().Context(),
// 		pages.Offset(), pages.Limit(), *params,
// 	)
// 	if er.StatusCode != 0 {
// 		return ctx.JSON(err.StatusCode, err)
// 	}
// 	pages.Items = views

// 	result := restypes.QueryResponse{
// 		Limit:      pages.Limit(),
// 		Offset:     pages.Offset(),
// 		Page:       pages.Page,
// 		TotalRows:  int64(pages.TotalCount),
// 		TotalPages: pages.PageCount,
// 		Items:      views,
// 	}

// 	return ctx.JSON(http.StatusOK, response.Success(result))
// }

package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/view/endpoints"
	"github.com/ARmo-BigBang/kit/pagination"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/ARmo-BigBang/kit/restypes"
	"github.com/labstack/echo/v4"
)

func (r *resource) query(ctx echo.Context) error {
	var params = &endpoints.ViewQueryRequestParams{}

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

	// Determine if we need to fetch all records
	page := ctx.QueryParam("page")
	perPage := ctx.QueryParam("per_page")

	// Fetch all records if page and per_page are "0"
	if page == "0" && perPage == "0" {
		allViews, er := r.service.Query(ctx.Request().Context(), 0, int(count), *params)
		if er.StatusCode != 0 {
			return ctx.JSON(er.StatusCode, er)
		}

		result := restypes.QueryResponse{
			Limit:      int(count),
			Offset:     0,
			Page:       1,
			TotalRows:  int64(count),
			TotalPages: 1,
			Items:      allViews,
		}
		return ctx.JSON(http.StatusOK, response.Success(result))
	}

	// Use pagination for other cases
	pages := pagination.NewFromRequest(ctx.Request(), int(count))
	paginatedViews, er := r.service.Query(ctx.Request().Context(), pages.Offset(), pages.Limit(), *params)
	if er.StatusCode != 0 {
		return ctx.JSON(er.StatusCode, er)
	}
	pages.Items = paginatedViews

	result := restypes.QueryResponse{
		Limit:      pages.Limit(),
		Offset:     pages.Offset(),
		Page:       pages.Page,
		TotalRows:  int64(pages.TotalCount),
		TotalPages: pages.PageCount,
		Items:      paginatedViews,
	}

	return ctx.JSON(http.StatusOK, response.Success(result))
}
