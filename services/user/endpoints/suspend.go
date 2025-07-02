package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/faker"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Suspend(ctx context.Context, id string) (models.User, response.ErrorResponse) {
	var user models.User
	if !policy.CanSuspendUser(ctx) {
		s.logger.With(ctx).Error("شما اجازه دسترسی به این کاربر را ندارید")
		err := response.ErrorForbidden("شما اجازه دسترسی به این کاربر را ندارید")
		return user, err
	}

	err := s.db.WithContext(ctx).First(&user, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}

	if user.SuspendedAt.Valid {
		user.SuspendedAt = dtp.NullTime{}
	} else {
		user.SuspendedAt = faker.SQLNow()
	}

	err = s.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "خطا در ذخیره کاربر")
	}

	return user, response.ErrorResponse{}
}
