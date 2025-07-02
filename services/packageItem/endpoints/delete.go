package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *packageItem) Delete(ctx context.Context, ids []int) ([]int, response.ErrorResponse) {
	if !policy.CanDeletePackageItem(ctx) {
		s.logger.With(ctx).Error("شما دسترسی حذف آیتم ندارید")
		return []int{}, response.ErrorForbidden("شما دسترسی حذف آیتم ندارید")
	}

	err := s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&models.PackageItem{}).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "خطا در حذف آیتم")
	}

	return ids, response.ErrorResponse{}
}
