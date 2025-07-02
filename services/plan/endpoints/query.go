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

type PlanQueryFilterType struct {
	Title           filter.FilterValue[string] `json:"title,omitempty"`
	Description     filter.FilterValue[string] `json:"description"`
	Price           filter.FilterValue[int]    `json:"price"`
	DiscountedPrice filter.FilterValue[int]    `json:"discounted_price"`
	CategoryID      filter.FilterValue[int]    `json:"category_id"`
	DayLength       filter.FilterValue[int]    `json:"day_length"`
	ProductLimit    filter.FilterValue[int]    `json:"product_limit"`
	StorageLimitMB  filter.FilterValue[int]    `json:"storage_limit_mb"`
	IconUrl         filter.FilterValue[string] `json:"icon_url"`
	CreatedAt       filter.FilterValue[string] `json:"created_at"`
}

type PlanQueryRequestParams struct {
	ID      string              `json:"id,omitempty"`
	Order   string              `json:"order,omitempty"`
	OrderBy string              `json:"order_by,omitempty"`
	Query   string              `json:"query,omitempty"`
	Filters PlanQueryFilterType `json:"filters,omitempty"`
}

func (data *PlanQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":       govalidity.New("id").Optional(),
		"query":    govalidity.New("query").Optional(),
		"order":    govalidity.New("order").Optional(),
		"order_by": govalidity.New("order_by").Optional(),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"description": govalidity.Schema{
				"op":    govalidity.New("filter.description.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.description.value").Optional(),
			},
			"price": govalidity.Schema{
				"op":    govalidity.New("filter.price.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.price.value").Optional(),
			},
			"discounted_price": govalidity.Schema{
				"op":    govalidity.New("filter.discounted_price.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.discounted_price.value").Optional(),
			},
			"category_id": govalidity.Schema{
				"op":    govalidity.New("filter.category_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.category_id.value").Optional(),
			},
			"day_length": govalidity.Schema{
				"op":    govalidity.New("filter.day_length.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.day_length.value").Optional(),
			},
			"product_limit": govalidity.Schema{
				"op":    govalidity.New("filter.product_limit.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.product_limit.value").Optional(),
			},
			"storage_limit_mb": govalidity.Schema{
				"op":    govalidity.New("filter.storage_limit_mb.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.storage_limit_mb.value").Optional(),
			},
			"icon_url": govalidity.Schema{
				"op":    govalidity.New("filter.icon_url.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.icon_url.value").Optional(),
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

func makeFilters(params PlanQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("id = %s", params.ID))
	}

	if params.Filters.Title.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Title.Op, params.Filters.Title.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("title %s %s", opValue.Operator, val))
	}
	if params.Filters.Description.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Description.Op, params.Filters.Description.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("description %s %s", opValue.Operator, val))
	}
	if params.Filters.Price.Op != "" {
		PriceStr := strconv.Itoa(params.Filters.Price.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.Price.Op, PriceStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("price %s %s", opValue.Operator, val))
	}
	if params.Filters.DiscountedPrice.Op != "" {
		DiscountedPriceStr := strconv.Itoa(params.Filters.DiscountedPrice.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.DiscountedPrice.Op, DiscountedPriceStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("discounted_price %s %s", opValue.Operator, val))
	}
	if params.Filters.CategoryID.Op != "" {
		CategoryIDStr := strconv.Itoa(params.Filters.CategoryID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.CategoryID.Op, CategoryIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where,
			fmt.Sprintf("EXISTS (SELECT 1 FROM plan_categories WHERE plan_categories.plan_id = plans.ID AND plan_categories.category_id %s %s)",
				opValue.Operator, val))
	}
	if params.Filters.DayLength.Op != "" {
		DayLengthStr := strconv.Itoa(params.Filters.DayLength.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.DayLength.Op, DayLengthStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("day_length %s %s", opValue.Operator, val))
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
	if params.Filters.IconUrl.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.IconUrl.Op, params.Filters.IconUrl.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("icon_url %s %s", opValue.Operator, val))
	}

	if params.Filters.CreatedAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CreatedAt.Op, params.Filters.CreatedAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("created_at %s %s", opValue.Operator, val))
	}

	return where

}

func (s *service) Query(
	ctx context.Context, offset, limit int, params PlanQueryRequestParams,
) ([]models.Plan, response.ErrorResponse) {
	var plans []models.Plan
	tx := s.db.WithContext(ctx).Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order)).
		Preload("Categories").
		Preload("Categories.Parent.Parent.Parent").
		Preload("Categories.Children.Children.Children")
		// .Preload("User")
	where := makeFilters(params)

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	err := tx.Find(&plans).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Plan{}, response.GormErrorResponse(err, "خطایی در یافتن طرح رخ داده است")
	}

	return plans, response.ErrorResponse{}
}
