package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/filter"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
)

type OrganizationProductQueryFilterType struct {
	CategoryID       filter.FilterValue[int]    `json:"category_id"`
	OrganizationID   filter.FilterValue[int]    `json:"organization_id"`
	ProductUID       filter.FilterValue[string] `json:"product_uid"`
	OrganizationUID  filter.FilterValue[string] `json:"organization_uid"`
	OrganizationType filter.FilterValue[string] `json:"organization_type"`
}

type OrganizationProductQueryRequestParams struct {
	ID      string                             `json:"id,omitempty"`
	Order   string                             `json:"order,omitempty"`
	OrderBy string                             `json:"order_by,omitempty"`
	Query   string                             `json:"query,omitempty"`
	Filters OrganizationProductQueryFilterType `json:"filters,omitempty"`
}

func (data *OrganizationProductQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":       govalidity.New("id").Optional(),
		"query":    govalidity.New("query").Optional(),
		"order":    govalidity.New("order").Optional(),
		"order_by": govalidity.New("order_by").Optional(),
		"filters": govalidity.Schema{
			"category_id": govalidity.Schema{
				"op":    govalidity.New("filter.category_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.category_id.value").Optional(),
			},
			"organization_id": govalidity.Schema{
				"op":    govalidity.New("filter.organization_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.organization_id.value").Optional(),
			},
			"product_uid": govalidity.Schema{
				"op":    govalidity.New("filter.product_uid.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.product_uid.value").Optional(),
			},
			"organization_uid": govalidity.Schema{
				"op":    govalidity.New("filter.organization_uid.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.organization_uid.value").Optional(),
			},
			"organization_type": govalidity.Schema{
				"op":    govalidity.New("filter.organization_type.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.organization_type.value").Optional(),
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

func makeOrganizationProductFilters(params OrganizationProductQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("products.id = %s", params.ID))
	}

	if params.Filters.CategoryID.Op != "" {
		CategoryIDStr := strconv.Itoa(params.Filters.CategoryID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.CategoryID.Op, CategoryIDStr)
		where = append(where, fmt.Sprintf("products.category_id %s %s",
			opValue.Operator, CategoryIDStr))
	}

	if params.Filters.OrganizationID.Op != "" {
		OrganizationIDStr := strconv.Itoa(params.Filters.OrganizationID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.OrganizationID.Op, OrganizationIDStr)
		where = append(where, fmt.Sprintf("products.organization_id %s %s",
			opValue.Operator, OrganizationIDStr))
	}

	if params.Filters.ProductUID.Op != "" {
		where = append(where, fmt.Sprintf("products.product_uid = '%s'::uuid",
			params.Filters.ProductUID.Value))
	}

	if params.Filters.OrganizationUID.Op != "" {
		where = append(where, fmt.Sprintf("organizations.organization_uid = '%s'::uuid",
			params.Filters.OrganizationUID.Value))
	}

	if params.Filters.OrganizationType.Op != "" {
		// اگر showroom یا recommender درخواست شده، enterprise و admin هم نمایش داده شوند
		if params.Filters.OrganizationType.Value == string(models.OrganizationTypeShowroom) ||
			params.Filters.OrganizationType.Value == string(models.OrganizationTypeRecommender) {
			where = append(where, fmt.Sprintf("(organizations.organization_type IN ('%s', '%s', '%s'))",
				params.Filters.OrganizationType.Value,
				models.OrganizationTypeEnterprise,
				models.OrganizationTypeAdmin))
		} else {
			opValue := filter.GetDBOperatorAndValue(params.Filters.OrganizationType.Op, params.Filters.OrganizationType.Value)
			where = append(where, fmt.Sprintf("organizations.organization_type %s '%s'",
				opValue.Operator, opValue.Value))
		}
	}

	return where
}

func (s *service) OrganizationProduct(
	ctx context.Context, offset, limit int, params OrganizationProductQueryRequestParams,
) ([]models.Product, response.ErrorResponse) {
	var products []models.Product

	// ابتدا سازمان رو پیدا می‌کنیم
	if params.Filters.OrganizationUID.Value != "" {
		var organization models.Organization
		err := s.db.WithContext(ctx).
			Where("organization_uid = ?", params.Filters.OrganizationUID.Value).
			First(&organization).Error

		if err != nil {
			s.logger.With(ctx).Error(err)
			return nil, response.GormErrorResponse(err, "سازمان مورد نظر یافت نشد")
		}

		// اگر basic باشد، هیچ داده‌ای نشان داده نمی‌شود
		if organization.OrganizationType == models.OrganizationTypeBasic {
			return nil, response.GormErrorResponse(err, "برای مشاهده لیست محصولات نیاز به ارتقاء سطح دسترسی دارید")
		}
	}

	tx := s.db.WithContext(ctx).
		Model(&models.Product{}).
		Joins("JOIN organizations ON products.organization_id = organizations.id").
		Offset(offset).
		Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order)).
		Preload("Documents").
		Preload("Category")

	where := makeOrganizationProductFilters(params)
	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	err := tx.Find(&products).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return nil, response.GormErrorResponse(err, "خطایی در یافتن محصول رخ داده است")
	}

	return products, response.ErrorResponse{}
}
