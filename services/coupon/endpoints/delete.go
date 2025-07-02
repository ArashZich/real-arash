package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Delete(ctx context.Context, ids []int) ([]int, response.ErrorResponse) {
	if !policy.CanDeleteCoupon(ctx) {
		s.logger.With(ctx).Error("شما دسترسی حذف دسته بندی ندارید")
		return []int{}, response.ErrorForbidden("شما دسترسی حذف دسته بندی ندارید")
	}

	err := s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&models.Coupon{}).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "خطا در حذف دسته")
	}

	return ids, response.ErrorResponse{}
}
