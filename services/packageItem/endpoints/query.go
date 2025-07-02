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

type PackageItemQueryFilterType struct {
	Title     filter.FilterValue[string] `json:"title,omitempty"`
	Price     filter.FilterValue[int]    `json:"price"`
	CreatedAt filter.FilterValue[string] `json:"created_at"`
}

type PackageItemQueryRequestParams struct {
	ID      string                     `json:"id,omitempty"`
	Order   string                     `json:"order,omitempty"`
	OrderBy string                     `json:"order_by,omitempty"`
	Query   string                     `json:"query,omitempty"`
	Filters PackageItemQueryFilterType `json:"filters,omitempty"`
}

func (data *PackageItemQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
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
			"price": govalidity.Schema{
				"op":    govalidity.New("filter.price.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.price.value").Optional(),
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

func makeFilters(params PackageItemQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("id = %s", params.ID))
	}

	if params.Filters.Title.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Title.Op, params.Filters.Title.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("title %s %s", opValue.Operator, val))
	}

	if params.Filters.Price.Op != "" {
		PriceStr := strconv.Itoa(params.Filters.Price.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.Price.Op, PriceStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("storage_limit_mb %s %s", opValue.Operator, val))
	}

	if params.Filters.CreatedAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CreatedAt.Op, params.Filters.CreatedAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("created_at %s %s", opValue.Operator, val))
	}

	return where

}

func (s *packageItem) Query(
	ctx context.Context, offset, limit int, params PackageItemQueryRequestParams,
) ([]models.PackageItem, response.ErrorResponse) {
	var packageItems []models.PackageItem
	tx := s.db.WithContext(ctx).Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order)).
		Preload("User")

	where := makeFilters(params)

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	err := tx.Find(&packageItems).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.PackageItem{}, response.GormErrorResponse(err, "خطایی در یافتن آیتم رخ داده است")
	}

	return packageItems, response.ErrorResponse{}
}
