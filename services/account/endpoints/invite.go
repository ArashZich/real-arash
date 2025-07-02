package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) useInviteCode(ctx context.Context, code string) (models.Invite, response.ErrorResponse) {
	var invite models.Invite
	err := s.db.WithContext(ctx).First(&invite, "code", code).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return invite, response.GormErrorResponse(err, "کد معرفی یافت نشد")
	}

	if !(invite.Limit > 0) && invite.Limit != -1 {
		return invite, response.ErrorBadRequest(nil, "کد معرف شما تمام شده است.")
	}
	if invite.Limit == -1 {
		return invite, response.ErrorResponse{}
	}
	if invite.Limit != -1 {
		invite.Limit = invite.Limit - 1
		err = s.db.WithContext(ctx).Save(&invite).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return invite, response.GormErrorResponse(err, "خطا در ثبت کد معرف")
		}
		return invite, response.ErrorResponse{}
	}
	return invite, response.ErrorResponse{}
}
