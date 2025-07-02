package endpoints

import (
	"context"
	"net/http"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type UpdateCouponRequest struct {
	Code                  string    `json:"code"`
	Description           string    `json:"description"`
	DiscountType          string    `json:"discount_type"` //fixed_amount percent
	Status                string    `json:"status"`        /* publish draft pending trash auto_trash //TODO implement auto_trash and pending status */
	DiscountingAmount     int       `json:"discounting_amount"`
	UsageLimit            int       `json:"usage_limit"`
	MaximumDiscountAmount int       `json:"maximum_discount_amount"`
	PlanID                int       `json:"plan_id"`
	ExpireDate            time.Time `json:"expire_date"` // Use string for date input, parse it in the function
	// IndividualUse      bool      `json:"individual_use"`
}

func (c *UpdateCouponRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"code":                    govalidity.New("code").Required(),
		"description":             govalidity.New("description").Optional(),
		"discount_type":           govalidity.New("discount_type").Required(),
		"status":                  govalidity.New("status").Required(),
		"discounting_amount":      govalidity.New("discounting_amount").Required(),
		"usage_limit":             govalidity.New("usage_limit").Optional(),
		"maximum_discount_amount": govalidity.New("maximum_discount_amount").Optional(),
		"plan_id":                 govalidity.New("plan_id").Optional(),
		"expire_date":             govalidity.New("expire_date").Optional(),
		// "individual_use":       govalidity.New("individual_use").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"code":                    "کد",
			"description":             "توضیحات",
			"discount_type":           "نوع تخفیف",
			"status":                  "وضعیت",
			"discounting_amount":      "مقدار تخفیف",
			"usage_limit":             "محدودیت استفاده",
			"maximum_discount_amount": "حداکثر تخفیف",
			"plan_id":                 "شناسه طرح",
			"expire_date":             "تاریخ انقضاء",
			// "individual_use":       "استفاده انفرادی",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Update(ctx context.Context, id string, input UpdateCouponRequest) (
	models.Coupon, response.ErrorResponse,
) {
	var coupon models.Coupon

	if !policy.CanUpdateCoupon(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ویرایش دسته بندی ندارید")
		return models.Coupon{}, response.ErrorForbidden(nil, "شما دسترسی ویرایش دسته بندی ندارید")
	}

	err := s.db.WithContext(ctx).First(&coupon, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Coupon{}, response.GormErrorResponse(err, "خطایی در یافتن دسته رخ داده است")
	}

	expireDate := dtp.NullTime{
		Time:  input.ExpireDate,
		Valid: !input.ExpireDate.IsZero(), // فقط چک می‌کنیم که تاریخ خالی نباشد
	}

	coupon.Code = input.Code
	coupon.Description = input.Description
	coupon.DiscountType = input.DiscountType
	coupon.Status = input.Status
	coupon.DiscountingAmount = input.DiscountingAmount
	coupon.UsageLimit = dtp.NullInt64{
		Int64: int64(input.UsageLimit),
		Valid: input.UsageLimit != 0,
	}
	// coupon.IndividualUse = input.IndividualUse
	coupon.MaximumDiscountAmount = dtp.NullInt64{
		Int64: int64(input.MaximumDiscountAmount),
		Valid: input.MaximumDiscountAmount != 0,
	}
	coupon.PlanID = dtp.NullInt64{
		Int64: int64(input.PlanID),
		Valid: input.PlanID != 0,
	}
	coupon.ExpireDate = expireDate

	err = s.db.WithContext(ctx).Save(&coupon).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Coupon{}, response.GormErrorResponse(err, "خطایی در بروزرسانی دسته رخ داده است")
	}
	return coupon, response.ErrorResponse{}
}
