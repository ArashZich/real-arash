// services/post/endpoints/delete.go
package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Delete(ctx context.Context, ids []uint) ([]uint, response.ErrorResponse) {
	if !policy.CanDeletePost(ctx) {
		s.logger.With(ctx).Error("شما دسترسی حذف پست را ندارید")
		return []uint{}, response.ErrorForbidden("شما دسترسی حذف پست را ندارید")
	}

	err := s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&models.Post{}).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []uint{}, response.GormErrorResponse(err, "خطا در حذف پست")
	}

	return ids, response.ErrorResponse{}
}
