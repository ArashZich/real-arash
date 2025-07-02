package endpoints

import (
	"context"
	"fmt"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) QueryByCategoryID(ctx context.Context, categoryID int) ([]models.Plan, response.ErrorResponse) {
	var plans []models.Plan
	where := fmt.Sprintf("EXISTS (SELECT 1 FROM plan_categories WHERE plan_categories.plan_id = plans.id AND plan_categories.category_id = %d)", categoryID)
	err := s.db.WithContext(ctx).Where(where).Find(&plans).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Plan{}, response.GormErrorResponse(err, "خطایی در یافتن طرح‌ها رخ داده است")
	}

	return plans, response.ErrorResponse{}
}
