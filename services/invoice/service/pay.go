package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"gitag.ir/armogroup/armo/services/reality/cashier"
	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/exp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type PayInvoiceRequest struct {
	InvoiceID int    `json:"invoice_id"`
	Gateway   string `json:"gateway"`
}

func (c *PayInvoiceRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"invoice_id": govalidity.New("invoice_id").Required(),
		"gateway":    govalidity.New("gateway").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"invoice_id": "شناسه فاکتور",
			"gateway":    "درگاه",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *invoice) Pay(ctx context.Context, input PayInvoiceRequest) (string, response.ErrorResponse) {
	var invoice models.Invoice
	err := s.db.WithContext(ctx).Preload("InvoiceItems").First(&invoice, "id =?", input.InvoiceID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return "", response.GormErrorResponse(err, "خطایی در یافتن فاکتور رخ داده است")
	}

	totalPrice := 0
	for _, invoiceItem := range invoice.InvoiceItems {
		totalPrice += exp.TerIf(invoiceItem.DiscountedPrice != 0, invoiceItem.DiscountedPrice, invoiceItem.TotalPrice)
	}

	discountedPrice := 0
	var coupon models.Coupon
	if invoice.CouponCode != "" {
		err := s.db.WithContext(ctx).First(&coupon, "code =?", invoice.CouponCode).Error
		if err != nil {
			s.logger.With(ctx).Error("خطایی در یافتن کوپن تخفیف رخ داده است")
			return "", response.ErrorBadRequest(err, "کد تخفیف شما نامعتبر است")
		}

		item := invoice.InvoiceItems[0]
		if coupon.PlanID.Valid && int(coupon.PlanID.Int64) != item.OwnerID {
			s.logger.With(ctx).Error("کد تخفیف شما معتبر نمی باشد")
			return "", response.ErrorBadRequest("کد تخفیف شما معتبر نمی باشد")
		}

		if coupon.Status != "publish" {
			s.logger.With(ctx).Error("کوپن تخفیف نامعتبر است")
			return "", response.ErrorBadRequest("کوپن تخفیف نامعتبر است")
		}

		if coupon.UsageLimit.Valid && int(coupon.UsageLimit.Int64) <= 0 {
			s.logger.With(ctx).Error("کد تخفیف شما منقضی شده است")
			return "", response.ErrorBadRequest("کد تخفیف شما منقضی شده است")
		}

		if coupon.UsageCount >= int(coupon.UsageLimit.Int64) {
			s.logger.With(ctx).Error("تعداد استفاده از کد تخفیف شما به حداکثر رسیده است")
			return "", response.ErrorBadRequest("تعداد استفاده از کد تخفیف شما به حداکثر رسیده است")
		} else {
			coupon.UsageCount = coupon.UsageCount + 1
			err = s.db.WithContext(ctx).Save(&coupon).Error
			if err != nil {
				s.logger.With(ctx).Error("خطایی در بروزرسانی کوپن تخفیف رخ داده است")
				return "", response.GormErrorResponse(err, "خطایی در اعمال کد تخفیف رخ داده است")
			}
		}

		if coupon.ExpireDate.Valid && coupon.ExpireDate.Time.Before(time.Now()) {
			s.logger.With(ctx).Error("کوپن تخفیف منقضی شده است")
			return "", response.ErrorBadRequest("کوپن تخفیف منقضی شده است")
		}

		dprice, errr := discount(coupon.DiscountType, totalPrice, coupon.DiscountingAmount)
		if errr != nil {
			s.logger.With(ctx).Error("خطایی در تخفیف کوپن رخ داده است")
			return "", response.GormErrorResponse(errr, "خطایی در تخفیف کوپن رخ داده است")
		}

		if coupon.DiscountType != "fixed_amount" && coupon.MaximumDiscountAmount.Valid && totalPrice-dprice > int(coupon.MaximumDiscountAmount.Int64) {
			discountedPrice = int(coupon.MaximumDiscountAmount.Int64)
		} else {
			discountedPrice = dprice
		}
	}

	payAmount := exp.TerIf(discountedPrice != 0, discountedPrice, totalPrice)

	// اعمال ۱۰ درصد مالیات بر مبلغ نهایی
	taxAmount := payAmount / 10
	payAmount += taxAmount

	cashr := cashier.NewCashier(input.Gateway)
	cashr.SetAmount(payAmount)
	cashr.SetDescription(invoice.InvoiceItems[0].Description)

	// اضافه کردن جزئیات بیشتر
	cashr.SetDetail("invoice_id", strconv.Itoa(int(invoice.ID)))
	cashr.SetDetail("invoice_unique_code", invoice.InvoiceUniqueCode)
	cashr.SetDetail("title", invoice.InvoiceItems[0].Title)
	cashr.SetDetail("buyer_name", invoice.ToName)
	cashr.SetDetail("buyer_phone", invoice.ToPhoneNumber)
	cashr.SetDetail("buyer_email", invoice.ToEmail)
	cashr.SetDetail("buyer_address", invoice.ToAddress)
	cashr.SetDetail("buyer_postal_code", invoice.ToPostalCode)
	cashr.SetDetail("item_price", strconv.Itoa(invoice.InvoiceItems[0].TotalPrice))
	cashr.SetDetail("discounted_price", strconv.Itoa(discountedPrice))
	cashr.SetDetail("tax_amount", strconv.Itoa(taxAmount))
	cashr.SetDetail("total_amount", strconv.Itoa(payAmount))

	if invoice.CouponCode != "" {
		cashr.SetDetail("coupon_code", invoice.CouponCode)
	}

	err = cashr.Purchase()
	if err != nil {
		s.logger.With(ctx).Error(err)
		return "", response.GormErrorResponse(err, "خطایی در ایجاد درخواست پرداخت رخ داده است")
	}

	invoice.CustomRefID = cashr.GetBill().GetUUID()
	invoice.FinalPaidAmount = payAmount

	err = s.db.WithContext(ctx).Save(&invoice).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return "", response.GormErrorResponse(err, "خطایی در ثبت فاکتور رخ داده است")
	}

	url := cashr.PayURL()

	return url, response.ErrorResponse{}
}
