package endpoints

import (
	"context"
	"strings"

	"net/http"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"

	"fmt"

	"strconv"

	"github.com/ARmo-BigBang/kit/exp"
	"github.com/ARmo-BigBang/kit/filter"
	"github.com/ARmo-BigBang/kit/response"

	"github.com/hoitek-go/govalidity"
)

type ViewGetProductFilterType struct {
	OrganizationID filter.FilterValue[int] `json:"organization_id"` // Add this line

}

type ViewGetProductRequestParams struct {
	ID       string                   `json:"id,omitempty"`
	Order    string                   `json:"order,omitempty"`
	OrderBy  string                   `json:"order_by,omitempty"`
	Query    string                   `json:"query,omitempty"`
	Duration string                   `json:"duration,omitempty"`
	Filters  ViewGetProductFilterType `json:"filters,omitempty"`
}

func (data *ViewGetProductRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":       govalidity.New("id").Optional(),
		"query":    govalidity.New("query").Optional(),
		"order":    govalidity.New("order").Optional(),
		"order_by": govalidity.New("order_by").Optional(),
		"duration": govalidity.New("duration").In([]string{"one_week", "one_month", "three_months", "six_months", "one_year"}).Optional(),
		"filters": govalidity.Schema{
			"organization_id": govalidity.Schema{
				"op":    govalidity.New("filter.organization_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.organization_id.value").Optional(),
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

func makeFiltersProduct(params ViewGetProductRequestParams) []string {
	var where []string

	if params.Filters.OrganizationID.Op != "" {
		OrganizationIDStr := strconv.Itoa(params.Filters.OrganizationID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.OrganizationID.Op, OrganizationIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("organization_id %s %s", opValue.Operator, val))
	}

	return where
}

func (s *service) GetProductView(ctx context.Context, params ViewGetProductRequestParams) ([]models.Product, response.ErrorResponse) {
	if !policy.CanGetProductView(ctx) {
		s.logger.With(ctx).Error("شما دسترسی حذف بازرسی را ندارید")
		return []models.Product{}, response.ErrorForbidden("شما دسترسی حذف بازرسی را ندارید")
	}

	var products []models.Product
	whereClauses := makeFiltersProduct(params) // Generate WHERE clauses
	whereCondition := strings.Join(whereClauses, " AND ")

	query := s.db.WithContext(ctx).
		Model(&models.Product{}).
		Preload("User").
		Preload("Category").
		Preload("Organization").
		Preload("Package").
		Order("view_count desc").
		Limit(5)

	if len(whereClauses) > 0 {
		query = query.Where(whereCondition)
	}

	err := query.Find(&products).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Product{}, response.GormErrorResponse(nil, "Error finding products.")
	}

	return products, response.ErrorResponse{}
}
