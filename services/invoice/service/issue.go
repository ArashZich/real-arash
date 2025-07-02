package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"

	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type InvoiceItemData struct {
	Title           string `json:"title" gorm:"size:100;"`
	Description     string `json:"description"`
	DiscountedPrice int    `json:"discounted_price"`
	TotalPrice      int    `json:"total_price"`
	InvoiceID       int    `json:"invoice_id"`
	OwnerID         int    `json:"owner_id"`
	OwnerType       string `json:"owner_type"`
	OrganizationID  int    `json:"organization_id"`
}

type CreateInvoiceRequest struct {
	InvoiceUniqueCode string            `json:"invoice_unique_code"`
	FromName          string            `json:"from_name"`
	FromAddress       string            `json:"from_address"`
	FromPhoneNumber   string            `json:"from_phone_number"`
	FromEmail         string            `json:"from_email"`
	Seller            string            `json:"seller"`
	EconomicID        string            `json:"economic_id"`
	RegisterNumber    string            `json:"register_number"`
	FromPostalCode    string            `json:"from_postal_code"`
	ToName            string            `json:"to_name"`
	ToAddress         string            `json:"to_address"`
	ToPhoneNumber     string            `json:"to_phone_number"`
	ToEmail           string            `json:"to_email"`
	ToPostalCode      string            `json:"to_postal_code"`
	Status            string            `json:"status"`
	CouponCode        string            `json:"coupon_code"`
	TaxPercentage     float32           `json:"tax_amount"`
	Suspended         bool              `json:"suspended_at"`
	OrganizationID    int               `json:"organization_id"`
	InvoiceItems      []InvoiceItemData `json:"invoice_items"`
}

func (c *CreateInvoiceRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"invoice_unique_code": govalidity.New("invoice_unique_code").Required(),
		"from_name":           govalidity.New("from_name").Required(),
		"from_address":        govalidity.New("from_address").Required(),
		"from_phone_number":   govalidity.New("from_phone_number").Required(),
		"from_email":          govalidity.New("from_email").Required(),
		"seller":              govalidity.New("seller").Required(),
		"economic_id":         govalidity.New("economic_id").Required(),
		"register_number":     govalidity.New("register_number").Required(),
		"from_postal_code":    govalidity.New("from_postal_code").Required(),
		"to_name":             govalidity.New("to_name").Required(),
		"to_address":          govalidity.New("to_address").Required(),
		"to_phone_number":     govalidity.New("to_phone_number").Required(),
		"to_email":            govalidity.New("to_email").Required(),
		"to_postal_code":      govalidity.New("to_postal_code").Required(),
		"status":              govalidity.New("status").Required(),
		"coupon_code":         govalidity.New("coupon_code").Optional(),
		"tax_amount":          govalidity.New("tax_amount").Optional(),
		"suspended_at":        govalidity.New("suspended_at").Optional(),
		"invoice_items":       govalidity.New("invoice_items").Required(),
		"organization_id":     govalidity.New("organization_id").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"invoice_unique_code": "شناسه فاکتور",
			"from_name":           "نام فرستنده",
			"from_address":        "آدرس فرستنده",
			"from_phone_number":   "شمار موبایل فرستنده",
			"from_email":          "ایمیل فرستنده",
			"seller":              "نام فروشنده",
			"economic_id":         "شناسه اقتصاد",
			"register_number":     "شناسه ثبت",
			"from_postal_code":    "کد پستی فرستنده",
			"to_name":             "نام گیرنده",
			"to_address":          "آدرس گیرنده",
			"to_phone_number":     "شماره تماس گیرنده",
			"to_email":            "ایمیل گیرنده",
			"to_postal_code":      "کد پستی گیرنده",
			"status":              "وضعیت فاکتور",
			"coupon_code":         "کد تخفیف",
			"total_amount":        "مبلغ",
			"tax_amount":          "مالیات",
			"suspended_at":        "تاریخ تعلیق",
			"invoice_items":       "آیتم های فاکتور",
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

func (s *invoice) Issue(ctx context.Context, input CreateInvoiceRequest) (models.Invoice, response.ErrorResponse) {

	//! Any person can issue invoice
	// if !policy.CanIssueInvoice(ctx) {
	// 	s.logger.With(ctx).Error("شما دسترسی ایجاد فاکتور را ندارید")
	// 	return models.Invoice{}, response.ErrorForbidden("شما دسترسی ایجاد فاکتور را ندارید")
	// }

	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	var invoiceItems []*models.InvoiceItem
	for _, invoice_item := range input.InvoiceItems {
		invoiceItems = append(
			invoiceItems, &models.InvoiceItem{
				Title:           invoice_item.Title,
				Description:     invoice_item.Description,
				TotalPrice:      invoice_item.TotalPrice,
				DiscountedPrice: invoice_item.DiscountedPrice,
				InvoiceID:       invoice_item.InvoiceID,
				OwnerID:         invoice_item.OwnerID,
				OwnerType:       invoice_item.OwnerType,
				OrganizationID:  invoice_item.OrganizationID,
				UserID:          id,
			},
		)
	}

	invoice := models.Invoice{
		UserID:            id,
		InvoiceUniqueCode: input.InvoiceUniqueCode,
		FromName:          input.FromName,
		FromAddress:       input.FromAddress,
		FromPhoneNumber:   input.FromPhoneNumber,
		FromEmail:         input.FromEmail,
		Seller:            input.Seller,
		EconomicID:        input.EconomicID,
		RegisterNumber:    input.RegisterNumber,
		FromPostalCode:    input.FromPostalCode,
		ToName:            input.ToName,
		ToAddress:         input.ToAddress,
		ToPhoneNumber:     input.ToPhoneNumber,
		ToEmail:           input.ToEmail,
		ToPostalCode:      input.ToPostalCode,
		Status:            input.Status,
		CouponCode:        input.CouponCode,
		TaxPercentage:     input.TaxPercentage,
		OrganizationID:    input.OrganizationID,
		InvoiceItems:      invoiceItems,
		SuspendedAt: dtp.NullTime{
			Time:  time.Now(),
			Valid: input.Suspended,
		},
	}

	err := s.db.WithContext(ctx).Create(&invoice).Preload("InvoiceItems").Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return invoice, response.GormErrorResponse(err, "خطایی در ایجاد فاکتور رخ داد")
	}
	return invoice, response.ErrorResponse{}
}
