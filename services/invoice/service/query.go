package service

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

type InvoiceQueryFilterType struct {
	InvoiceUniqueCode filter.FilterValue[string] `json:"invoice_unique_code"`
	FromName          filter.FilterValue[string] `json:"from_name"`
	FromAddress       filter.FilterValue[string] `json:"from_address"`
	FromPhoneNumber   filter.FilterValue[string] `json:"from_phone_number"`
	FromEmail         filter.FilterValue[string] `json:"from_email"`
	Seller            filter.FilterValue[string] `json:"seller"`
	EconomicID        filter.FilterValue[string] `json:"economic_id"`
	RegisterNumber    filter.FilterValue[string] `json:"register_number"`
	FromPostalCode    filter.FilterValue[string] `json:"from_postal_code"`
	ToName            filter.FilterValue[string] `json:"to_name"`
	ToAddress         filter.FilterValue[string] `json:"to_address"`
	ToPhoneNumber     filter.FilterValue[string] `json:"to_phone_number"`
	ToEmail           filter.FilterValue[string] `json:"to_email"`
	ToPostalCode      filter.FilterValue[string] `json:"to_postal_code"`
	Status            filter.FilterValue[string] `json:"status"`
	TotalAmount       filter.FilterValue[int]    `json:"total_amount"`
	DiscountAmount    filter.FilterValue[int]    `json:"discount_amount"`
	CouponID          filter.FilterValue[int]    `json:"coupon_id"`
	TaxPercentage     filter.FilterValue[string] `json:"tax_amount"`
	OrganizationID    filter.FilterValue[int]    `json:"organization_id"`
	SuspendedAt       filter.FilterValue[string] `json:"suspended_at"`
	CreatedAt         filter.FilterValue[string] `json:"created_at"`
}

type InvoiceQueryRequestParams struct {
	ID      string                 `json:"id,omitempty"`
	Order   string                 `json:"order,omitempty"`
	OrderBy string                 `json:"order_by,omitempty"`
	Query   string                 `json:"query,omitempty"`
	Filters InvoiceQueryFilterType `json:"filters,omitempty"`
}

