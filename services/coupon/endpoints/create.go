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

type CreateCouponRequest struct {
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

func (c *CreateCouponRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
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

func (s *service) Create(ctx context.Context, input CreateCouponRequest) (models.Coupon, response.ErrorResponse) {

	if !policy.CanCreateCoupon(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ایجاد کوپن تخفیف را ندارید")
		return models.Coupon{}, response.ErrorForbidden("شما دسترسی ایجاد کوپن تخفیف را ندارید")
	}

	code := ""
	if input.Code != "" {
		if len(input.Code) <= 5 {
			s.logger.With(ctx).Error("طول کد تخفیف باید 6 کاراکتر یا بیشتر باشد")
			return models.Coupon{}, response.ErrorBadRequest("طول کد تخفیف باید 6 کاراکتر یا بیشتر باشد")
		}
		code = input.Code
	} else {
		generateCode, err := generateCode(6, 1, Alphanumeric, PatternChar, "P", "@")
		if err != nil {
			s.logger.With(ctx).Error("خطایی در ساخت کد تخفیف رخ داده است")
			return models.Coupon{}, response.ErrorBadRequest("خطایی در ساخت کد تخفیف رخ داده است")
		}
		code = generateCode
	}

	if input.DiscountType == "percent" && input.DiscountingAmount > 100 {
		s.logger.With(ctx).Error("مبلغ تخفیف نباید بیشتر از 100% باشد")
		return models.Coupon{}, response.ErrorBadRequest("مبلغ تخفیف نباید بیشتر از 100% باشد")
	}

	if input.DiscountType == "fixed_amount" && input.DiscountingAmount <= 0 {
		s.logger.With(ctx).Error("مبلغ تخفیف نباید 0 باشد")
		return models.Coupon{}, response.ErrorBadRequest("مبلغ تخفیف نباید 0 باشد")
	}

	expireDate := dtp.NullTime{
		Time:  input.ExpireDate,
		Valid: input.ExpireDate.After(time.Now()),
	}

	coupon := models.Coupon{
		Code:              code,
		Description:       input.Description,
		DiscountType:      input.DiscountType,
		Status:            input.Status,
		DiscountingAmount: input.DiscountingAmount,
		UsageCount:        0,
		UsageLimit: dtp.NullInt64{
			Int64: int64(input.UsageLimit),
			Valid: input.UsageLimit != 0,
		},
		MaximumDiscountAmount: dtp.NullInt64{
			Int64: int64(input.MaximumDiscountAmount),
			Valid: input.MaximumDiscountAmount != 0,
		},
		PlanID: dtp.NullInt64{
			Int64: int64(input.PlanID),
			Valid: input.PlanID != 0,
		},
		ExpireDate: expireDate,
		// IndividualUse:      input.IndividualUse,
	}

	err := s.db.WithContext(ctx).Create(&coupon).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return coupon, response.GormErrorResponse(err, "خطایی در ایجاد کوپن تخفیف رخ داد")
	}
	return coupon, response.ErrorResponse{}
}
