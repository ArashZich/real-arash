package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Delete(ctx context.Context, ids []int) ([]int, response.ErrorResponse) {

	if !policy.CanDeleteView(ctx) {
		s.logger.With(ctx).Error("شما دسترسی حذف بازرسی را ندارید")
		return []int{}, response.ErrorForbidden("شما دسترسی حذف بازرسی را ندارید")
	}

	err := s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&models.View{}).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "خطا در حذف بازرسی")
	}

	return ids, response.ErrorResponse{}
}
