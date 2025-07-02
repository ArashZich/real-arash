package service

import (
	"context"
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/cashier"
	"gitag.ir/armogroup/armo/services/reality/cashier/gateway/payping"
	"gitag.ir/armogroup/armo/services/reality/config"
	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/services/invoice/template"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type VerifyInvoiceRequest struct {
	RefId       string `json:"refid"`
	ClientRefId string `json:"clientrefid"`
	CardNumber  string `json:"cardnumber"`
	CardHashPan string `json:"cardhashpan"`
}

func (c *VerifyInvoiceRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"refid":       govalidity.New("refid").Optional(),
		"clientrefid": govalidity.New("clientrefid").Optional(),
		"cardnumber":  govalidity.New("cardnumber").Optional(),
		"cardhashpan": govalidity.New("cardhashpan").Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"refid":       "شناسه فاکتور",
			"clientrefid": "شناسه فاکتور در سامانه مالی",
			"cardnumber":  "شماره کارت",
			"cardhashpan": "شماره رمزی کارت",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)
	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *invoice) Verify(ctx context.Context, input VerifyInvoiceRequest) (string, response.ErrorResponse) {

	targetWebsite := config.AppConfig.PayPingCallbackTargetWebsite

	// Find invoice
	var invoice models.Invoice
	err := s.db.WithContext(ctx).Where("custom_ref_id = ?", input.ClientRefId).Preload("InvoiceItems").First(&invoice).Error
	if err != nil {
		msg := "خطایی در یافتن فاکتور رخ داده است. هزینه کسر شده حد اکثر طی ۷۲ ساعت به حساب شما بازگردانده خواهد شد."
		s.logger.With(ctx).Error(err)
		return template.RenderRedirectTemplate(false, msg, targetWebsite), response.ErrorResponse{}
	}

	// Verify payment
	cashr := cashier.NewCashier("payping")

	rcpt, err := cashr.Verify(&payping.VerifyRequest{PaymentRefID: input.RefId})
	if err != nil {
		s.logger.With(ctx).Errorf("Verification failed for RefID: %s, Error: %v", input.RefId, err)
		msg := "پرداخت شما موفقیت آمیز نبود. اگر مبلغی از حساب شما کسر شده است، ظرف 72 ساعت عودت داده خواهد شد."
		return template.RenderRedirectTemplate(false, msg, targetWebsite), response.ErrorResponse{}
	}
	s.logger.With(ctx).Infof("Verification successful for RefID: %s", input.RefId)

	// Create package after successful payment verification
	item := invoice.InvoiceItems[0]
	_, er := s.CreatePackage(ctx, CreateInvoicePackageRequest{
		PlanID:         item.OwnerID,
		OrganizationID: item.OrganizationID,
	})
	if er.StatusCode != 0 {
		s.logger.With(ctx).Error(er)
		msg := "خطایی در پردازش خرید رخ داده است. هزینه کسر شده حد اکثر طی ۷۲ ساعت به حساب شما بازگردانده خواهد شد."
		return template.RenderRedirectTemplate(false, msg, targetWebsite), response.ErrorResponse{}
	}

	// Update invoice status after successful payment and package creation
	invoice.RefID = rcpt.GetReferenceID()
	invoice.Status = "paid"
	err = s.db.WithContext(ctx).Save(&invoice).Error
	if err != nil {
		s.logger.With(ctx).Error(err, rcpt)
		msg := "خطایی در پردازش فاکتور شما رخ داده است. پلن شما فعال می‌باشد. لطفا با پشتیبانی تماس بگیرید."
		return template.RenderRedirectTemplate(true, msg, targetWebsite), response.ErrorResponse{}
	}

	msg := "فاکتور شما با موفقیت پرداخت شد. پلن شما فعال است."
	return template.RenderRedirectTemplate(true, msg, targetWebsite), response.ErrorResponse{}
}
