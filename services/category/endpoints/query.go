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

type CategoryQueryFilterType struct {
	Title     filter.FilterValue[string] `json:"title,omitempty"`
	ParentID  filter.FilterValue[int]    `json:"parent_id"`
	IconUrl   filter.FilterValue[string] `json:"icon_url"`
	CreatedAt filter.FilterValue[string] `json:"created_at"`
}

type CategoryQueryRequestParams struct {
	ID      string                  `json:"id,omitempty"`
	Order   string                  `json:"order,omitempty"`
	OrderBy string                  `json:"order_by,omitempty"`
	Query   string                  `json:"query,omitempty"`
	Filters CategoryQueryFilterType `json:"filters,omitempty"`
}

func (data *CategoryQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
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
			"parent_id": govalidity.Schema{
				"op":    govalidity.New("filter.parent_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.parent_id.value").Optional(),
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

func makeFilters(params CategoryQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("id = %s", params.ID))
	}

	if params.Filters.Title.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Title.Op, params.Filters.Title.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("title %s %s", opValue.Operator, val))
	}

	if params.Filters.ParentID.Op != "" {
		ParentIDStr := strconv.Itoa(params.Filters.ParentID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.ParentID.Op, ParentIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("parent_id %s %s", opValue.Operator, val))
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
	ctx context.Context, offset, limit int, params CategoryQueryRequestParams,
) ([]models.Category, response.ErrorResponse) {
	var categories []models.Category
	tx := s.db.WithContext(ctx).Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order)).
		Preload("User").
		Preload("Parent.Parent.Parent").
		Preload("Children.Children.Children").
		Preload("Plans").
		Preload("Products")

	where := makeFilters(params)

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	err := tx.Find(&categories).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Category{}, response.GormErrorResponse(err, "خطایی در یافتن دسته رخ داده است")
	}

	return categories, response.ErrorResponse{}
}
