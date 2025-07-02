package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Check(ctx context.Context, code string) (models.Invite, response.ErrorResponse) {
	var invite models.Invite

	err := s.db.WithContext(ctx).First(&invite, "code", code).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return invite, response.GormErrorResponse(nil, "دعوتنامه‌ای با این کد یافت نشد")
	}

	return invite, response.ErrorResponse{}
}
