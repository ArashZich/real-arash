package service

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

type UpdateInvoiceRequest struct {
	InvoiceUniqueCode string  `json:"invoice_unique_code"`
	FromName          string  `json:"from_name"`
	FromAddress       string  `json:"from_address"`
	FromPhoneNumber   string  `json:"from_phone_number"`
	FromEmail         string  `json:"from_email"`
	Seller            string  `json:"seller"`
	EconomicID        string  `json:"economic_id"`
	RegisterNumber    string  `json:"register_number"`
	FromPostalCode    string  `json:"from_postal_code"`
	ToName            string  `json:"to_name"`
	ToAddress         string  `json:"to_address"`
	ToPhoneNumber     string  `json:"to_phone_number"`
	ToEmail           string  `json:"to_email"`
	ToPostalCode      string  `json:"to_postal_code"`
	Status            string  `json:"status"`
	SuspendedAt       bool    `json:"suspended_at"`
	TaxPercentage     float32 `json:"tax_amount"`
	OrganizationID    int     `json:"organization_id"`
}

func (c *UpdateInvoiceRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"invoice_unique_code": govalidity.New("invoice_unique_code").Required(),
		"from_name":           govalidity.New("from_name").Optional(),
		"from_address":        govalidity.New("from_address").Optional(),
		"from_phone_number":   govalidity.New("from_phone_number").Optional(),
		"from_email":          govalidity.New("from_email").Optional(),
		"seller":              govalidity.New("seller").Optional(),
		"economic_id":         govalidity.New("economic_id").Optional(),
		"register_number":     govalidity.New("register_number").Optional(),
		"from_postal_code":    govalidity.New("from_postal_code").Optional(),
		"to_name":             govalidity.New("to_name").Optional(),
		"to_address":          govalidity.New("to_address").Optional(),
		"to_phone_number":     govalidity.New("to_phone_number").Optional(),
		"to_email":            govalidity.New("to_email").Optional(),
		"to_postal_code":      govalidity.New("to_postal_code").Optional(),
		"status":              govalidity.New("status").Optional(),
		"suspended_at":        govalidity.New("suspended_at").Optional(),
		"tax_amount":          govalidity.New("tax_amount").Optional(),
		"organization_id":     govalidity.New("organization_id").Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"invoice_unique_code": "شناسه فاکتور",
			"from_name":           "نام فرستنده",
			"from_address":        "آدرس فرستنده",
			"from_phone_number":   "شمار موبایل فرستنده",
			"from_email":          "ایمیل فرستنده",
			"from_postal_code":    "کد پستی فرستنده",
			"seller":              "فروشنده",
			"economic_id":         "شناسه اقتصادی",
			"register_number":     "شماره ثبت",
			"to_name":             "نام گیرنده",
			"to_address":          "آدرس گیرنده",
			"to_phone_number":     "شماره تماس گیرنده",
			"to_email":            "ایمیل گیرنده",
			"to_postal_code":      "کد پستی گیرنده",
			"status":              "وضعیت فاکتور",
			"suspended_at":        "تاریخ تعلیق",
			"tax_amount":          "مالیات",
			"organization_id":     "شناسه سازمان",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *invoice) Update(ctx context.Context, id string, input UpdateInvoiceRequest) (
	models.Invoice, response.ErrorResponse,
) {
	var invoice models.Invoice

	if !policy.CanUpdateInvoice(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ویرایش فاکتور را ندارید")
		return models.Invoice{}, response.ErrorForbidden(nil, "شما دسترسی ویرایش فاکتور را ندارید")
	}

	err := s.db.WithContext(ctx).First(&invoice, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Invoice{}, response.GormErrorResponse(err, "خطایی در یافتن فاکتور رخ داده است")
	}

	invoice.SuspendedAt = dtp.NullTime{
		Time:  time.Now(),
		Valid: input.SuspendedAt,
	}

	invoice.InvoiceUniqueCode = input.InvoiceUniqueCode
	invoice.FromName = input.FromName
	invoice.FromAddress = input.FromAddress
	invoice.FromPhoneNumber = input.FromPhoneNumber
	invoice.FromPhoneNumber = input.FromPhoneNumber
	invoice.FromEmail = input.FromEmail
	invoice.FromPostalCode = input.FromPostalCode
	invoice.Seller = input.Seller
	invoice.EconomicID = input.EconomicID
	invoice.RegisterNumber = input.RegisterNumber
	invoice.ToName = input.ToName
	invoice.ToAddress = input.ToAddress
	invoice.ToPhoneNumber = input.ToPhoneNumber
	invoice.ToEmail = input.ToEmail
	invoice.ToPostalCode = input.ToPostalCode
	invoice.Status = input.Status
	invoice.OrganizationID = input.OrganizationID
	invoice.TaxPercentage = input.TaxPercentage

	err = s.db.WithContext(ctx).Save(&invoice).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Invoice{}, response.GormErrorResponse(err, "خطایی در بروزرسانی فاکتور رخ داده است")
	}
	return invoice, response.ErrorResponse{}
}
