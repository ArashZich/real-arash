package transports

import (
	"encoding/json"
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
	suspendedAt := ctx.QueryParam("suspended_at")
	isOfficial := ctx.QueryParam("is_official")
	isProfileCompleted := ctx.QueryParam("is_profile_completed")
	hasOrganization := ctx.QueryParam("has_organization") // اضافه شده
	hasPackages := ctx.QueryParam("has_packages")         // اضافه شده
	affiliateCodesJSON := ctx.QueryParam("affiliate_codes")
	var affiliateCodes []string
	if affiliateCodesJSON != "" {
		err := json.Unmarshal([]byte(affiliateCodesJSON), &affiliateCodes)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid affiliate_codes format"})
		}
	}
	if orderBy == "" {
		orderBy = "id"
	}
	if order == "" {
		order = "desc"
	}

	c := ctx.Request().Context()

	count, err := r.service.Count(c, query, suspendedAt, isOfficial, isProfileCompleted, hasOrganization, hasPackages, affiliateCodes)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	pages := pagination.NewFromRequest(ctx.Request(), int(count))
	users, er := r.service.Query(
		c, pages.Offset(), pages.Limit(),
		orderBy, order, query, suspendedAt, isOfficial, isProfileCompleted, hasOrganization, hasPackages, affiliateCodes,
	)
	if er.StatusCode != 0 {
		return ctx.JSON(er.StatusCode, er)
	}

	result := restypes.QueryResponse{
		Limit:      pages.Limit(),
		Offset:     pages.Offset(),
		Page:       pages.Page,
		TotalRows:  int64(pages.TotalCount),
		TotalPages: pages.PageCount,
		Items:      users,
	}

	return ctx.JSON(http.StatusOK, response.Success(result))
}
