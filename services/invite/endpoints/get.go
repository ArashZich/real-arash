package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Get(ctx context.Context, id string) (models.Invite, response.ErrorResponse) {
	var invite models.Invite
	if !policy.CanGetInvite(ctx) {
		s.logger.With(ctx).Error("شما اجازه دسترسی به این کد معرف را ندارید.")
		err := response.ErrorForbidden("شما دسترسی لازم برای انجام این کار را ندارید")
		return invite, err
	}
	err := s.db.WithContext(ctx).First(&invite, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return invite, response.GormErrorResponse(err, "کد معرف یافت نشد.")
	}

	return invite, response.ErrorResponse{}
}
