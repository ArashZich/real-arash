package endpoints

import (
	"context"
	"math"
	"net/http"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type SetEnterpriseRequest struct {
	UserID         int `json:"user_id"`
	OrganizationID int `json:"organization_id"`
}

func (c *SetEnterpriseRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"user_id":         govalidity.New("user_id").Required(),
		"organization_id": govalidity.New("organization_id").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"user_id":         "شناسه کاربر",
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

func (s *service) SetEnterprise(ctx context.Context, input SetEnterpriseRequest) (models.Package, response.ErrorResponse) {

	if !policy.CanCreatePackage(ctx) {
		s.logger.With(ctx).Error("شما دسترسی خرید بسته را ندارید")
		return models.Package{}, response.ErrorForbidden("شما دسترسی ایجاد سازمان را ندارید")
	}

	var plan models.Plan
	err := s.db.WithContext(ctx).Preload("Categories").First(&plan, "title = ?", "enterprise").Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Package{}, response.GormErrorResponse(err, "Error finding enterprise plan")
	}

	var organization models.Organization
	err = s.db.WithContext(ctx).Preload("Category").Preload("Packages").First(&organization, "id = ?", input.OrganizationID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Package{}, response.GormErrorResponse(err, "Error finding organization")
	}

	var pkg models.Package
	exists := false
	err = s.db.WithContext(ctx).Preload("Plan").Where("user_id = ? AND organization_id = ?", input.UserID, input.OrganizationID).First(&pkg).Error
	if err == nil {
		exists = true
	} else {
		s.logger.With(ctx).Info("پلن فعالی برای کاربر ثبت نشده است.")
	}

	if exists {
		// Calculate the remaining time until the package expires
		var remainingTime time.Duration
		if pkg.ExpiredAt.Valid {
			remainingTime = time.Until(pkg.ExpiredAt.Time)
		}

		// Calculate the remaining days, rounding up to include the last day fully
		remainingDays := int(math.Ceil(remainingTime.Hours() / 24.0))

		// Directly assign the enterprise plan's limits to the package
		pkg.PlanID = int(plan.ID)
		pkg.ProductLimit = plan.ProductLimit
		pkg.StorageLimitMB = plan.StorageLimitMB
		pkg.Price = plan.Price
		pkg.Plan = &plan

		// Adding the total number of days of the new plan to the number of days remaining
		todayStart := time.Now().Truncate(24 * time.Hour)
		newExpirationDate := todayStart.AddDate(0, 0, plan.DayLength+remainingDays)

		pkg.ExpiredAt = dtp.NullTime{Time: newExpirationDate, Valid: true}

		err = s.db.WithContext(ctx).Save(&pkg).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return models.Package{}, response.GormErrorResponse(err, "خطا در بروزرسانی پلن به اینترپرایز")
		}
	} else {
		// Create a new package
		pkg = models.Package{
			UserID:         input.UserID,
			PlanID:         int(plan.ID),
			ProductLimit:   plan.ProductLimit,
			StorageLimitMB: plan.StorageLimitMB,
			Price:          plan.Price,
			OrganizationID: input.OrganizationID,
			CategoryID:     organization.CategoryID,
			ExpiredAt: dtp.NullTime{
				Time:  time.Now().AddDate(0, 0, plan.DayLength),
				Valid: true,
			},
		}
		err = s.db.WithContext(ctx).Create(&pkg).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return models.Package{}, response.GormErrorResponse(err, "خطا در ایجاد پلن اینترپرایز")
		}
	}

	s.logger.With(ctx).Info("Enterprise package processed", "UserID", pkg.UserID, "PlanID", pkg.PlanID, "ProductLimit", pkg.ProductLimit, "StorageLimitMB", pkg.StorageLimitMB)
	return pkg, response.ErrorResponse{}
}
