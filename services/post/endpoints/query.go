package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/exp"
	"github.com/ARmo-BigBang/kit/filter"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
)

type FilterValueBool struct {
	Op    string `json:"op,omitempty"`
	Value bool   `json:"value,omitempty"`
}

type PostQueryFilterType struct {
	Title     filter.FilterValue[string] `json:"title,omitempty"`
	CreatedAt filter.FilterValue[string] `json:"created_at"`
	Published FilterValueBool            `json:"published"`
}

type PostQueryRequestParams struct {
	ID      string              `json:"id,omitempty"`
	Order   string              `json:"order,omitempty"`
	OrderBy string              `json:"order_by,omitempty"`
	Query   string              `json:"query,omitempty"`
	Filters PostQueryFilterType `json:"filters,omitempty"`
}

func (data *PostQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
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
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
			"published": govalidity.Schema{
				"op":    govalidity.New("filter.published.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.published.value").Optional(),
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

func makeFilters(params PostQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("id = %s", params.ID))
	}

	if params.Filters.Title.Op != "" {
		decodedTitle, err := url.QueryUnescape(params.Filters.Title.Value)
		if err != nil {
			decodedTitle = params.Filters.Title.Value
		}
		opValue := filter.GetDBOperatorAndValue(params.Filters.Title.Op, decodedTitle)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("title %s %s", opValue.Operator, val))
	}

	if params.Filters.CreatedAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CreatedAt.Op, params.Filters.CreatedAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("created_at %s %s", opValue.Operator, val))
	}

	if params.Filters.Published.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Published.Op, fmt.Sprintf("%t", params.Filters.Published.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("%s", opValue.Value))
		where = append(where, fmt.Sprintf("published %s %s", opValue.Operator, val))
	}

	return where
}

func (s *service) Query(
	ctx context.Context, offset, limit int, params PostQueryRequestParams,
) ([]models.Post, response.ErrorResponse) {
	var posts []models.Post
	tx := s.db.WithContext(ctx).Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order))

	where := makeFilters(params)

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	fmt.Println("Querying with filters:", where) // افزودن لاگ برای بررسی فیلترها

	err := tx.Find(&posts).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Post{}, response.GormErrorResponse(err, "خطایی در یافتن پست‌ها رخ داده است")
	}

	fmt.Println("Posts found:", posts) // افزودن لاگ برای بررسی نتایج

	return posts, response.ErrorResponse{}
}
