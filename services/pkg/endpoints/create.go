package endpoints

import (
	"context"
	"math"
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

type CreatePackageRequest struct {
	PlanID         int `json:"plan_id"`
	OrganizationID int `json:"organization_id"`
}

func (c *CreatePackageRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"plan_id":         govalidity.New("plan_id").Required(),
		"organization_id": govalidity.New("organization_id").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
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

func (s *service) Create(ctx context.Context, input CreatePackageRequest) (models.Package, response.ErrorResponse) {

	// if !policy.CanCreatePackage(ctx) {
	// 	s.logger.With(ctx).Error("شما دسترسی خرید بسته را ندارید")
	// 	return models.Package{}, response.ErrorForbidden("شما دسترسی ایجاد سازمان را ندارید")
	// }

	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	var plan models.Plan
	err := s.db.WithContext(ctx).Preload("Categories").First(&plan, "id = ?", input.PlanID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Package{}, response.GormErrorResponse(err, "Error finding plan")
	}

	var organization models.Organization
	err = s.db.WithContext(ctx).Preload("Category").Preload("Packages").First(&organization, "id = ?", input.OrganizationID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Package{}, response.GormErrorResponse(err, "Error finding organization")
	}

	var pkg models.Package
	exists := false
	err = s.db.WithContext(ctx).Preload("Plan").Where("user_id = ? AND organization_id = ?", id, input.OrganizationID).First(&pkg).Error
	if err == nil {
		exists = true
	} else {
		s.logger.With(ctx).Info("پلن فعالی برای شما ثبت نشده است.")
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
			s.logger.With(ctx).Error("خطایی در سیستم ره داده است. لطفا دوباره امتحان کنید.")
			return models.Package{}, response.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "خطایی در سیستم ره داده است. لطفا دوباره امتحان کنید.",
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

		// Directly assign the new plan's limits to the package
		pkg.PlanID = int(plan.ID)
		pkg.ProductLimit = plan.ProductLimit - consumedProductLimit
		pkg.StorageLimitMB = plan.StorageLimitMB - consumedStorageLimitMB
		pkg.Price = plan.Price
		pkg.Plan = &plan
		// Adding the total number of days of the new plan to the number of days remaining
		// Ensure to start calculation from the beginning of the current day to include today fully
		todayStart := time.Now().Truncate(24 * time.Hour)
		newExpirationDate := todayStart.AddDate(0, 0, plan.DayLength+remainingDays)

		pkg.ExpiredAt = dtp.NullTime{Time: newExpirationDate, Valid: true}

		err = s.db.WithContext(ctx).Save(&pkg).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return models.Package{}, response.GormErrorResponse(err, "خطا در بروزرسانی پلن")
		}
	} else {
		// Create a new package
		pkg = models.Package{
			UserID:         id,
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
			return models.Package{}, response.GormErrorResponse(err, "خطا در خرید پلن")
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

	s.logger.With(ctx).Info("Package processed", "PlanID", pkg.PlanID, "ProductLimit", pkg.ProductLimit, "StorageLimitMB", pkg.StorageLimitMB)
	return pkg, response.ErrorResponse{}

}
