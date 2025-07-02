package endpoints

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/google/uuid"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type UpdateVisitDurationRequest struct {
	ProductUID     uuid.UUID `json:"product_uid"`
	OrganizationID int       `json:"organization_id"`
	VisitDuration  int64     `json:"visit_duration"` // Duration of the visit in seconds
	VisitUID       uuid.UUID `json:"visit_uid"`
}

func (c *UpdateVisitDurationRequest) ValidateDuration(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"visit_duration":  govalidity.New("visit_duration").Required().Min(0),
		"product_uid":     govalidity.New("product_uid").Optional(),
		"organization_id": govalidity.New("organization_id").Required(),
		"visit_uid":       govalidity.New("visit_uid").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"visit_duration":  "مدت زمان بازدید",
			"product_uid":     "محصول",
			"organization_id": "سازمان",
			"visit_uid":       "شناسه بازدید",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)
	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}
	return nil
}

func (s *service) CreateOrUpdateDuration(ctx context.Context, input UpdateVisitDurationRequest) (models.View, response.ErrorResponse) {
	// Find the existing view by VisitUID
	var view models.View
	err := s.db.WithContext(ctx).Where("visit_uid = ?", input.VisitUID).First(&view).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// If no view matches the VisitUID, return an error indicating it doesn't exist
		return models.View{}, response.GormErrorResponse(err, "No view found with the specified VisitUID")
	} else if err != nil {
		// Handle any other database error
		s.logger.With(ctx).Error(err)
		return models.View{}, response.GormErrorResponse(err, "Error finding view")
	}

	// Update the VisitDuration of the found view
	view.VisitDuration += input.VisitDuration

	// Save the updated view record
	if err := s.db.WithContext(ctx).Save(&view).Error; err != nil {
		s.logger.With(ctx).Error(err)
		return models.View{}, response.GormErrorResponse(err, "Error updating view duration")
	}

	return view, response.ErrorResponse{}
}
