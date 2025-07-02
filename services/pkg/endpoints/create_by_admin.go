// create_by_admin-endpoints.go
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

type CreatePackageByAdminRequest struct {
	UserID         int `json:"user_id"`
	PlanID         int `json:"plan_id"`
	OrganizationID int `json:"organization_id"`
}

func (c *CreatePackageByAdminRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"user_id":         govalidity.New("user_id").Required(),
		"plan_id":         govalidity.New("plan_id").Required(),
		"organization_id": govalidity.New("organization_id").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"user_id":         "شناسه کاربر",
			"plan_id":         "شناسه طرح",
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

func (s *service) CreateByAdmin(ctx context.Context, input CreatePackageByAdminRequest) (models.Package, response.ErrorResponse) {
	// Check if user is super admin
	if !policy.CanCreatePackage(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ایجاد بسته برای کاربران را ندارید")
		return models.Package{}, response.ErrorForbidden("شما دسترسی ایجاد بسته برای کاربران را ندارید")
	}

	// Verify user exists
	var user models.User
	err := s.db.WithContext(ctx).First(&user, "id = ?", input.UserID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Package{}, response.GormErrorResponse(err, "Error finding user")
	}

	// Find plan with its categories
	var plan models.Plan
	err = s.db.WithContext(ctx).Preload("Categories").First(&plan, "id = ?", input.PlanID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Package{}, response.GormErrorResponse(err, "Error finding plan")
	}

	// Find organization with its category and packages
	var organization models.Organization
	err = s.db.WithContext(ctx).Preload("Category").Preload("Packages").First(&organization, "id = ?", input.OrganizationID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Package{}, response.GormErrorResponse(err, "Error finding organization")
	}

	// Check for existing package
	var pkg models.Package
	exists := false
	err = s.db.WithContext(ctx).Preload("Plan").Where("user_id = ? AND organization_id = ?", input.UserID, input.OrganizationID).First(&pkg).Error
	if err == nil {
		exists = true
		s.logger.With(ctx).Info("Existing package found for user", "UserID", input.UserID)
	}

	// Define plan hierarchy
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
			return models.Package{}, response.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "خطایی در سیستم رخ داده است. لطفا دوباره امتحان کنید.",
			}
		}
	}

	if exists {
		// Check if the new plan is a downgrade or the same plan for renewal
		isDowngrade := planHierarchy[plan.Title] < planHierarchy[currentPlan.Title]
		isRenewal := plan.ID == currentPlan.ID

		if isDowngrade && !isRenewal {
			return models.Package{}, response.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "امکان خرید پلن با رده پایین تر از پلن فعلی وجود ندارد.",
			}
		}

		// Calculate the remaining time until the package expires
		var remainingTime time.Duration
		if pkg.ExpiredAt.Valid {
			remainingTime = time.Until(pkg.ExpiredAt.Time)
		}

		// Calculate the remaining days, rounding up to include the last day fully
		remainingDays := int(math.Ceil(remainingTime.Hours() / 24.0))

		// Calculate the consumed amounts if applicable
		consumedProductLimit := pkg.Plan.ProductLimit - pkg.ProductLimit
		consumedStorageLimitMB := pkg.Plan.StorageLimitMB - pkg.StorageLimitMB

		// Update existing package
		pkg.PlanID = int(plan.ID)
		pkg.ProductLimit = plan.ProductLimit - consumedProductLimit
		pkg.StorageLimitMB = plan.StorageLimitMB - consumedStorageLimitMB
		pkg.Price = plan.Price
		pkg.Plan = &plan

		// Adding the total number of days of the new plan to the number of days remaining
		todayStart := time.Now().Truncate(24 * time.Hour)
		newExpirationDate := todayStart.AddDate(0, 0, plan.DayLength+remainingDays)
		pkg.ExpiredAt = dtp.NullTime{Time: newExpirationDate, Valid: true}

		err = s.db.WithContext(ctx).Save(&pkg).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return models.Package{}, response.GormErrorResponse(err, "خطا در بروزرسانی پلن")
		}
	} else {
		// Create new package
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
			return models.Package{}, response.GormErrorResponse(err, "خطا در ایجاد پلن")
		}
	}

	// Save or create the package in the database
	if exists {
		err = s.db.WithContext(ctx).Save(&pkg).Error
	}

	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Package{}, response.GormErrorResponse(err, "خطا در ذخیره سازی پلن")
	}

	s.logger.With(ctx).Info("Package processed by admin", "UserID", pkg.UserID, "PlanID", pkg.PlanID, "ProductLimit", pkg.ProductLimit, "StorageLimitMB", pkg.StorageLimitMB)
	return pkg, response.ErrorResponse{}
}