func (data *InvoiceQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":       govalidity.New("id").Optional(),
		"query":    govalidity.New("query").Optional(),
		"order":    govalidity.New("order").Optional(),
		"order_by": govalidity.New("order_by").Optional(),
		"filters": govalidity.Schema{
			"organization_id": govalidity.Schema{
				"op":    govalidity.New("filter.organization_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.organization_id.value").Optional(),
			},
			"invoice_unique_code": govalidity.Schema{
				"op":    govalidity.New("filter.invoice_unique_code.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.invoice_unique_code.value").Optional(),
			},
			"FromName": govalidity.Schema{
				"op":    govalidity.New("filter.FromName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.FromName.value").Optional(),
			},
			"FromAddress": govalidity.Schema{
				"op":    govalidity.New("filter.FromAddress.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.FromAddress.value").Optional(),
			},
			"FromPhoneNumber": govalidity.Schema{
				"op":    govalidity.New("filter.FromPhoneNumber.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.FromPhoneNumber.value").Optional(),
			},
			"FromEmail": govalidity.Schema{
				"op":    govalidity.New("filter.FromEmail.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.FromEmail.value").Optional(),
			},
			"Seller": govalidity.Schema{
				"op":    govalidity.New("filter.Seller.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.Seller.value").Optional(),
			},
			"EconomicID": govalidity.Schema{
				"op":    govalidity.New("filter.EconomicID.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.EconomicID.value").Optional(),
			},
			"RegisterNumber": govalidity.Schema{
				"op":    govalidity.New("filter.RegisterNumber.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.RegisterNumber.value").Optional(),
			},
			"FromPostalCode": govalidity.Schema{
				"op":    govalidity.New("filter.FromPostalCode.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.FromPostalCode.value").Optional(),
			},
			"ToName": govalidity.Schema{
				"op":    govalidity.New("filter.ToName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.ToName.value").Optional(),
			},
			"ToAddress": govalidity.Schema{
				"op":    govalidity.New("filter.ToAddress.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.ToAddress.value").Optional(),
			},
			"ToPhoneNumber": govalidity.Schema{
				"op":    govalidity.New("filter.ToPhoneNumber.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.ToPhoneNumber.value").Optional(),
			},
			"ToEmail": govalidity.Schema{
				"op":    govalidity.New("filter.ToEmail.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.ToEmail.value").Optional(),
			},
			"ToPostalCode": govalidity.Schema{
				"op":    govalidity.New("filter.ToPostalCode.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.ToPostalCode.value").Optional(),
			},
			"Status": govalidity.Schema{
				"op":    govalidity.New("filter.Status.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.Status.value").Optional(),
			},
			"DiscountAmount": govalidity.Schema{
				"op":    govalidity.New("filter.DiscountAmount.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.DiscountAmount.value").Optional(),
			},
			"total_amount": govalidity.Schema{
				"op":    govalidity.New("filter.total_amount.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.total_amount.value").Optional(),
			},
			"coupon_id": govalidity.Schema{
				"op":    govalidity.New("filter.coupon_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.coupon_id.value").Optional(),
			},
			"tax_amount": govalidity.Schema{
				"op":    govalidity.New("filter.tax_amount.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.tax_amount.value").Optional(),
			},
			"suspended_at": govalidity.Schema{
				"op":    govalidity.New("filter.suspended_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.suspended_at.value").Optional(),
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

func makeFilters(params InvoiceQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("id = %s", params.ID))
	}

	if params.Filters.OrganizationID.Op != "" {
		OrganizationIDStr := strconv.Itoa(params.Filters.OrganizationID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.OrganizationID.Op, OrganizationIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("organization_id %s %s", opValue.Operator, val))
	}

	if params.Filters.InvoiceUniqueCode.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.InvoiceUniqueCode.Op, params.Filters.InvoiceUniqueCode.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("invoice_unique_code %s %s", opValue.Operator, val))
	}
	if params.Filters.FromName.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.FromName.Op, params.Filters.FromName.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("from_name %s %s", opValue.Operator, val))
	}
	if params.Filters.FromAddress.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.FromAddress.Op, params.Filters.FromAddress.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("from_address %s %s", opValue.Operator, val))
	}
	if params.Filters.FromPhoneNumber.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.FromPhoneNumber.Op, params.Filters.FromPhoneNumber.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("from_email %s %s", opValue.Operator, val))
	}
	if params.Filters.FromEmail.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.FromEmail.Op, params.Filters.FromEmail.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("from_email %s %s", opValue.Operator, val))
	}
	if params.Filters.ToName.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.ToName.Op, params.Filters.ToName.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("to_name %s %s", opValue.Operator, val))
	}
	if params.Filters.ToAddress.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.ToAddress.Op, params.Filters.ToAddress.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("to_address %s %s", opValue.Operator, val))
	}
	if params.Filters.ToPhoneNumber.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.ToPhoneNumber.Op, params.Filters.ToPhoneNumber.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("to_phone_number %s %s", opValue.Operator, val))
	}
	if params.Filters.ToEmail.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.ToEmail.Op, params.Filters.ToEmail.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("to_email %s %s", opValue.Operator, val))
	}
	if params.Filters.ToPostalCode.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.ToPostalCode.Op, params.Filters.ToPostalCode.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("to_postal_code %s %s", opValue.Operator, val))
	}
	if params.Filters.Status.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Status.Op, params.Filters.Status.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("status %s %s", opValue.Operator, val))
	}
	if params.Filters.DiscountAmount.Op != "" {
		DiscountAmountStr := strconv.Itoa(params.Filters.DiscountAmount.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.DiscountAmount.Op, DiscountAmountStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("discount_amount %s %s", opValue.Operator, val))
	}
	if params.Filters.TotalAmount.Op != "" {
		TotalAmountStr := strconv.Itoa(params.Filters.TotalAmount.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.TotalAmount.Op, TotalAmountStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("total_amount %s %s", opValue.Operator, val))
	}
	if params.Filters.CouponID.Op != "" {
		CouponIDStr := strconv.Itoa(params.Filters.CouponID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.CouponID.Op, CouponIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("coupon_id %s %s", opValue.Operator, val))
	}
	if params.Filters.TaxPercentage.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.DiscountAmount.Op, params.Filters.TaxPercentage.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("tax_amount %s %s", opValue.Operator, val))
	}
	if params.Filters.SuspendedAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.SuspendedAt.Op, params.Filters.SuspendedAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("suspended_at %s %s", opValue.Operator, val))
	}
	if params.Filters.CreatedAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CreatedAt.Op, params.Filters.CreatedAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("created_at %s %s", opValue.Operator, val))
	}

	return where

}

func (s *invoice) Query(
	ctx context.Context, offset, limit int, params InvoiceQueryRequestParams,
) ([]models.Invoice, response.ErrorResponse) {
	var invoices []models.Invoice
	tx := s.db.WithContext(ctx).Offset(offset).Preload("InvoiceItems").Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order)).
		Preload("User")

	where := makeFilters(params)

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	// if user not admin add where invoice paid

	err := tx.Find(&invoices).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Invoice{}, response.GormErrorResponse(err, "خطایی در یافتن فاکتور رخ داده است")
	}

	return invoices, response.ErrorResponse{}
}
