package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) UserInfo(ctx context.Context, accessToken string) (models.User, response.ErrorResponse) {
	Id := policy.ExtractIdClaim(ctx)
	var user models.User
	var token models.Token
	err := s.db.WithContext(ctx).First(&token, "access_token", accessToken).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.ErrorUnAuthorized(nil, "توکن ارسالی معتبر نمی‌باشد")
	}

	err = s.db.WithContext(ctx).
		Preload("Roles").
		Preload("Invite").
		Preload("Organizations.Packages.Category").
		First(&user, "id", Id).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.ErrorUnAuthorized(nil, "کاربری با این مشخصات یافت نشد")
	}

	if user.SuspendedAt.Valid {
		s.logger.With(ctx).Error(err)
		return user, response.ErrorUnAuthorized(nil, "حساب کاربری شما مسدود شده است")
	}

	return user, response.ErrorResponse{}
}
