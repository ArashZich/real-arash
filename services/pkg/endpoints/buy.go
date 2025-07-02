package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	invoiceSvc "gitag.ir/armogroup/armo/services/reality/services/invoice/service"
	"github.com/ARmo-BigBang/kit/exp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/google/uuid"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type BuyPackageRequest struct {
	PlanID         int    `json:"plan_id"`
	CouponCode     string `json:"coupon_code"`
	OrganizationID int    `json:"organization_id"`
}

func (c *BuyPackageRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"plan_id":         govalidity.New("plan_id").Required(),
		"coupon_code":     govalidity.New("coupon_code").Optional(),
		"organization_id": govalidity.New("organization_id").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"plan_id":         "شناسه طرح",
			"coupon_id":       "کد تخفیف",
			"organization_id": "سازمان",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Buy(ctx context.Context, input BuyPackageRequest) (string, response.ErrorResponse) {
	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	// find plan from plan id
	var plan models.Plan
	err := s.db.WithContext(ctx).Preload("Categories").First(&plan, "id =?", input.PlanID).Error
	if err != nil {
		s.logger.With(ctx).Error("Failed to find plan", "PlanID", input.PlanID, "error", err)
		return "", response.GormErrorResponse(err, "خطایی در یافتن طرح رخ داده است")
	}

	var organization models.Organization
	err = s.db.WithContext(ctx).Preload("Category").Preload("User").First(&organization, "id =?", input.OrganizationID).Error
	if err != nil {
		s.logger.With(ctx).Error("Failed to find organization", "OrganizationID", input.OrganizationID, "error", err)
		return "", response.GormErrorResponse(err, "خطایی در یافتن سازمان رخ داده است")
	}

	// Check if organization already has an active package
	var pkg models.Package
	exists := false
	err = s.db.WithContext(ctx).Preload("Plan").Where("user_id = ? AND organization_id = ?", id, input.OrganizationID).First(&pkg).Error
	if err == nil {
		exists = true
		s.logger.With(ctx).Info("Existing package found", "PackageID", pkg.ID, "PlanID", pkg.PlanID)
	}

	planHierarchy := map[string]int{
		"starter":    1,
		"pro":        2,
		"premium":    3,
		"enterprise": 4,
	}

	// Fetch the current (or "old") plan details if the package already exists
	var currentPlan models.Plan
	if exists && pkg.PlanID != 0 {
		err = s.db.WithContext(ctx).First(&currentPlan, "id = ?", pkg.PlanID).Error
		if err != nil {
			s.logger.With(ctx).Error("خطایی در سیستم رخ داده است. لطفا دوباره امتحان کنید.")
			return "", response.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "خطایی در سیستم رخ داده است. لطفا دوباره امتحان کنید.",
			}
		}
	}

	if exists {
		// Check if the new plan is a downgrade, same plan (renewal), or upgrade
		isDowngrade := planHierarchy[plan.Title] < planHierarchy[currentPlan.Title]
		isRenewal := plan.ID == currentPlan.ID

		// Prevent downgrades (unless it's a renewal)
		if isDowngrade && !isRenewal {
			return "", response.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "امکان خرید پلن با رده پایین تر از پلن فعلی وجود ندارد.",
			}
		}

		// NEW: Prevent renewal if package is still active (not expired)
		if isRenewal {
			// Check if package is still active (not expired)
			if pkg.ExpiredAt.Valid && pkg.ExpiredAt.Time.After(time.Now()) {
				return "", response.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message: fmt.Sprintf("پکیج فعلی شما هنوز تا تاریخ %s فعال است. امکان تمدید قبل از انقضا وجود ندارد. برای ارتقای پلن، پلن بالاتری انتخاب کنید.",
						pkg.ExpiredAt.Time.Format("2006/01/02")),
				}
			}
		}

		// Allow upgrades regardless of expiration status
		// زمان برای upgrade ها اضافه میشه در service layer (invoice service)
	}

	// if organization has package we cant proceed with creating package
	if len(organization.Packages) > 0 {
		s.logger.With(ctx).Error(err)
		return "", response.ErrorBadRequest("سازمان مورد نظر قبلا بسته خریداری کرده است")
	}

	// TODO: update invoice creation with new org data ayman is doing
	// issue invoice
	invoice, er := s.invoice.Issue(ctx, invoiceSvc.CreateInvoiceRequest{
		InvoiceUniqueCode: uuid.New().String(),
		FromName:          "آرمو",
		FromAddress:       "ایران، تهران، میدان آزادی، اتوبان لشگری، بعد از ایستگاه مترو بیمه، پلاک ۳۱، کارخانه نوآوری آزادی",
		FromPhoneNumber:   "02128424173",
		FromEmail:         "hello@armogroup.tech",
		FromPostalCode:    "1391950015",
		Seller:            "آفاق روشم مرژهای واقعیت",
		EconomicID:        "-",
		RegisterNumber:    "635177",
		ToName:            exp.TerIf(organization.IsIndividual, organization.User.Name+" "+organization.User.LastName, organization.CompanyName),
		ToAddress:         exp.TerIf(organization.IsIndividual, organization.IndividualAddress, organization.LegalAddress),
		ToPhoneNumber:     exp.TerIf(organization.IsIndividual, organization.User.Phone, organization.PhoneNumber),
		ToEmail:           exp.TerIf(organization.IsIndividual, organization.User.Email.String, organization.Email),
		ToPostalCode:      organization.ZipCode,
		Status:            "pending",
		CouponCode:        input.CouponCode,
		TaxPercentage:     10,
		Suspended:         false,
		OrganizationID:    int(organization.ID),
		InvoiceItems: []invoiceSvc.InvoiceItemData{
			{
				Title: plan.Title,
				Description: fmt.Sprintf("خرید طرح %s با شناسه پلن (ID: %d) و دسته‌بندی %s توسط کاربر %s (ID: %d) برای سازمان %s (ID: %d)",
					plan.Title,
					plan.ID,
					exp.TerIf(len(plan.Categories) > 0, plan.Categories[0].Title, "بدون دسته‌بندی"),
					organization.User.Name+" "+organization.User.LastName,
					organization.User.ID,
					organization.CompanyName,
					organization.ID),

				TotalPrice:      plan.Price,
				DiscountedPrice: plan.DiscountedPrice,
				OwnerID:         int(plan.ID),
				OwnerType:       "plan",
				OrganizationID:  int(organization.ID),
			},
		},
	})
	if er.StatusCode != 0 {
		return "", er
	}

	redirectLink, er := s.invoice.Pay(ctx, invoiceSvc.PayInvoiceRequest{
		InvoiceID: int(invoice.ID),
		Gateway:   "payping",
	})
	if er.StatusCode != 0 {
		return "", er
	}

	return redirectLink, response.ErrorResponse{}
}
