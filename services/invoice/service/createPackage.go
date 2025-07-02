package service

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type CreateInvoicePackageRequest struct {
	PlanID         int `json:"plan_id"`
	OrganizationID int `json:"organization_id"`
}

func (c *CreateInvoicePackageRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"plan_id":         govalidity.New("plan_id").Required(),
		"organization_id": govalidity.New("organization_id").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"plan_id":         "Plan ID",
			"organization_id": "Organization",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)
	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}
	return nil
}

func (s *invoice) CreatePackage(ctx context.Context, input CreateInvoicePackageRequest) (models.Package, response.ErrorResponse) {
	var plan models.Plan
	err := s.db.WithContext(ctx).Preload("Categories").First(&plan, "id =?", input.PlanID).Error
	if err != nil {
		return models.Package{}, response.GormErrorResponse(err, "Error finding the plan")
	}

	var organization models.Organization
	err = s.db.WithContext(ctx).Preload("Category").Preload("Packages").First(&organization, "id =?", input.OrganizationID).Error
	if err != nil {
		return models.Package{}, response.GormErrorResponse(err, "Error finding the organization")
	}

	var pkg models.Package
	exists := false
	err = s.db.WithContext(ctx).Preload("Plan").Where("organization_id = ?", input.OrganizationID).First(&pkg).Error
	if err == nil {
		exists = true
	} else {
		s.logger.With(ctx).Info("No active plan found for you.")
	}

	if exists {
		updatedPkg := updatePackageDetails(&pkg, plan)
		err = s.db.WithContext(ctx).Save(updatedPkg).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return models.Package{}, response.GormErrorResponse(err, "Error updating package")
		}
		s.logger.With(ctx).Info("Package processed", "PlanID", updatedPkg.PlanID, "ProductLimit", updatedPkg.ProductLimit, "StorageLimitMB", updatedPkg.StorageLimitMB)
		return *updatedPkg, response.ErrorResponse{}
	} else {
		// Create a new package
		newPkg := createNewPackage(plan, organization, input)
		err = s.db.WithContext(ctx).Create(&newPkg).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return models.Package{}, response.GormErrorResponse(err, "Error creating package")
		}
		s.logger.With(ctx).Info("Package processed", "PlanID", newPkg.PlanID, "ProductLimit", newPkg.ProductLimit, "StorageLimitMB", newPkg.StorageLimitMB)
		return newPkg, response.ErrorResponse{}
	}
}

func updatePackageDetails(pkg *models.Package, plan models.Plan) *models.Package {

	remainingTime := time.Until(pkg.ExpiredAt.Time)
	remainingDays := int(math.Ceil(remainingTime.Hours() / 24.0))

	// Calculate the consumed amounts if applicable
	consumedProductLimit := pkg.Plan.ProductLimit - pkg.ProductLimit
	consumedStorageLimitMB := pkg.Plan.StorageLimitMB - pkg.StorageLimitMB

	pkg.PlanID = int(plan.ID)
	pkg.ProductLimit = plan.ProductLimit - consumedProductLimit
	pkg.StorageLimitMB = plan.StorageLimitMB - consumedStorageLimitMB
	pkg.Price = plan.Price

	// Update the plan details directly
	pkg.Plan = &models.Plan{
		ID:             plan.ID,
		Title:          plan.Title,
		ProductLimit:   plan.ProductLimit,
		StorageLimitMB: plan.StorageLimitMB,
		Price:          plan.Price,
		DayLength:      plan.DayLength,
		Categories:     plan.Categories,
	}

	todayStart := time.Now().Truncate(24 * time.Hour)
	newExpirationDate := todayStart.AddDate(0, 0, plan.DayLength+remainingDays)
	pkg.ExpiredAt = dtp.NullTime{Time: newExpirationDate, Valid: true}

	fmt.Printf("After update: Package = %+v\n", pkg)

	return pkg
}

func createNewPackage(plan models.Plan, organization models.Organization, input CreateInvoicePackageRequest) models.Package {
	return models.Package{
		UserID:         organization.UserID,
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
}
