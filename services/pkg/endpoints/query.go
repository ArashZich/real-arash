package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/exp"
	"github.com/ARmo-BigBang/kit/filter"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
)

type PackageQueryFilterType struct {
	ProductLimit   filter.FilterValue[int]    `json:"product_limit"`
	StorageLimitMB filter.FilterValue[int]    `json:"storage_limit_mb"`
	Price          filter.FilterValue[int]    `json:"price"`
	PlanID         filter.FilterValue[int]    `json:"plan_id"`
	CategoryID     filter.FilterValue[int]    `json:"category_id"`
	OrganizationID filter.FilterValue[int]    `json:"organization_id"`
	ExpiredAt      filter.FilterValue[string] `json:"expired_at"`
	CreatedAt      filter.FilterValue[string] `json:"created_at"`
}

type PackageQueryRequestParams struct {
	ID      string                 `json:"id,omitempty"`
	Order   string                 `json:"order,omitempty"`
	OrderBy string                 `json:"order_by,omitempty"`
	Query   string                 `json:"query,omitempty"`
	Filters PackageQueryFilterType `json:"filters,omitempty"`
}

func (data *PackageQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":       govalidity.New("id").Optional(),
		"query":    govalidity.New("query").Optional(),
		"order":    govalidity.New("order").Optional(),
		"order_by": govalidity.New("order_by").Optional(),
		"filters": govalidity.Schema{
			"ProductLimit": govalidity.Schema{
				"op":    govalidity.New("filter.ProductLimit.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.ProductLimit.value").Optional(),
			},
			"StorageLimitMB": govalidity.Schema{
				"op":    govalidity.New("filter.StorageLimitMB.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.StorageLimitMB.value").Optional(),
			},
			"price": govalidity.Schema{
				"op":    govalidity.New("filter.price.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.price.value").Optional(),
			},
			"plan_id": govalidity.Schema{
				"op":    govalidity.New("filter.plan_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.plan_id.value").Optional(),
			},
			"category_id": govalidity.Schema{
				"op":    govalidity.New("filter.category_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.category_id.value").Optional(),
			},
			"organization_id": govalidity.Schema{
				"op":    govalidity.New("filter.organization_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.organization_id.value").Optional(),
			},
			"expired_at": govalidity.Schema{
				"op":    govalidity.New("filter.expired_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.expired_at.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
		},
	}
	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

func makeFilters(params PackageQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("id = %s", params.ID))
	}
	if params.Filters.ProductLimit.Op != "" {
		ProductLimitStr := strconv.Itoa(params.Filters.ProductLimit.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.ProductLimit.Op, ProductLimitStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("product_limit %s %s", opValue.Operator, val))
	}

	if params.Filters.StorageLimitMB.Op != "" {
		StorageLimitMBStr := strconv.Itoa(params.Filters.StorageLimitMB.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.StorageLimitMB.Op, StorageLimitMBStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("storage_limit_mb %s %s", opValue.Operator, val))
	}

	if params.Filters.Price.Op != "" {
		PricesStr := strconv.Itoa(params.Filters.Price.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.Price.Op, PricesStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("price %s %s", opValue.Operator, val))
	}

	if params.Filters.PlanID.Op != "" {
		PlanIDStr := strconv.Itoa(params.Filters.PlanID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.PlanID.Op, PlanIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("plan_id %s %s", opValue.Operator, val))
	}

	if params.Filters.CategoryID.Op != "" {
		CategoryIDStr := strconv.Itoa(params.Filters.CategoryID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.CategoryID.Op, CategoryIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("category_id %s %s", opValue.Operator, val))
	}

	if params.Filters.OrganizationID.Op != "" {
		OrganizationIDStr := strconv.Itoa(params.Filters.OrganizationID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.OrganizationID.Op, OrganizationIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("organization_id %s %s", opValue.Operator, val))
	}

	if params.Filters.ExpiredAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.ExpiredAt.Op, params.Filters.ExpiredAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("expired_at %s %s", opValue.Operator, val))
	}

	if params.Filters.CreatedAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CreatedAt.Op, params.Filters.CreatedAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("created_at %s %s", opValue.Operator, val))
	}

	return where

}

func (s *service) Query(
	ctx context.Context, offset, limit int, params PackageQueryRequestParams,
) ([]models.Package, response.ErrorResponse) {
	var packages []models.Package
	tx := s.db.WithContext(ctx).Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order)).
		Preload("User").Preload("Plan.Categories").Preload("Organization").Preload("Products").Preload("Plan")

	where := makeFilters(params)

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	err := tx.Find(&packages).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Package{}, response.GormErrorResponse(err, "خطایی در یافتن بسته رخ داده است")
	}

	return packages, response.ErrorResponse{}
}
