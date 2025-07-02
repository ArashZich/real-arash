package service

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *invoice) Delete(ctx context.Context, ids []int) ([]int, response.ErrorResponse) {
	if !policy.CanDeleteInvoice(ctx) {
		s.logger.With(ctx).Error("شما دسترسی حذف فاکتور ندارید")
		return []int{}, response.ErrorForbidden("شما دسترسی حذف فاکتور ندارید")
	}

	err := s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&models.Invoice{}).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "خطا در حذف فاکتور")
	}

	return ids, response.ErrorResponse{}
}
