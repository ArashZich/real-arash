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

type CouponQueryFilterType struct {
	Code                  filter.FilterValue[string] `json:"code,omitempty"`
	DiscountType          filter.FilterValue[string] `json:"discount_type"`
	Status                filter.FilterValue[string] `json:"status"`
	DiscountAmount        filter.FilterValue[string] `json:"Discounting_amount"`
	UsageCount            filter.FilterValue[int]    `json:"usage_count"`
	UsageLimit            filter.FilterValue[int]    `json:"usage_limit"`
	MaximumDiscountAmount filter.FilterValue[int]    `json:"maximum_discount_amount"`
	PlanID                filter.FilterValue[int]    `json:"plan_id"`
	ExpireDate            filter.FilterValue[string] `json:"expire_date"`
	CreatedAt             filter.FilterValue[string] `json:"created_at"`
}

type CouponQueryRequestParams struct {
	ID      string                `json:"id,omitempty"`
	Order   string                `json:"order,omitempty"`
	OrderBy string                `json:"order_by,omitempty"`
	Query   string                `json:"query,omitempty"`
	Filters CouponQueryFilterType `json:"filters,omitempty"`
}

func (data *CouponQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":       govalidity.New("id").Optional(),
		"query":    govalidity.New("query").Optional(),
		"order":    govalidity.New("order").Optional(),
		"order_by": govalidity.New("order_by").Optional(),
		"filters": govalidity.Schema{
			"code": govalidity.Schema{
				"op":    govalidity.New("filter.code.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.code.value").Optional(),
			},
			"discount_type": govalidity.Schema{
				"op":    govalidity.New("filter.discount_type.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.discount_type.value").Optional(),
			},
			"status": govalidity.Schema{
				"op":    govalidity.New("filter.status.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.status.value").Optional(),
			},
			"Discounting_amount": govalidity.Schema{
				"op":    govalidity.New("filter.Discounting_amount.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.Discounting_amount.value").Optional(),
			},
			"usage_count": govalidity.Schema{
				"op":    govalidity.New("filter.usage_count.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.usage_count.value").Optional(),
			},
			"usage_limit": govalidity.Schema{
				"op":    govalidity.New("filter.usage_limit.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.usage_limit.value").Optional(),
			},
			"plan_id": govalidity.Schema{
				"op":    govalidity.New("filter.plan_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.plan_id.value").Optional(),
			},
			"maximum_discount_amount": govalidity.Schema{
				"op":    govalidity.New("filter.maximum_discount_amount.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.maximum_discount_amount.value").Optional(),
			},
			"expire_date": govalidity.Schema{
				"op":    govalidity.New("filter.expire_date.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.expire_date.value").Optional(),
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

func makeFilters(params CouponQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("id = %s", params.ID))
	}

	if params.Filters.Code.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Code.Op, params.Filters.Code.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("code %s %s", opValue.Operator, val))
	}
	if params.Filters.DiscountType.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.DiscountType.Op, params.Filters.DiscountType.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("discount_type %s %s", opValue.Operator, val))
	}
	if params.Filters.Status.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Status.Op, params.Filters.Status.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("status %s %s", opValue.Operator, val))
	}
	if params.Filters.DiscountAmount.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.DiscountAmount.Op, params.Filters.DiscountAmount.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("Discounting_amount %s %s", opValue.Operator, val))
	}
	if params.Filters.UsageCount.Op != "" {
		UsageCountStr := strconv.Itoa(params.Filters.UsageCount.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.UsageCount.Op, UsageCountStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("usage_count %s %s", opValue.Operator, val))
	}
	if params.Filters.UsageLimit.Op != "" {
		UsageLimitStr := strconv.Itoa(params.Filters.UsageLimit.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.UsageLimit.Op, UsageLimitStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("usage_limit %s %s", opValue.Operator, val))
	}
	if params.Filters.MaximumDiscountAmount.Op != "" {
		MaximumDiscountAmountStr := strconv.Itoa(params.Filters.MaximumDiscountAmount.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.MaximumDiscountAmount.Op, MaximumDiscountAmountStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("maximum_amount %s %s", opValue.Operator, val))
	}
	if params.Filters.ExpireDate.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.ExpireDate.Op, params.Filters.ExpireDate.Value)
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
	ctx context.Context, offset, limit int, params CouponQueryRequestParams,
) ([]models.Coupon, response.ErrorResponse) {
	var coupons []models.Coupon
	tx := s.db.WithContext(ctx).Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order)).
		Preload("Plan")

	where := makeFilters(params)

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	err := tx.Find(&coupons).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Coupon{}, response.GormErrorResponse(err, "خطایی در یافتن دسته رخ داده است")
	}

	return coupons, response.ErrorResponse{}
}
