package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Delete(ctx context.Context, ids []int) ([]int, response.ErrorResponse) {
	var res []int
	if !policy.CanDeleteInvite(ctx) {
		s.logger.With(ctx).Error("شما اجازه حذف کد معرف ندارید.")
		return res, response.ErrorForbidden("شما اجازه حذف کد معرف ندارید.")
	}

	err := s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&models.Invite{}).Error
	if err != nil {
		return res, response.GormErrorResponse(err, "خطایی در حذف کد معرف رخ داده است.")
	}
	return ids, response.ErrorResponse{}
}
