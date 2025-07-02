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

type ProductQueryFilterType struct {
	Name           filter.FilterValue[string] `json:"name,omitempty"`
	ThumbnailURI   filter.FilterValue[string] `json:"thumbnail_uri"`
	CategoryID     filter.FilterValue[int]    `json:"category_id"`
	PackageID      filter.FilterValue[int]    `json:"package_id"`
	OrganizationID filter.FilterValue[int]    `json:"organization_id"`
	DisabledAt     filter.FilterValue[string] `json:"disabled_at"`
	CreatedAt      filter.FilterValue[string] `json:"created_at"`
	ProductUID     filter.FilterValue[string] `json:"product_uid"`
}

type ProductQueryRequestParams struct {
	ID      string                 `json:"id,omitempty"`
	Order   string                 `json:"order,omitempty"`
	OrderBy string                 `json:"order_by,omitempty"`
	Query   string                 `json:"query,omitempty"`
	Filters ProductQueryFilterType `json:"filters,omitempty"`
}

func (data *ProductQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":       govalidity.New("id").Optional(),
		"query":    govalidity.New("query").Optional(),
		"order":    govalidity.New("order").Optional(),
		"order_by": govalidity.New("order_by").Optional(),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"thumbnail_url": govalidity.Schema{
				"op":    govalidity.New("filter.thumbnail_uri.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.thumbnail_uri.value").Optional(),
			},
			"category_id": govalidity.Schema{
				"op":    govalidity.New("filter.category_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.category_id.value").Optional(),
			},
			"organization_id": govalidity.Schema{
				"op":    govalidity.New("filter.organization_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.organization_id.value").Optional(),
			},
			"package_id": govalidity.Schema{
				"op":    govalidity.New("filter.package_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.package_id.value").Optional(),
			},
			"disabled_at": govalidity.Schema{
				"op":    govalidity.New("filter.disabled_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.disabled_at.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
			"product_uid": govalidity.Schema{
				"op":    govalidity.New("filter.product_uid.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.product_uid.value").Optional(),
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

func makeFilters(params ProductQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("id = %s", params.ID))
	}

	if params.Filters.Name.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Name.Op, params.Filters.Name.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("name %s %s", opValue.Operator, val))
	}
	if params.Filters.ThumbnailURI.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.ThumbnailURI.Op, params.Filters.ThumbnailURI.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("thumbnail_uri %s %s", opValue.Operator, val))
	}
	if params.Filters.CategoryID.Op != "" {
		CategoryIDStr := strconv.Itoa(params.Filters.CategoryID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.CategoryID.Op, CategoryIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("category_id %s %s", opValue.Operator, val))
	}
	if params.Filters.PackageID.Op != "" {
		PackageIDStr := strconv.Itoa(params.Filters.PackageID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.PackageID.Op, PackageIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("package_id %s %s", opValue.Operator, val))
	}
	if params.Filters.OrganizationID.Op != "" {
		OrganizationIDStr := strconv.Itoa(params.Filters.OrganizationID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.OrganizationID.Op, OrganizationIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("organization_id %s %s", opValue.Operator, val))
	}
	if params.Filters.DisabledAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.DisabledAt.Op, params.Filters.DisabledAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("disabled_at %s %s", opValue.Operator, val))
	}
	if params.Filters.CreatedAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CreatedAt.Op, params.Filters.CreatedAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("created_at %s %s", opValue.Operator, val))
	}
	if params.Filters.ProductUID.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.ProductUID.Op, params.Filters.ProductUID.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("product_uid %s %s", opValue.Operator, val))
	}

	return where
}

func (s *service) Query(
	ctx context.Context, offset, limit int, params ProductQueryRequestParams,
) ([]models.Product, response.ErrorResponse) {
	var products []models.Product
	tx := s.db.WithContext(ctx).Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order)).
		Preload("User").Preload("Documents").Preload("Category") // اصلاح ارجاع به Documents

	where := makeFilters(params)

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	err := tx.Find(&products).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Product{}, response.GormErrorResponse(err, "خطایی در یافتن محصول رخ داده است")
	}

	return products, response.ErrorResponse{}
}
