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

type OrganizationQueryFilterType struct {
	Name                      filter.FilterValue[string] `json:"name,omitempty"`
	Domain                    filter.FilterValue[string] `json:"domain,omitempty"`
	NationalCode              filter.FilterValue[string] `json:"national_code"`
	IsIndividual              filter.FilterValue[string] `json:"is_individual"`
	CompanyName               filter.FilterValue[string] `json:"company_name"`
	CompanyRegistrationNumber filter.FilterValue[string] `json:"company_registration_number"`
	Industry                  filter.FilterValue[string] `json:"industry"`
	CompanySize               filter.FilterValue[int]    `json:"company_size"`
	PhoneNumber               filter.FilterValue[string] `json:"phone_number"`
	Website                   filter.FilterValue[string] `json:"website"`
	CompanyLogo               filter.FilterValue[string] `json:"company_logo"`
	CategoryID                filter.FilterValue[int]    `json:"category_id"`
	CreatedAt                 filter.FilterValue[string] `json:"created_at"`
}

type OrganizationQueryRequestParams struct {
	ID      string                      `json:"id,omitempty"`
	Order   string                      `json:"order,omitempty"`
	OrderBy string                      `json:"order_by,omitempty"`
	Query   string                      `json:"query,omitempty"`
	Filters OrganizationQueryFilterType `json:"filters,omitempty"`
}

func (data *OrganizationQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
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
			"domain": govalidity.Schema{
				"op":    govalidity.New("filter.domain.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.domain.value").Optional(),
			},
			"national_code": govalidity.Schema{
				"op":    govalidity.New("filter.national_code.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.national_code.value").Optional(),
			},
			"is_individual": govalidity.Schema{
				"op":    govalidity.New("filter.is_individual.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.is_individual.value").Optional(),
			},
			"company_name": govalidity.Schema{
				"op":    govalidity.New("filter.company_name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.company_name.value").Optional(),
			},
			"company_registration_number": govalidity.Schema{
				"op":    govalidity.New("filter.company_registration_number.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.company_registration_number.value").Optional(),
			},
			"industry": govalidity.Schema{
				"op":    govalidity.New("filter.industry.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.industry.value").Optional(),
			},
			"company_size": govalidity.Schema{
				"op":    govalidity.New("filter.company_size.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.company_size.value").Optional(),
			},
			"phone_number": govalidity.Schema{
				"op":    govalidity.New("filter.phone_number.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.phone_number.value").Optional(),
			},
			"website": govalidity.Schema{
				"op":    govalidity.New("filter.website.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.website.value").Optional(),
			},
			"company_logo": govalidity.Schema{
				"op":    govalidity.New("filter.company_logo.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.company_logo.value").Optional(),
			},
			"category_id": govalidity.Schema{
				"op":    govalidity.New("filter.category_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.category_id.value").Optional(),
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

func makeFilters(params OrganizationQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("id = %s", params.ID))
	}

	if params.Filters.Name.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Name.Op, params.Filters.Name.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("name %s %s", opValue.Operator, val))
	}
	if params.Filters.Domain.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Domain.Op, params.Filters.Domain.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("domain %s %s", opValue.Operator, val))
	}
	if params.Filters.Industry.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Industry.Op, params.Filters.Industry.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("industry %s %s", opValue.Operator, val))
	}
	if params.Filters.NationalCode.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.NationalCode.Op, params.Filters.NationalCode.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("national_code %s %s", opValue.Operator, val))
	}
	if params.Filters.IsIndividual.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.IsIndividual.Op, params.Filters.IsIndividual.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("is_individual %s %s", opValue.Operator, val))
	}
	if params.Filters.CompanyName.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CompanyName.Op, params.Filters.CompanyName.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("company_name %s %s", opValue.Operator, val))
	}
	if params.Filters.CompanyRegistrationNumber.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CompanyRegistrationNumber.Op, params.Filters.CompanyRegistrationNumber.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("company_registration_number %s %s", opValue.Operator, val))
	}
	if params.Filters.CompanySize.Op != "" {
		CompanySizeStr := strconv.Itoa(params.Filters.CompanySize.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.CompanySize.Op, CompanySizeStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("company_size %s %s", opValue.Operator, val))
	}

	if params.Filters.PhoneNumber.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.PhoneNumber.Op, params.Filters.PhoneNumber.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("phone_number %s %s", opValue.Operator, val))
	}

	if params.Filters.Website.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Website.Op, params.Filters.Website.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("website %s %s", opValue.Operator, val))
	}
	if params.Filters.CompanyLogo.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CompanyLogo.Op, params.Filters.CompanyLogo.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("company_logo %s %s", opValue.Operator, val))
	}

	if params.Filters.CategoryID.Op != "" {
		CategoryIDStr := strconv.Itoa(params.Filters.CategoryID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.CategoryID.Op, CategoryIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("category_id %s %s", opValue.Operator, val))
	}

	if params.Filters.CreatedAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CreatedAt.Op, params.Filters.CreatedAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("created_at %s %s", opValue.Operator, val))
	}

	return where

}

func (s *service) Query(
	ctx context.Context, offset, limit int, params OrganizationQueryRequestParams,
) ([]models.Organization, response.ErrorResponse) {
	var organizations []models.Organization
	tx := s.db.WithContext(ctx).Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order)).
		Preload("User").
		Preload("Products").
		Preload("Packages")

	where := makeFilters(params)

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	err := tx.Find(&organizations).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Organization{}, response.GormErrorResponse(err, "خطایی در یافتن سازمان رخ داده است")
	}

	return organizations, response.ErrorResponse{}
}
